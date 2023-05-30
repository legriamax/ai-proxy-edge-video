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
// Restarts always when something goes wron