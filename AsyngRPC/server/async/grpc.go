package async

import (
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/jaden7856/go-grpcUpload/AsyngRPC/protobuf"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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
		if err := svr.sendData(stream, req.Payload[:1]); err != nil {
			return err
		}
	}

	return nil
}

// 데이터 수신 함수
func (svr *GrpcServer) receiveData(stream pb.Ipcgrpc_SendDataServer) (*pb.IpcRequest, error) {
	req, err := stream.Recv()
	if err != nil {
		if errors.Is(err, io.EOF) {
			if svr.debugMode > 0 {
				fmt.Println("[FATAL] Read EOF")
			}
			return nil, io.EOF
		}
		if svr.debugMode > 0 {
			fmt.Println("[FATAL] Read error:", err)
		}
		return nil, err
	}

	// 수신 데이터 카운트 증가
	svr.readCntTotal++
	svr.readSzTotal += int(req.Size)

	if svr.debugMode >= 1 {
		fmt.Printf("Read (RSz: %d) (RCT: %d) (RST: %d)\n", req.Size, svr.readCntTotal, svr.readSzTotal)
	}
	if svr.debugMode >= 2 {
		fmt.Println("Data:", req.Payload)
	}

	return req, nil
}

// 데이터 송신 함수
func (svr *GrpcServer) sendData(stream pb.Ipcgrpc_SendDataServer, data []byte) error {
	res := &pb.IpcReply{Payload: data}
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
		fmt.Println("Sent Data:", res.Payload)
	}

	return nil
}
