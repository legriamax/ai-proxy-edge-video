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
func (gih *grpcImageHandler) Storage(ctx context.Context, req *pb.StorageRequest) (*pb.StorageResponse, error) {

	deviceID := req.DeviceId

	if deviceID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "device id required")
	}

	info, err := gih.processManager.Info(deviceID)
	if err != nil {
		g.Log.Error("failed to get deviceID info", err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	if info.RTMPEndpoint == "" {
		return nil, status.Errorf(codes.InvalidArgument, "device "+deviceID+" doesn't have an associated RTMP stream")
	}

	apiErr := gih.enableDisableStorageAPICall(req.Start, info.RTMPEndpoint)
	if apiErr != nil {
		if apiErr == models.ErrForbidden {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("cannot enable