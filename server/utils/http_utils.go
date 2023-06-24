package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	microCrypto "github.com/chryscloud/go-microkit-plugins/crypto"
	g "github.com/chryscloud/video-edge-ai-proxy/gl