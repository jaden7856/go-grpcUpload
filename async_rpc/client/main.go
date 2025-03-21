package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"

	"github.com/jaden7856/go-grpcUpload/async_rpc/client/async"
	"github.com/pkg/errors"
)

func validateFlags(ip, port, mode string, packSize, debugMode int) error {
	if net.ParseIP(ip) == nil {
		return errors.New("invalid IP address")
	}

	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return errors.New("invalid port number")
	}

	if packSize <= 0 {
		return errors.New("packet size must be greater than 0")
	}

	if debugMode < 0 || debugMode > 2 {
		return errors.New("debug mode must be 0, 1, or 2")
	}

	if mode != "grpc" && mode != "socket" {
		return errors.New("mode must be 'grpc' or 'socket'")
	}

	return nil
}

func main() {
	ip := flag.String("ip", "localhost", "client ip")
	port := flag.String("port", "50057", "client port (gRPC: 50057, Socket: 50056)")
	packSize := flag.Int("size", 512, "packet size")
	pnPackCount := flag.Int("count", 1000000, "packet count")
	pstLogTime := flag.String("logtime", "ztime.json", "logtime name")
	debugMode := flag.Int("debug", 0, "Debug mode (0: off, 1: basic, 2: verbose)")
	mode := flag.String("mode", "grpc", "Select communication mode: 'grpc' or 'socket'")
	flag.Parse()

	if err := validateFlags(*ip, *port, *mode, *packSize, *debugMode); err != nil {
		fmt.Println("[ERROR] ", err)
		return
	}

	address := fmt.Sprintf("%s:%s", *ip, *port)
	client, err := async.NewClient(*mode, address, *pstLogTime, *packSize, *pnPackCount, *debugMode)
	if err != nil {
		fmt.Println("[ERROR] ", err)
		return
	}

	if err = client.Run(); err != nil {
		fmt.Println("[FATAL] ", err)
	}
}
