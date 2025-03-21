package async

import "github.com/pkg/errors"

const (
	BUFSIZE = 8388608  // 1024 * 1024 * 8
	QUESIZE = 41943040 // 1024 * 1024 * 8 * 5
	WINSIZE = 4194304  // 1024 * 1024 * 4
)

type Server interface {
	Run() error
}

func NewServer(mode, address string, packSize, debugMode int) (Server, error) {
	if mode == "socket" {
		return NewSocketServer(address, debugMode, packSize), nil
	} else if mode == "grpc" {
		return NewGrpcServer(address, debugMode), nil
	}
	return nil, errors.New("invalid server mode")
}
