package grpcapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	g "github.com/chryscloud/video-edge-ai-proxy/globals"
	"github.com/chryscloud/video-edge-ai-proxy/models"
	pb "github.com/chryscloud/video-edge-ai-proxy/proto"
	"github.com/chryscloud/video-edge-ai-proxy/utils"
	"github.com/go-resty/resty/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Storage - enable disable storage on Chrysalis Cloud
func (gih *grpcImageHandler) Storage(ctx context.Context, req *pb.StorageRequest) (*pb.Sto