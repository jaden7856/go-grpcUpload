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

		wgMain         sync.WaitGroup
		startTime      time.Time
		srtTimeElapsed []timeElapsed
		nElapsedCnt    = 0
		nElapsedIx     = 0

		nReadSzTotal, nReadCntTotal, nSendSzTotal, nSendCntTotal = 0, 0, 0, 0
		nReadSz, nSendSz, ix                                     = 0, 0, 0
	)

	// 초기화
	pstAddress := flag.String("add", "master:50056", "ip:port")
	pnPackSize := flag.Int("size", 512, "packet size")
	pnPackCount := flag.Int("count", 1000000, "packet count")
	pstLogTime := flag.String("logtime", "ztime.json", "logtime name")
	pnDebugMode := flag.Int("debug", 0, "debug mode - 0,1,2")
	flag.Parse()

	bsBufR := make([]byte, 1)
	bsBufS := make([]byte, *pnPackSize)

	for ix := 0; ix < *pnPackSize; ix++ {
		bsBufS[ix] = 'a'
	}

	nElapsedCnt = *pnPackCount/1000 + 30
	srtTimeElapsed = make([]timeElapsed, nElapsedCnt)

	fpTime, _ := os.OpenFile(*pstLogTime, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer func(fpTime *os.File) {
		if err := fpTime.Close(); err != nil {
			fmt.Println("Failed open file")
		}
	}(fpTime)

	// 시작 시간
	startTime = time.Now()

	// 서버 연결
	client, err = net.Dial("tcp", *pstAddress)
	if err != nil {
		fmt.Printf("[F] Dial (%v)\n", err)
	}

	if err = client.(*net.TCPConn).SetReadBuffer(WINSIZE); err != nil {
		fmt.Printf("SetReadBuffer err %s", err)
	}
	if err = client.(*net.TCPConn).SetWriteBuffer(WINSIZE); err != nil {
		fmt.Printf("SetWriteBuffer err %s", err)
	}
	if *pnDebugMode > 0 {
		fmt.Printf("Start (%s)\n", *pstAddress)
	}

	// 대기 등록
	wgMain.Add(2)

	// 수신
	go func() {
		defer wgMain.Done()

		for {
			nReadSz, err = client.Read(bsBufR)
			if err != nil {
				fmt.Printf("[F] Read (%v)\n", err)
				break
			}
			nReadCntTotal++
			nReadSzTotal += nReadSz

			switch *pnDebugMode {
			case 1:
				fmt.Printf("Read (RSz:%d) (RCT:%d)(RST:%d)\n", nReadSz, nReadCntTotal, nReadSzTotal)
			case 2:
				fmt.Printf("Read (RSz:%d) (RCT:%d)(RST:%d)\n(%v)\n", nReadSz, nReadCntTotal, nReadSzTotal, bsBufR[:nReadSz])
			}

			// 빠져나가기
			if nReadSzTotal >= *pnPackCount {
				break
			}
		}
	}()

	// 송신
	go func() {
		defer wgMain.Done()

		for ix = 0; ix < *pnPackCount; ix++ {
			// 송신
			nSendSz, err = client.Write(bsBufS)
			if err != nil {
				fmt.Printf("[F] Write (%v)\n", err)
				break
			}
			nSendCntTotal++
			nSendSzTotal += nSendSz

			switch *pnDebugMode {
			case 1:
				fmt.Printf("Write (SSz:%d) (SCT:%d)(SST:%d)\n", nSendSz, nSendCntTotal, nSendSzTotal)
			case 2:
				fmt.Printf("Write (SSz:%d) (SCT:%d)(SST:%d)\n(%v)\n", nSendSz, nSendCntTotal, nSendSzTotal, bsBufS)
			}

			// 경과 시간 저장
			if nSendCntTotal == 1 ||
				(nSendCntTotal >= 10 && nSendCntTotal < 100 && nSendCntTotal%10 == 0) ||
				(nSendCntTotal >= 100 && nSendCntTotal < 1000 && nSendCntTotal%100 == 0) ||
				(nSendCntTotal%1000 == 0) {
				if nElapsedIx < nElapsedCnt {
					srtTimeElapsed[nElapsedIx].nSendCntTotal = nSendCntTotal
					srtTimeElapsed[nElapsedIx].elapsedTime = time.Since(startTime)
					nElapsedIx++
				} else {
					fmt.Printf("[W] ElapsedIx(%d) < ELASPIX\n", nElapsedIx)
				}
			}
		}
	}()

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

	// 타임로그 작성
	if nElapsedIx < nElapsedCnt {
		srtTimeElapsed[nElapsedIx].nSendCntTotal = nSendCntTotal + 1
		srtTimeElapsed[nElapsedIx].elapsedTime = time.Since(startTime)
		nElapsedIx++
	}
	for ix = 0; ix < nElapsedCnt; ix++ {
		if srtTimeElapsed[ix].nSendCntTotal > 0 {
			fmt.Fprintf(fpTime, "{ \"index\" : { \"_index\" : \"commspeed\", \"_type\" : \"record\", \"_id\" : \"%v\" } }\n{\"sync\":\"async\", \"packsize\":%d, \"packcnt\":%d, \"escount\":%d, \"estime\":%v}\n",
				time.Now().UnixNano(), *pnPackSize, *pnPackCount, srtTimeElapsed[ix].nSendCntTotal, srtTimeElapsed[ix].elapsedTime)
		}
	}

	// 종료
	fmt.Println("End")
}
