package async

import (
	"github.com/pkg/errors"
	"time"
)

const WINSIZE = 4194304 // 1024 * 1024 * 4

type Client interface {
	Run() error
}

func NewClient(mode, address, logTime string, packSize, packCount, debugMode int) (Client, error) {
	if mode == "socket" {
		return NewSocketClient(address, logTime, packSize, packCount, debugMode), nil
	} else if mode == "grpc" {
		return NewGrpcClient(address, logTime, packSize, packCount, debugMode), nil
	}
	return nil, errors.New("invalid server mode")
}

type timeElapsed struct {
	elapsedTime   time.Duration
	nSendCntTotal int
}
