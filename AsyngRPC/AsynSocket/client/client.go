/*
성능 테스트 비교를 위한 예제
rpc 처럼 1패킷 보내고 1패킷 받지만, 송수신 자체를 동기화 시키지 않는다.
서버에서는 한꺼번에 받아서 내부에서 루프 돌면서 처리
송신, 수신은 각각 고루틴으로 돈다.

usage : exe -add=주소 -tls=보안 -size=패킷크기 -count=패킷송신횟수 -loop=연결송신끊기 -logtime=타임로그이름 -debug=디버그모드
 - add : ip:port
 - tls : none (평문), tls(보안)
 - size : 패킷 크기, 기본 512 byte
 - count : 패킷 보내는 횟수, 기본 1000번
 - loop : 연결 해제 횟수, 기본 1번
 - logtime : 타임로그 이름, 기본 ztime.json
 - debug : 디버그 모드, 0,1,2, 기본 0

ex) exe -add=192.168.124.131:50056 -tls=tls -size=512 -count=1000000 -loop=1 -logtime=ztime.json -debug=0
*/

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
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

	var (
		err    error
		client net.Conn
		cert   tls.Certificate
		conf   *tls.Config

		wgMain         sync.WaitGroup
		startTime      time.Time
		srtTimeElapeed []timeElapsed
		nElapsedCnt    int = 0
		nElapsedIx     int = 0

		nReadSzTotal, nReadCntTotal, nSendSzTotal, nSendCntTotal int = 0, 0, 0, 0
		nReadSz, nSendSz, ix                                     int = 0, 0, 0
	)

	// 초기화
	pstAddress := flag.String("add", "master:50056", "ip:port")
	pstTls := flag.String("tls", "none", "none or tls")
	pnPackSize := flag.Int("size", 512, "packet size")
	pnPackCount := flag.Int("count", 1000000, "packet count")
	pnLoopCount := flag.Int("loop", 1, "loop count")
	pstLogTime := flag.String("logtime", "ztime.json", "logtime name")
	pnDebugMode := flag.Int("debug", 0, "debug mode - 0,1,2")
	flag.Parse()
	bsBufR := make([]byte, 1)
	bsBufS := make([]byte, *pnPackSize)
	for ix := 0; ix < *pnPackSize; ix++ {
		bsBufS[ix] = 'a'
	}
	nElapsedCnt = *pnPackCount**pnLoopCount/1000 + 30
	srtTimeElapeed = make([]timeElapsed, nElapsedCnt)
	fpTime, _ := os.OpenFile(*pstLogTime, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if *pstTls == "tls" {
		cert, err = tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
		checkErr(0, "LoadX509KeyPair", err)
		conf = &tls.Config{
			Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true,
		}
	}
	fmt.Println("Begin")

	// 시작 시간
	startTime = time.Now()

	// Loop 개수만큼 연결 끊기
	for nLoop := 0; nLoop < *pnLoopCount; nLoop++ {
		// 서버 연결
		if *pstTls == "tls" {
			client, err = tls.Dial("tcp", *pstAddress, conf)
			if err != nil {
				fmt.Printf("[F] Dial (%v)\n", err)
				break
			}
			//client.(*net.TCPConn).SetReadBuffer(WINSIZE)
			//client.(*net.TCPConn).SetWriteBuffer(WINSIZE)
			if *pnDebugMode > 0 {
				fmt.Printf("Start TLS (%s)\n", *pstAddress)
			}
		} else {
			client, err = net.Dial("tcp", *pstAddress)
			if err != nil {
				fmt.Printf("[F] Dial (%v)\n", err)
				break
			}
			client.(*net.TCPConn).SetReadBuffer(WINSIZE)
			client.(*net.TCPConn).SetWriteBuffer(WINSIZE)
			if *pnDebugMode > 0 {
				fmt.Printf("Start (%s)\n", *pstAddress)
			}

			fmt.Printf("Connect (%v)\n", nLoop) //[swc-test]
		}
		// client.SetDeadline(time.Now().Add(time.Second * 1))
		time.Sleep(time.Second * 1) //[swc-zzz]

		// 대기 등록
		wgMain.Add(1)

		// 수신
		go func(cli net.Conn) {
			defer wgMain.Done()

			for {
				nReadSz, err = client.Read(bsBufR)
				if err != nil {
					fmt.Printf("[F] Read (%v)\n", err)
					break
				}
				nReadCntTotal++
				nReadSzTotal += nReadSz
				if *pnDebugMode == 1 {
					fmt.Printf("Read (RSz:%d) (RCT:%d)(RST:%d)\n", nReadSz, nReadCntTotal, nReadSzTotal)
				} else if *pnDebugMode == 2 {
					fmt.Printf("Read (RSz:%d) (RCT:%d)(RST:%d)\n(%v)\n", nReadSz, nReadCntTotal, nReadSzTotal, bsBufR[:nReadSz])
				}

				// 빠져나가기
				if nReadSzTotal >= *pnPackCount {
					//fmt.Printf("[W] Read break (%d)(%d)\n", nReadSzTotal, *pnPackCount)
					break
				}
			}

			if *pnLoopCount <= 1 {
				fmt.Println("Stop Read")
			}
		}(client)

		// 송신
		go func(cli net.Conn) {
			for ix = 0; ix < *pnPackCount; ix++ {
				// 송신
				nSendSz, err = client.Write(bsBufS)
				if err != nil {
					fmt.Printf("[F] Write (%v)\n", err)
					break
				}
				nSendCntTotal++
				nSendSzTotal += nSendSz
				if *pnDebugMode == 1 {
					fmt.Printf("Write (SSz:%d) (SCT:%d)(SST:%d)\n", nSendSz, nSendCntTotal, nSendSzTotal)
				} else if *pnDebugMode == 2 {
					fmt.Printf("Write (SSz:%d) (SCT:%d)(SST:%d)\n(%v)\n", nSendSz, nSendCntTotal, nSendSzTotal, bsBufS)
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

			// 대기
			if *pnDebugMode > 0 {
				fmt.Println("Wait")
			}
			wgMain.Wait()

			// 종료
			if *pnDebugMode > 0 {
				fmt.Println("Stop")
			}
			client.Close()
		}(client)
	}

	// 타임로그 작성
	if nElapsedIx < nElapsedCnt {
		srtTimeElapeed[nElapsedIx].nSendCntTotal = nSendCntTotal + 1
		srtTimeElapeed[nElapsedIx].elapsedTime = time.Since(startTime)
		nElapsedIx++
	}
	for ix = 0; ix < nElapsedCnt; ix++ {
		if srtTimeElapeed[ix].nSendCntTotal > 0 {
			fmt.Fprintf(fpTime, "{ \"index\" : { \"_index\" : \"commspeed\", \"_type\" : \"record\", \"_id\" : \"%v\" } }\n{\"sender\":\"AsynTCP_%s\", \"tls\":\"%s\", \"sync\":\"async\", \"packsize\":%d, \"loopcnt\":%d, \"packcnt\":%d, \"escount\":%d, \"estime\":%v}\n",
				time.Now().UnixNano(), *pstTls, *pstTls, *pnPackSize, *pnLoopCount, *pnPackCount, srtTimeElapeed[ix].nSendCntTotal, srtTimeElapeed[ix].elapsedTime.Nanoseconds())
		}
	}

	// 종료
	_ = fpTime.Close()
	fmt.Println("End")
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
