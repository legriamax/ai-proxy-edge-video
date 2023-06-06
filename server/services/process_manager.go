// Copyright 2020 Wearless Tech Inc All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package services

import (
	"context"
	"encoding/json"
	"errors"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/chryscloud/go-microkit-plugins/docker"
	g "github.com/chryscloud/video-edge-ai-proxy/globals"
	"github.com/chryscloud/video-edge-ai-proxy/models"
	"github.com/dgraph-io/badger/v2"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	dockerErrors "github.com/docker/docker/client"
	"github.com/go-redis/redis/v7"
)

const (
	// Resource: https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
	ArchitectureAmd64 = "amd64"
	ArchitectureArm64 = "arm64"
)

var (
	ArchitectureSuffixMap = map[string]string{ArchitectureAmd64: "", ArchitectureArm64: "-arm64v8"}
)

// ProcessManager - start, stop of docker containers
type ProcessManager struct {
	storage *Storage
	rdb     *redis.Client
}

func NewProcessManager(storage *Storage, rdb *redis.Client) *ProcessManager {
	return &ProcessManager{
		storage: storage,
		rdb:     rdb,
	}
}

// Start - starts the docker container with rtsp, device_id and possibly rtmp environment variables.
// Restarts always when something goes wrong within the streaming process
func (pm *ProcessManager) Start(process *models.StreamProcess, imageUpgrade *models.ImageUpgrade) error {

	// detect architecture
	arch := runtime.GOARCH

	if _, ok := ArchitectureSuffixMap[arch]; !ok {
		return errors.New("architecture currently not supported")
	}

	if process.Name == "" || process.RTSPEndpoint == "" {
		return errors.New("required parameters missing")
	}

	if !imageUpgrade.HasImage && !imageUpgrade.HasUpgrade {
		return errors.New("no camera container found. Please refer to documentation on how to pull a docker image manually")
	}

	settingsTagBytes, err := pm.storage.Get(models.PrefixSettingsDockerTagVersions, "rtsp")
	if err != nil {
		if err == badger.ErrKeyNotFound {

			// if no docker tag version stored in database but image does exist on disk, then store settings docker tag version with that image
			tag := models.CameraTypeToImageTag["rtsp"]
			if imageUpgrade == nil {
				return errors.New("Image not found. Please check the docs and pull the docker image manually.")
			}
			maximumExistingTag := tag + ":" + imageUpgrade.CurrentVersion
			// store to database
			g.Log.Info("maximum existing tag od disk found: ", maximumExistingTag)

			settingsTagVersion := &models.SettingDockerTagVersion{
				CameraType: "rtsp",
				Tag:        tag,
				Version:    imageUpgrade.CurrentVersion,
			}
			stb, sErr := pm.storeSettingsTagVersion(settingsTagVersion)
			if sErr != nil {
				g.Log.Error("failed to store new settings tag version ", sErr)
				return sErr
			}

			settingsTagBytes = stb
		} else {
			g.Log.Error("failed to read rtsp tag from settings", err)
			return err
		}
	}

	var settingsTag models.SettingDockerTagVersion
	err = json.Unmarshal(settingsTagBytes, &settingsTag)
	if err != nil {
		g.Log.Error("failed to unamrshal settings tag", err)
		return err
	}
	process.ImageTag = settingsTag.Tag + ":" + settingsTag.Version

	// Check the latest version that exists on the disk (and if is the same as the one in settings)
	// if is not, correct the latest version stored (most likely user chose to manually deleted the newer version)
	if imageUpgrade.CurrentVersion != settingsTag.Version {
		settingsTag.Version = imageUpgrade.CurrentVersion

		process.ImageTag = imageUpgrade.Name + ":" + imageUpgrade.CurrentVersion

		_, sErr := pm.storeSettingsTagVersion(&settingsTag)
		if sErr != nil {
			g.Log.Error("failed to store new settings tag", sErr, ", image version: ", settingsTag.Tag, settingsTag.Version)
			return sErr
		}
	}

	cl := docker.NewSocketClient(docker.Log(g.Log), docker.Host("unix:///var/run/docker.sock"))

	fl := filters.NewArgs()
	pruneReport, pruneErr := cl.ContainersPrune(fl)
	if pruneErr != nil {
		g.Log.Error("container prunning fialed", pruneErr)
		return pruneErr
	}
	g.Log.Info("prune successfull. Report and space reclaimed:", pruneReport.ContainersDeleted, pruneReport.SpaceReclaimed)

	hostConfig := &container.HostConfig{
		// PortBindings: mappingPorts,
		LogConfig: container.LogConfig{
			Type:   "json-file",
			Config: map[string]string{"max-file": "3", "max-size": "3M"},
		},
		RestartPolicy: container.RestartPolicy{Name: "always"},
		Resources: container.Resources{
			CPUShares: 1024, // equal weight to all containers. check here the docs here:  https://docs.docker.com/config/containers/resource_constraints/
		},
		NetworkMode: container.NetworkMode("chrysnet"),
	}

	if g.Conf.Buffer.OnDisk {
		mounts := make([]mount.Mount, 0)
		mount := mount.Mount{
			Type:     mount.TypeBind,
			Source:   g.Conf.Buffer.OnDiskFolder,
			Target:   g.Conf.Buffer.OnDiskFolder,
			ReadOnly: false,
		}
		mounts = append(mounts, mount)

		hostConfig.Mounts = mounts
	}

	envVars := []string{"rtsp_endpoint=" + process.RTSPEndpoint, "device_id=" + process.Name, "in_memory_buffer=" + strconv.Itoa(g.Conf.Buffer.InMemory)}
	if process.RTMPEndpoint != "" {
		envVars = append(envVars, "rtmp_endpoint="+process.RTMPEndpoint)
	}
	if g.Conf.Buffer.OnDisk {
		envVars = append(envVars, "disk_buffer_path="+g.Conf.Buffer.OnDiskFolder)
		envVars = append(envVars, "disk_cleanup_rate="+g.Conf.Buffer.OnDiskCleanupOlderThan)
	}
	if g.Conf.Redis.Connection != "" {
		host := strings.Split(g.Conf.Redis.Connection, ":")
		if len(host) == 2 {
			envVars = append(envVars, "redis_host="+host[0])
			envVars = append(envVars, "redis_port="+host[1])
		}
	}
	if g.Conf.Buffer.InMemoryScale != "" {
		envVars = append(envVars, "memory_scale="+g.Conf.Buffer.InMemoryScale)
	}

	envVars = append(envVars, "PYTHONUNBUFFERED=0") // for output to console

	_, ccErr := cl.ContainerCreate(strings.ToLower(process.Name), &container.Config{
		Image: process.ImageTag,
		Env:   envVars,
	}, hostConfig, nil)

	if ccErr != nil {
		g.Log.Error("failed to create container ", process.Name, ccErr)
		return ccErr
	}

	err = cl.ContainerStart(process.Name)
	if err != nil {
		g.Log.Error("failed to start container", process.Name, err)
		return err
	}

	process.Status = "running"
	process.Created = time.Now().Unix() * 1000

	// set default value in redis if RTMP streaming enabled
	if process.RTMPEndpoint != "" {
		valMap := make(map[string]interface{}, 0)
		valMap[models.RedisLastAccessQueryTimeKey] = time.Now().Unix() * 1000
		valMap[models.RedisProxyRTMPKey] = true

		rErr := pm.rdb.HSet(models.RedisLastAccessPrefix+process.Name, valMap).Err()
		if rErr != nil {
			g.Log.Error("failed to store startproxy value map to redis", rErr)
			return rErr
		}
		if process.RTMPStreamStatus == nil {
			process.RTMPStreamStatus = &models.RTMPStreamStatus{}
		}
		process.RTMPStreamStatus.Streaming = true
	}

	obj, err := json.Marshal(process)
	if err != nil {
		g.Log.Error("failed to marshal process json", err)
		return err
	}

	err = pm.storage.Put(models.PrefixRTSPProcess, process.Name, obj)
	if err != nil {
		g.Log.Error("failed to store process into datastore", err)
		return err
	}

	return nil
}

