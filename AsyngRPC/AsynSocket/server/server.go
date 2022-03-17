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

const (
	BUFSIZE = 8388608  // 1024 * 1024 * 8
	QUESIZE = 41943040 // 1024 * 1024 * 8 * 5
	WINSIZE = 4194304  // 1024 * 1024 * 4
)

func main() {
	var (
		err    error
		listen net.Listener
		conn   net.Conn

		nReadSzTotal, nReadCntTotal, nSendSzTotal, nSendCntTotal int
		nReadSz, nSendSz                                         int
		nIxRead, nIxWrite, nIxCount, nIxTemp, nIxMod             int
	)

	// 초기화
	pstAddress := flag.String("add", "master:50056", "ip:port")
	pnPackSize := flag.Int("size", 512, "packet size")
	pnDebugMode := flag.Int("debug", 0, "debug mode - 0,1,2")
	flag.Parse()

	bsBufR := make([]byte, BUFSIZE)
	bsBufP := make([]byte, *pnPackSize)

	// 연결 대기
	listen, err = net.Listen("tcp", *pstAddress)
	if err != nil {
		fmt.Printf("failed listen to %s", *pstAddress)
		panic(err)
	}
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

		if err := conn.(*net.TCPConn).SetReadBuffer(WINSIZE); err != nil {
			fmt.Printf("SetReadBuffer err %s", err)
		}
		if err = conn.(*net.TCPConn).SetWriteBuffer(WINSIZE); err != nil {
			fmt.Printf("SetWriteBuffer err %s", err)
		}

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

			switch *pnDebugMode {
			case 1:
				fmt.Printf("Read (RSz:%d) (RCT:%d)(RST:%d)\n", nReadSz, nReadCntTotal, nReadSzTotal)
			case 2:
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
				nIxWrite += nReadSz
				if nIxWrite >= QUESIZE {
					nIxWrite = 0
				}
			} else {
				nIxWrite = nReadSz - nIxTemp
			}

			// 원형버퍼큐 처리 -> 송신
			for {
				if nIxCount < *pnPackSize {
					break
				}
				nIxTemp = QUESIZE - nIxRead
				if nIxTemp >= *pnPackSize {
					nIxRead += *pnPackSize
					if nIxRead >= QUESIZE {
						nIxRead = 0
					}
				} else {
					nIxMod = *pnPackSize - nIxTemp
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
				switch *pnDebugMode {
				case 1:
					fmt.Printf("Write (SSz:%d) (SCT:%d)(SST:%d)\n", nSendSz, nSendCntTotal, nSendSzTotal)
				case 2:
					fmt.Printf("Write (SSz:%d) (SCT:%d)(SST:%d)\n(%v)", nSendSz, nSendCntTotal, nSendSzTotal, bsBufP[:1])
				}
			}
		}

		conn.Close()
		// 종료
		fmt.Println("Stop")
	}
}
