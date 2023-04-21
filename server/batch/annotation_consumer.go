package batch

import (
	"time"

	"github.com/adjust/rmq/v2"
	"github.com/chryscloud/go-microkit-plugins/models/ai"
	g "github.com/chryscloud/video-edge-ai-proxy/globals"
	pb "github.com/chryscloud/video-edge-ai-proxy/proto"
	"github.com/chryscloud/video-edge-ai-proxy/services"
	"github.com/chryscloud/video-edge-ai-proxy/utils"
	"github.com/go-resty/resty/v2"
	"github.com/golang/protobuf/proto"
)

type AnnotationConsumer struct {
	settingsService *services.SettingsManager
	restClient      *resty.Client
	msgQueue        rmq.Queue
}

func NewAnnotationConsumer(tag int, settingsService *services.SettingsManager, msgQueue rmq.Queue) *AnnotationConsumer {
	restClient := resty.New().SetRetryCount(3)

	ac := &AnnotationConsumer{
		settingsService: settingsService,
		restClient:      restClient,
		msgQueue:        msgQueue,
	}

	// check every 5 seconds if any rejected annotations
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				ac.failedAnnotationsTryRedo(<-ticker.C)
			}
		}
	}()

	return ac
}

// checks if any messages need to be reinstated that have failed before (in case of internet outage on the edge)
func (ac *AnnotationConsumer) failedAnnotationsTryRedo(tick time.Time) {
	cnt := ac.msgQueue.ReturnAllRejected()
	if cnt > 0 {
		g.Log.Info("re-queued ", cnt, " of previously rejected annotatons")
	}
}

func (ac *AnnotationConsumer) Consume(batch rmq.Deliveries) {

	if g.Conf.Annotation.Endpoint == "" {
		g.Log.