// Stop - stops the docker container by the name of deviceID and removed from local datastore
// databasePrefix = models.PrefixRTSPProcess or models.PrefixAppProcess
func (pm *ProcessManager) Stop(deviceID string, databasePrefix string) error {
	cl := docker.NewSocketClient(docker.Log(g.Log), docker.Host("unix:///var/run/docker.sock"))

	container, err := cl.ContainerGet(deviceID)
	if err != nil {
		if dockerErrors.IsErrNotFound(err) {
			g.Log.Info("container not found to be stopeed", err)
			return models.ErrProcessNotFound
		}
	}

	// waits up to 10 minutes to stop the container, otherwise kills after 30 seconds
	killAfer := time.Second * 5
	err = cl.ContainerStop(container.ID, &killAfer)
	if err != nil {
		if dockerErrors.IsErrNotFound(err) {
			g.Log.Warn("container doesn't exist. probably stopped before", err)
			return nil
		}
	}

	err = pm.storage.Del(databasePrefix, deviceID)
	if err != nil {
		g.Log.Error("Failed to delete rtsp proces", err)
		return err
	}

	fl := filters.NewArgs()
	pruneReport, pruneErr := cl.ContainersPrune(fl)
	if pruneErr != nil {
		g.Log.Error("container prunning fialed", pruneErr)
		return pruneErr
	}
	g.Log.Info("prune successfull. Report and space reclaimed:", pruneReport.ContainersDeleted, pruneReport.SpaceReclaimed)

	return nil
}

// ListStream - GRPC method for list all streams (doesn't alter the actual processes)
func (pm *ProcessManager) ListStream(ctx context.Context, fo