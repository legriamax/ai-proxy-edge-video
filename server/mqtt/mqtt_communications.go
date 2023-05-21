
package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	g "github.com/chryscloud/video-edge-ai-proxy/globals"
	"github.com/chryscloud/video-edge-ai-proxy/models"
	"github.com/chryscloud/video-edge-ai-proxy/services"
	"github.com/chryscloud/video-edge-ai-proxy/utils"
	badger "github.com/dgraph-io/badger/v2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	qtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v7"
)

const (
	mqttBrokerURL   = "tls://mqtt.googleapis.com:8883"
	protocolVersion = 4 // corresponds to MQTT 3.1.1
)

// ProcessManager - start, stop of docker containers
type mqttManager struct {
	rdb                      *redis.Client
	settingsService          *services.SettingsManager
	processService           *services.ProcessManager
	appService               *services.AppProcessManager
	client                   *qtt.Client
	clientOpts               *qtt.ClientOptions
	stop                     chan bool
	gatewayID                string
	projectID                string
	jwt                      string
	processEvents            sync.Map
	lastProcessEventNotified sync.Map
	mutex                    sync.Mutex
}

func NewMqttManager(rdb *redis.Client, settingsService *services.SettingsManager, processService *services.ProcessManager, appService *services.AppProcessManager) *mqttManager {
	return &mqttManager{
		rdb:                      rdb,
		settingsService:          settingsService,
		processService:           processService,
		appService:               appService,
		processEvents:            sync.Map{},
		lastProcessEventNotified: sync.Map{},
		mutex:                    sync.Mutex{},
	}
}

func (mqtt *mqttManager) onConnect(client qtt.Client) {
	g.Log.Info("MQTT client connected", client.IsConnected())
}

func (mqtt *mqttManager) onMessage(client qtt.Client, msg qtt.Message) {
	g.Log.Info("Command received from Chrysalis Cloud:", msg.Topic())

	var edgeConfig models.EdgeCommandPayload
	err := json.Unmarshal(msg.Payload(), &edgeConfig)
	if err != nil {
		g.Log.Error("failed to unmarshal config payload", err, string(msg.Payload()))
		return
	}

	// mapping to local process types for cameras
	operation := ""
	if edgeConfig.Type == models.ProcessTypeRTSP {

		if edgeConfig.Operation == "a" {
			operation = models.DeviceOperationStart
		} else if edgeConfig.Operation == "r" {
			operation = models.DeviceOperationDelete
		} else {
			g.Log.Error("camera command operation not supported: ", edgeConfig.Name, edgeConfig.ImageTag, edgeConfig.Operation)
			return
		}
	} else {
		// mapping to local process types for applications
		operation = edgeConfig.Operation
	}
	err = utils.PublishToRedis(mqtt.rdb, edgeConfig.Name, models.MQTTProcessOperation(operation), edgeConfig.Type, msg.Payload())
	if err != nil {
		g.Log.Error("failed to process starting of the new device on the edge", err)
	}
}

func (mqtt *mqttManager) onConnectionLost(client qtt.Client, err error) {
	g.Log.Error("MQTT connection lost", err)
}

func (mqtt *mqttManager) configHandler(client qtt.Client, msg qtt.Message) {
	g.Log.Info("Received config request: ", msg.Topic())
	g.Log.Info("Message: ", string(msg.Payload()))
}

// StartGatewayListener checks every 15 seconds if there are any settings for connection to gateway
func (mqtt *mqttManager) StartGatewayListener() error {

	delay := time.Second * 15
	go func() {

		for {
			_, err := mqtt.getMQTTSettings()
			if err == nil {
				mqttErr := mqtt.run()
				if mqttErr != nil {
					g.Log.Error("Failed to init mqtt", mqttErr)
				}
				// exit the waiting function
				break
			}

			select {
			case <-time.After(delay):
			case <-mqtt.stop:
				g.Log.Info("MQTT cron job stopped")
				return
			}
		}
	}()

	return nil
}

func (mqtt *mqttManager) run() error {
	err := mqtt.gatewayInit()

	if err != nil {
		if err == ErrNoMQTTSettings {
			return nil
		}
		g.Log.Error("failed to connect gateway and report presence", err)
		return err
	}

	// init redis listener for local messages (this is only for active local changes)
	// e.g. Device/process added, removed, ...
	sub := mqtt.rdb.Subscribe(models.RedisLocalMQTTChannel)

	go func(sub *redis.PubSub) {

		defer sub.Close()

		for {
			val, err := sub.ReceiveMessage()
			if err != nil {
				g.Log.Error("failed to receive mqtt local pubsub message", err)
			} else {
				g.Log.Info("redis message received: ", val)
				payload := []byte(val.Payload)
				var localMsg models.MQTTMessage
				err := json.Unmarshal(payload, &localMsg)
				if err != nil {
					g.Log.Error("failed to unmarshal internal redis pubsub message", err)
				} else {
					g.Log.Info("Received message object from redis pubsub for mqtt: ", localMsg.DeviceID)
					var opErr error
					if localMsg.ProcessType == models.MQTTProcessType(models.ProcessTypeRTSP) {
						if localMsg.ProcessOperation == models.MQTTProcessOperation(models.DeviceOperationAdd) {

							opErr = mqtt.bindDevice(localMsg.DeviceID, models.MQTTProcessType(models.ProcessTypeRTSP))

						} else if localMsg.ProcessOperation == models.MQTTProcessOperation(models.DeviceOperationRemove) {

							opErr = mqtt.unbindDevice(localMsg.DeviceID, models.MQTTProcessType(models.ProcessTypeRTSP))

						} else if localMsg.ProcessOperation == models.MQTTProcessOperation(models.DeviceOperationUpgradeAvailable) {
							// TODO: TBD
							g.Log.Warn("TBD: process operation upgrade available")
						} else if localMsg.ProcessOperation == models.MQTTProcessOperation(models.DeviceOperationUpgradeFinished) {
							// TODO: TBD
							g.Log.Warn("TBD: process operation upgrade completed/finished")
						} else if localMsg.ProcessOperation == models.MQTTProcessOperation(models.DeviceOperationStart) {

							opErr = mqtt.StartCamera(localMsg.Message)
						} else if localMsg.ProcessOperation == models.MQTTProcessOperation(models.DeviceOperationDelete) {

							opErr = mqtt.StopCamera(localMsg.Message)
						} else if localMsg.ProcessOperation == models.MQTTProcessOperation(models.DeviceInternalTesting) {

							// **********
							// internal testing operations
							// **********
							testErr := mqtt.reportDeviceStateChange(localMsg.DeviceID, models.ProcessStatusRestarting)