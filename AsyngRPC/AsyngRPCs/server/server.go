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
	"net"
)

//
type server struct {
	pb.UnimplementedIpcgrpcServer
	//server *grpc.Server

	nReadCntTotal int
	nSendCntTotal int
	nDebugMode    int
	req           *pb.IpcRequest
	res           *pb.IpcReply
	err           error
}

const WINSIZE = 4194304 // 1024 * 1024 * 4

func main() {

	var (
		err  error
		opts []grpc.ServerOption
	)

	// 초기화
	pstAddress := flag.String("add", "master:50057", "ip:port")
	pnPackSize := flag.Int("size", 512, "packet size")
	*pnPackSize = 512
	pnDebugMode := flag.Int("debug", 0, "debug mode - 0,1,2")
	flag.Parse()

	// 연결 대기
	listen, err := net.Listen("tcp", *pstAddress)
	checkErr(0, "Listen", err)

	opts = []grpc.ServerOption{grpc.MaxRecvMsgSize(WINSIZE), grpc.MaxSendMsgSize(WINSIZE)}
	fmt.Printf("Start (%v)\n", *pstAddress)

	// 호출 대기
	s := grpc.NewServer(opts...)
	pb.RegisterIpcgrpcServer(s, &server{nDebugMode: *pnDebugMode})
	reflection.Register(s)
	err = s.Serve(listen)
	if err != nil {
		fmt.Printf("[F] Serve\n")
	}

	// 종료
	listen.Close()
	fmt.Println("Stop")
}

//
func (svr *server) SendData(stream pb.Ipcgrpc_SendDataServer) error {
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
			svr.err = stream.Send(svr.res)
			checkErr(0, "Send", svr.err)
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

// 에러 체크
func checkErr(nMode int, strTitle string, err error) {
	if err != nil {
		fmt.Println(strTitle + " -> " + err.Error())
		if nMode == 0 {
			panic(err)
		} else {
		}
	} else {
		// 에러 아니여 ~
	}
}
