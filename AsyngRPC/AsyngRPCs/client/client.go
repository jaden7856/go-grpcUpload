/*
성능 테스트 비교를 위한 예제
rpc 처럼 1패킷 보내고 1패킷 받지만, 송수신 자체를 동기화 시키지 않는다.
서버에서는 한꺼번에 받아서 내부에서 루프 돌면서 처리
송신, 수신은 각각 고루틴으로 돈다.
단, rpc 특성상 수신하는 서버에서 패킷을 정해진 크기 외에는 한꺼번에 받지 못 한다.

usage : exe -add=주소 -tls=보안 -size=패킷크기 -count=패킷송신횟수 -loop=연결송신끊기 -logtime=타임로그이름 -debug=디버그모드
 - add : ip:port
 - tls : none (평문), tls(보안)
 - size : 패킷 크기, 기본 512 byte
 - count : 패킷 보내는 횟수, 기본 1000번
 - loop : 연결 해제 횟수, 기본 1번
 - logtime : 타임로그 이름, 기본 ztime.json
 - debug : 디버그 모드, 0,1,2, 기본 0

ex) exe -add=centos7ora2:50057 -tls=tls -size=512 -count=1000000 -loop=1 -logtime=ztime.json -debug=0
*/

package main

import (
	"flag"
	"fmt"
	pb "github.com/jaden7856/go-grpcUpload/AsyngRPC/AsyngRPCs/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"os"
	"sync"
	"time"
)

type timeElapsed struct {
	elapsedTime   time.Duration
	nSendCntTotal int
}

const WINSIZE = 4194304 // 1024 * 1024 * 4

func main() {

	// 변수 초기화
	var (
		err    error
		opts   []grpc.DialOption
		conn   *grpc.ClientConn
		client pb.IpcgrpcClient
		req    pb.IpcRequest
		res    *pb.IpcReply
		stream pb.Ipcgrpc_SendDataClient

		wgMain         sync.WaitGroup
		startTime      time.Time
		srtTimeElapeed []timeElapsed
		nElapsedCnt    = 0
		nElapsedIx     = 0

		nReadCntTotal, nSendCntTotal = 0, 0
		ix                           = 0
	)

	// 초기화
	pstAddress := flag.String("add", "master:50057", "ip:port")
	pnPackSize := flag.Int("size", 512, "packet size")
	pnPackCount := flag.Int("count", 1000000, "packet count")
	pstLogTime := flag.String("logtime", "ztime.json", "logtime name")
	pnDebugMode := flag.Int("debug", 0, "debug mode - 0,1,2")
	flag.Parse()
	bsBufS := make([]byte, *pnPackSize)
	for ix := 0; ix < *pnPackSize; ix++ {
		bsBufS[ix] = 'a'
	}
	req.Bsreq = bsBufS
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*60000)
	//defer cancel()
	nElapsedCnt = *pnPackCount/1000 + 30
	srtTimeElapeed = make([]timeElapsed, nElapsedCnt)
	fpTime, _ := os.OpenFile(*pstLogTime, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	opts = []grpc.DialOption{grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(WINSIZE), grpc.MaxCallRecvMsgSize(WINSIZE))}
	if *pnDebugMode > 0 {
		fmt.Printf("Start (%v)\n", *pstAddress)
	}

	fmt.Println("Begin")

	// 시작 시간
	startTime = time.Now()

	// Loop 개수만큼 연결 끊기
	for {
		// 서버 연결
		conn, err = grpc.Dial(*pstAddress, opts...)
		if err != nil {
			fmt.Printf("[F] Dial (%v)\n", err)
			break
		}
		client = pb.NewIpcgrpcClient(conn)

		// 스트림 생성
		//stream, err = client.SendData(ctx)
		stream, err = client.SendData(context.Background())
		if err != nil {
			fmt.Printf("[F] SendDataCTX (%v)\n", err)
			break
		}

		// 대기 등록
		wgMain.Add(1)

		// 수신
		go func() {
			defer wgMain.Done()

			for {
				res, err = stream.Recv()
				if err == io.EOF {
					if *pnDebugMode > 0 {
						fmt.Printf("Read EOF\n\n")
					}
					break
				} else if err != nil {
					fmt.Printf("Read ERR (%v)\n\n", err)
					return
				} else {
					nReadCntTotal++
					if *pnDebugMode == 1 {
						fmt.Printf("Read (RSz:%d) (RCT:%d)\n", len(res.Bsres), nReadCntTotal)
					} else if *pnDebugMode == 2 {
						fmt.Printf("Read (RSz:%d) (RCT:%d)\n(%v)\n", len(res.Bsres), nReadCntTotal, res.Bsres)
					}

					// 빠져나가기
					if nReadCntTotal >= *pnPackCount {
						//fmt.Printf("[W] Read break (%d)(%d)\n", nReadSzSum, *pnPackCount)
						break
					}
				}
			}
		}()

		// 송신
		for ix = 0; ix < *pnPackCount; ix++ {
			err = stream.Send(&req)
			if err != nil {
				//fmt.Printf("[F] Send ERR (%v)\n", err)
				break
			} else {
				nSendCntTotal++
				if *pnDebugMode == 1 {
					fmt.Printf("SendData (SSz:%d) (SCT:%d)\n", len(req.Bsreq), nSendCntTotal)
				} else if *pnDebugMode == 2 {
					fmt.Printf("SendData (SSz:%d) (SCT:%d)\n(%v)\n", len(req.Bsreq), nSendCntTotal, req.Bsreq)
				}

				// 경과 시간 저장
				if nSendCntTotal == 1 ||
					(nSendCntTotal >= 10 && nSendCntTotal < 100 && nSendCntTotal%10 == 0) ||
					(nSendCntTotal >= 100 && nSendCntTotal < 1000 && nSendCntTotal%100 == 0) ||
					(nSendCntTotal%1000 == 0) {
					if nElapsedIx < nElapsedCnt {
						srtTimeElapeed[nElapsedIx].nSendCntTotal = nSendCntTotal
						srtTimeElapeed[nElapsedIx].elapsedTime = time.Since(startTime)
						nElapsedIx++
					} else {
						fmt.Printf("[W] ElapsedIx(%d) < ELASPIX\n", nElapsedIx)
					}
				}
			}
		}

		// 대기
		if *pnDebugMode > 0 {
			fmt.Println("Wait")
		}
		wgMain.Wait()

		// 종료
		if *pnDebugMode > 0 {
			fmt.Println("Stop")
		}

		conn.Close()
	}

	// 타임로그 작성
	if nElapsedIx < nElapsedCnt {
		srtTimeElapeed[nElapsedIx].nSendCntTotal = nSendCntTotal + 1
		srtTimeElapeed[nElapsedIx].elapsedTime = time.Since(startTime)
		nElapsedIx++
	}
	for ix = 0; ix < nElapsedCnt; ix++ {
		if srtTimeElapeed[ix].nSendCntTotal > 0 {
			fmt.Fprintf(fpTime, "{ \"index\" : { \"_index\" : \"commspeed\", \"_type\" : \"record\", \"_id\" : \"%v\" } }\n{\"sync\":\"async\", \"packsize\":%d, \"packcnt\":%d, \"escount\":%d, \"estime\":%v}\n",
				time.Now().UnixNano(), *pnPackSize, *pnPackCount, srtTimeElapeed[ix].nSendCntTotal, srtTimeElapeed[ix].elapsedTime)
		}
	}

	// 종료
	_ = fpTime.Close()
	fmt.Println("End")
}