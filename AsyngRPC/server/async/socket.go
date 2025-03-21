package async

import (
	"fmt"
	"log"
	"net"

	"github.com/jaden7856/go-grpcUpload/AsyngRPC/protobuf"
)

type SocketServer struct {
	pb.UnimplementedIpcgrpcServer

	conn      net.Conn
	address   string
	debugMode int
	packSize  int

	readCntTotal int
	sendCntTotal int
	readSzTotal  int
	sendSzTotal  int
}

// NewSocketServer 생성자
func NewSocketServer(address string, debugMode, packSize int) *SocketServer {
	return &SocketServer{
		address:   address,
		debugMode: debugMode,
		packSize:  packSize,
	}
}

// Run gRPC 서버 실행 함수
func (svr *SocketServer) Run() error {
	listener, err := net.Listen("tcp", svr.address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", svr.address, err)
	}
	defer listener.Close()

	log.Printf("SocketServer started on %s", svr.address)

	// 연결 요청 처리
	for {
		svr.conn, err = listener.Accept()
		if err != nil {
			if svr.debugMode > 0 {
				fmt.Println("[FATAL] Accept:", err)
			}
			continue
		}

		go svr.handleConnection()
	}
}

func (svr *SocketServer) handleConnection() {
	defer func() {
		svr.conn.Close()
		fmt.Println("Stop")
	}()

	// 설정 버퍼 크기 조정
	tcpConn, ok := svr.conn.(*net.TCPConn)
	if ok {
		// tcpConn 과 conn 은 같은 메모리 주소를 가리키는 변수입니다.
		tcpConn.SetReadBuffer(WINSIZE)
		tcpConn.SetWriteBuffer(WINSIZE)
	}

	bsBufR := make([]byte, BUFSIZE)
	bsBufP := make([]byte, svr.packSize)

	for {
		nReadSz, err := svr.conn.Read(bsBufR)
		if err != nil {
			if svr.debugMode > 0 {
				fmt.Println("[FATAL] Read error:", err)
			}
			return
		}

		svr.readCntTotal++
		svr.readSzTotal += nReadSz

		switch svr.debugMode {
		case 1:
			fmt.Printf("Read (RSz:%d) (RCT:%d)(RST:%d)\n", nReadSz, svr.readCntTotal, svr.readSzTotal)
		case 2:
			fmt.Printf("Read (RSz:%d) (RCT:%d)(RST:%d)\n(%v)\n", nReadSz, svr.readCntTotal, svr.readSzTotal, bsBufR[:nReadSz])
		}

		svr.processReceivedData(bsBufP)
	}
}

func (svr *SocketServer) processReceivedData(bsBufP []byte) {
	var nIxRead, nIxWrite, nIxCount, nIxTemp, nIxMod int

	// 패킷 단위로 데이터 처리
	for {
		if nIxCount < svr.packSize {
			break
		}
		nIxTemp = QUESIZE - nIxRead
		if nIxTemp >= svr.packSize {
			nIxRead += svr.packSize
			if nIxRead >= QUESIZE {
				nIxRead = 0
			}
		} else {
			nIxMod = svr.packSize - nIxTemp
			nIxRead = nIxMod
		}
		nIxCount -= svr.packSize
		if svr.debugMode > 0 {
			fmt.Printf("Pack (Cnt:%d)(Wx:%d)(Ix:%d)\n", nIxCount, nIxWrite, nIxRead)
		}

		// 송신 (1 byte)
		nSendSz, err := svr.conn.Write(bsBufP[:1])
		if err != nil {
			fmt.Println("[FATAL] Write:", err)
			return
		}

		svr.sendCntTotal++
		svr.sendSzTotal += nSendSz

		switch svr.debugMode {
		case 1:
			fmt.Printf("Write (SSz:%d) (SCT:%d)(SST:%d)\n", nSendSz, svr.sendCntTotal, svr.sendSzTotal)
		case 2:
			fmt.Printf("Write (SSz:%d) (SCT:%d)(SST:%d)\n(%v)\n", nSendSz, svr.sendCntTotal, svr.sendSzTotal, bsBufP[:1])
		}
	}
}
