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

	ac := &AnnotationC