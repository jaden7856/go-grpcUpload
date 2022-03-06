/*
성능 테스트 비교를 위한 예제

usage : exe -add=주소 -tls=보안 -size=패킷크기 -debug=디버그모드
 - add : ip:port
 - tls : none (평문), tls(보안)
 - size : 패킷 크기, 기본 512 byte
 - debug : 디버그 모드, 0,1,2, 기본 0

ex) exe -add=192.168.124.131:50056 -tls=tls -size=512 -debug=0

주의 : 다중연결 처리는 안 해 놨다.
*/

package main

import (
	"flag"
	"fmt"
	"net"
)

const BUFSIZE = 8388608  // 1024 * 1024 * 8
const QUESIZE = 41943040 // 1024 * 1024 * 8 * 5
const WINSIZE = 4194304  // 1024 * 1024 * 4

func main() {

	var (
		err    error
		listen net.Listener
		conn   net.Conn

		nReadSzTotal, nReadCntTotal, nSendSzTotal, nSendCntTotal = 0, 0, 0, 0
		nReadSz, nSendSz                                         = 0, 0
		nIxRead, nIxWrite, nIxCount, nIxTemp, nIxMod             = 0, 0, 0, 0, 0
		//nTemp                                                    = 0
	)

	// 초기화
	pstAddress := flag.String("add", "master:50056", "ip:port")
	pnPackSize := flag.Int("size", 512, "packet size")
	pnDebugMode := flag.Int("debug", 0, "debug mode - 0,1,2")
	flag.Parse()
	bsBufR := make([]byte, BUFSIZE)
	bsBufQ := make([]byte, QUESIZE)
	bsBufP := make([]byte, *pnPackSize)

	// 연결 대기
	listen, err = net.Listen("tcp", *pstAddress)
	checkErr(0, "Listen", err)
	fmt.Printf("Start (%s)\n", *pstAddress)

	defer listen.Close()

	// 연결 요청 처리
	for {
		conn, err = listen.Accept()
		if err != nil {
			if *pnDebugMode > 0 {
				fmt.Printf("[F] Accept (%v)\n", err)
			}
			continue
		}

		nIxCount, nIxRead, nIxWrite, nIxTemp, nIxMod = 0, 0, 0, 0, 0

		conn.(*net.TCPConn).SetReadBuffer(WINSIZE)
		conn.(*net.TCPConn).SetWriteBuffer(WINSIZE)

		//nTemp = nTemp + 1
		//fmt.Printf("Accept (%v)\n", nTemp)

		//var syncWait sync.WaitGroup
		//syncWait.Add(1)
		//
		//go func() {
		//	// 송수신
		//	defer syncWait.Done()

		for {
			// 수신
			nReadSz, err = conn.Read(bsBufR)
			if err != nil {
				if *pnDebugMode > 0 {
					fmt.Printf("[F] Read (%v)\n", err)
				}
				break
			}
			nReadCntTotal++
			nReadSzTotal += nReadSz
			if *pnDebugMode == 1 {
				fmt.Printf("Read (RSz:%d) (RCT:%d)(RST:%d)\n", nReadSz, nReadCntTotal, nReadSzTotal)
			} else if *pnDebugMode == 2 {
				fmt.Printf("Read (RSz:%d) (RCT:%d)(RST:%d)\n(%v)\n", nReadSz, nReadCntTotal, nReadSzTotal, bsBufR[:nReadSz])
			}

			// 원형버퍼큐 저장
			nIxCount += nReadSz
			if nIxCount > QUESIZE {
				fmt.Printf("[F] IxCount (%d)\n", nIxCount)
				break
			}
			nIxTemp = QUESIZE - nIxWrite
			if nIxTemp >= nReadSz {
				copy(bsBufQ[nIxWrite:], bsBufR[:nReadSz])
				nIxWrite += nReadSz
				if nIxWrite >= QUESIZE {
					nIxWrite = 0
				}
			} else {
				copy(bsBufQ[nIxWrite:], bsBufR[:nIxTemp])
				copy(bsBufQ[:], bsBufR[nIxTemp:nReadSz])
				nIxWrite = nReadSz - nIxTemp
			}

			// 원형버퍼큐 처리 -> 송신
			for {
				if nIxCount < *pnPackSize {
					if *pnDebugMode > 0 {
						// fmt.Printf("[W] Write IxCount(%d) < PACK_SIZE\n", nIxCount)
					}
					break
				}
				nIxTemp = QUESIZE - nIxRead
				if nIxTemp >= *pnPackSize {
					copy(bsBufP[:], bsBufQ[nIxRead:nIxRead+*pnPackSize])
					nIxRead += *pnPackSize
					if nIxRead >= QUESIZE {
						nIxRead = 0
					}
				} else {
					nIxMod = *pnPackSize - nIxTemp
					copy(bsBufP[:], bsBufQ[nIxRead:])
					copy(bsBufP[nIxTemp:], bsBufQ[:nIxMod])
					nIxRead = nIxMod
				}
				nIxCount -= *pnPackSize
				if *pnDebugMode > 0 {
					fmt.Printf("Pack (Cnt:%d)(Wx:%d)(Ix:%d)\n", nIxCount, nIxWrite, nIxRead)
				}

				// 송신 - 1 byte
				nSendSz, err = conn.Write(bsBufP[:1])
				if err != nil {
					fmt.Printf("[F] Write (%v)\n", err)
					break
				}
				nSendCntTotal++
				nSendSzTotal += nSendSz
				if *pnDebugMode == 1 {
					fmt.Printf("Write (SSz:%d) (SCT:%d)(SST:%d)\n", nSendSz, nSendCntTotal, nSendSzTotal)
				} else if *pnDebugMode == 2 {
					fmt.Printf("Write (SSz:%d) (SCT:%d)(SST:%d)\n(%v)", nSendSz, nSendCntTotal, nSendSzTotal, bsBufP[:1])
				}
			}
		}
		//}()

		conn.Close()
		//syncWait.Wait()
	}

	// 종료
	fmt.Println("Stop")
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
