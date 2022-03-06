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
	pb "github.com/jaden7856/go-grpcUpload/AsyngRPC/AsyngRPCs/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
)

const WINSIZE = 4194304 // 1024 * 1024 * 4

type stServer struct {
	*pb.UnimplementedIpcgrpcServer
	nReadCntTotal int
	nSendCntTotal int
	nDebugMode    int
	req           *pb.IpcRequest
	res           *pb.IpcReply
	err           error
}

func (svr *stServer) SendData(stream pb.Ipcgrpc_SendDataServer) error {
	log.Println("start new server")

	for {
		// 수신
		svr.req, svr.err = stream.Recv()
		if svr.err == io.EOF {
			if svr.nDebugMode > 0 {
				fmt.Printf("Read EOF\n\n")
			}
			break
		} else if svr.err != nil {
			if svr.nDebugMode > 0 {
				fmt.Printf("Read ERR (%v)\n\n", svr.err)
			}
			return svr.err
		} else {
			svr.nReadCntTotal++
			if svr.nDebugMode == 1 {
				fmt.Printf("Read (RSz:%d) (RCT:%d)\n", len(svr.req.Bsreq), svr.nReadCntTotal)
			} else if svr.nDebugMode == 2 {
				fmt.Printf("Read (RSz:%d) (RCT:%d)\n(%v)\n", len(svr.req.Bsreq), svr.nReadCntTotal, svr.req.Bsreq)
			}

			// 송신
			svr.res = &pb.IpcReply{Bsres: svr.req.Bsreq[0:1]}
			if svr.err = stream.Send(svr.res); svr.err != nil {
				fmt.Printf("send error %v", svr.err)
			}

			svr.nSendCntTotal++
			if svr.nDebugMode == 1 {
				fmt.Printf("Send (SSz:%d) (SCT:%d)\n", len(svr.res.Bsres), svr.nSendCntTotal)
			} else if svr.nDebugMode == 2 {
				fmt.Printf("Send (SSz:%d) (SCT:%d)\n(%v)\n", len(svr.res.Bsres), svr.nSendCntTotal, svr.res.Bsres)
			}
		}
	}
	return nil
}

func main() {
	var (
		err  error
		opts []grpc.ServerOption
	)

	pstAddress := flag.String("add", "localhost:50057", "ip:port")
	pnPackSize := flag.Int("size", 512, "packet size")
	pnDebugMode := flag.Int("debug", 0, "debug mode - 0,1,2")
	*pnPackSize = 512
	flag.Parse()

	lis, err := net.Listen("tcp", *pstAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	opts = []grpc.ServerOption{grpc.MaxRecvMsgSize(WINSIZE), grpc.MaxSendMsgSize(WINSIZE)}

	// gRPC 서버 생성
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterIpcgrpcServer(grpcServer, &stServer{nReadCntTotal: 0, nSendCntTotal: 0, nDebugMode: *pnDebugMode,
		req: nil, res: nil, err: nil})
	reflection.Register(grpcServer)

	log.Printf("start gRPC server on %s port", *pstAddress)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
