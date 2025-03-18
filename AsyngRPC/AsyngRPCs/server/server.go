/*
성능 테스트 비교를 위한 예제

usage : exe -add=주소 -tls=보안 -size=패킷크기 -debug=디버그모드
 - add : ip:port
 - tls : none (평문), tls(보안)
 - size : RPC 는 size 필요없지만 일관성을 위해
 - debug : 디버그 모드, 0,1,2, 기본 0

ex) exe -add=192.168.124.131:50057 -tls=tls -size=512 -debug=0
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/jaden7856/go-grpcUpload/AsyngRPC/AsyngRPCs/protobuf"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const WINSIZE = 4194304 // 1024 * 1024 * 4

type GrpcServer struct {
	pb.UnimplementedIpcgrpcServer
	address      string
	debugMode    int
	readCntTotal int
	sendCntTotal int
	readSzTotal  int
	sendSzTotal  int
}

// NewGrpcServer 생성자
func NewGrpcServer(address string, debugMode int) *GrpcServer {
	return &GrpcServer{
		address:   address,
		debugMode: debugMode,
	}
}

// 데이터 수신 함수
func (svr *GrpcServer) receiveData(stream pb.Ipcgrpc_SendDataServer) (*pb.IpcRequest, error) {
	req, err := stream.Recv()
	if err != nil {
		if errors.Is(err, io.EOF) {
			if svr.debugMode > 0 {
				fmt.Println("Read EOF")
			}
			return nil, io.EOF
		}
		if svr.debugMode > 0 {
			fmt.Printf("Read ERR: %v\n", err)
		}
		return nil, err
	}

	// 수신 데이터 카운트 증가
	svr.readCntTotal++
	svr.readSzTotal += int(req.Nsize)

	if svr.debugMode >= 1 {
		fmt.Printf("Read (RSz: %d) (RCT: %d) (RST: %d)\n", req.Nsize, svr.readCntTotal, svr.readSzTotal)
	}
	if svr.debugMode >= 2 {
		fmt.Printf("Data: %v\n", req.Bsreq)
	}

	return req, nil
}

// 데이터 송신 함수
func (svr *GrpcServer) sendData(stream pb.Ipcgrpc_SendDataServer, data []byte) error {
	res := &pb.IpcReply{Bsres: data}
	err := stream.Send(res)
	if err != nil {
		fmt.Printf("Send error: %v\n", err)
		return err
	}

	// 송신 데이터 카운트 증가
	svr.sendCntTotal++
	svr.sendSzTotal++

	if svr.debugMode >= 1 {
		fmt.Printf("Send (SSz: %d) (SCT: %d) (SST: %d)\n", len(data), svr.sendCntTotal, svr.sendSzTotal)
	}
	if svr.debugMode >= 2 {
		fmt.Printf("Sent Data: %v\n", res.Bsres)
	}

	return nil
}

// SendData gRPC 메서드 구현
func (svr *GrpcServer) SendData(stream pb.Ipcgrpc_SendDataServer) error {
	log.Println("Start new gRPC server session")

	for {
		req, err := svr.receiveData(stream)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// 받은 데이터 일부만 송신 (예제에서는 첫 바이트만 전송)
		if err := svr.sendData(stream, req.Bsreq[:1]); err != nil {
			return err
		}
	}

	return nil
}

// Run gRPC 서버 실행 함수
func (svr *GrpcServer) Run() error {
	listener, err := net.Listen("tcp", svr.address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", svr.address, err)
	}
	defer listener.Close()

	serverOptions := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(WINSIZE),
		grpc.MaxSendMsgSize(WINSIZE),
	}

	grpcServer := grpc.NewServer(serverOptions...)
	pb.RegisterIpcgrpcServer(grpcServer, svr)
	reflection.Register(grpcServer)

	log.Printf("Starting gRPC server on %s", svr.address)
	return grpcServer.Serve(listener)
}

func main() {
	address := flag.String("add", "localhost:50057", "gRPC server address (ip:port)")
	debugMode := flag.Int("debug", 0, "Debug mode (0: off, 1: basic, 2: verbose)")
	flag.Parse()

	// 서버 실행
	server := NewGrpcServer(*address, *debugMode)
	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
