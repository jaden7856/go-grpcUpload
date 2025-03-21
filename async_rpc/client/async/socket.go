package async

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type SocketClient struct {
	wg sync.WaitGroup

	conn      net.Conn
	address   string
	packSize  int
	packCount int
	logTime   string
	debugMode int

	startTime      time.Time
	timeElapsedLog []timeElapsed
	nElapsedIx     int
}

// NewSocketClient 생성자
func NewSocketClient(address, logTime string, packSize, packCount, debugMode int) *SocketClient {
	elapsedCnt := packCount/1000 + 30
	return &SocketClient{
		wg:             sync.WaitGroup{},
		address:        address,
		packSize:       packSize,
		packCount:      packCount,
		logTime:        logTime,
		debugMode:      debugMode,
		timeElapsedLog: make([]timeElapsed, elapsedCnt),
	}
}

// Run gRPC 서버 실행 함수
func (c *SocketClient) Run() error {
	recvBuffer := make([]byte, 1)
	sendBuffer := make([]byte, c.packSize)

	c.addBuffer(sendBuffer)

	logFile, err := os.OpenFile(c.logTime, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("[FATAL] failed to open log file: %v", err)
	}

	defer func(logFile *os.File) {
		if err := logFile.Close(); err != nil {
			fmt.Println("[WARNING] Failed close file :", err)
		}
	}(logFile)

	// 시작 시간
	c.startTime = time.Now()

	// 서버 연결
	c.conn, err = net.Dial("tcp", c.address)
	if err != nil {
		return fmt.Errorf("[FATAL] Dial err : %v", err)
	}
	defer c.conn.Close()

	tcpConn, ok := c.conn.(*net.TCPConn)
	if ok {
		// tcpConn 과 conn 은 같은 메모리 주소를 가리키는 변수입니다.
		tcpConn.SetReadBuffer(WINSIZE)
		tcpConn.SetWriteBuffer(WINSIZE)
	}

	if c.debugMode > 0 {
		fmt.Printf("Start (%s)\n", c.address)
	}

	c.wg.Add(2)
	go c.receiveData(recvBuffer)
	go c.sendData(sendBuffer)
	c.wg.Wait()

	c.logElapsedTimes(logFile)

	// 종료
	if c.debugMode > 0 {
		fmt.Println("End")
	}
	return nil
}

func (c *SocketClient) addBuffer(bsBufS []byte) {
	for ix := 0; ix < c.packSize; ix++ {
		bsBufS[ix] = 'a'
	}
}

func (c *SocketClient) receiveData(recvBuffer []byte) {
	defer c.wg.Done()
	var totalReadCount, totalReadSize int

	for totalReadSize < c.packCount {
		n, err := c.conn.Read(recvBuffer)
		if err != nil {
			fmt.Printf("[FATAL] Read error: %v\n", err)
			return
		}

		totalReadCount++
		totalReadSize += n

		switch c.debugMode {
		case 1:
			fmt.Printf("Read (RSz:%d) (RCT:%d)(RST:%d)\n", n, totalReadCount, totalReadSize)
		case 2:
			fmt.Printf("Read (RSz:%d) (RCT:%d)(RST:%d)\n(%v)\n", n, totalReadCount, totalReadSize, recvBuffer[:n])
		}
	}
}

// sendData 송신
func (c *SocketClient) sendData(sendBuffer []byte) {
	defer c.wg.Done()
	var totalSentSize, totalSentCount int

	for i := 0; i < c.packCount; i++ {
		n, err := c.conn.Write(sendBuffer)
		if err != nil {
			fmt.Printf("[FATAL] Write error: %v\n", err)
			return
		}
		totalSentSize += n
		totalSentCount++

		switch c.debugMode {
		case 1:
			fmt.Printf("Write (SSz:%d) (SCT:%d)(SST:%d)\n", n, totalSentCount, totalSentSize)
		case 2:
			fmt.Printf("Write (SSz:%d) (SCT:%d)(SST:%d)\n(%v)\n", n, totalSentCount, totalSentSize, sendBuffer)
		}

		// 경과 시간 저장 (샘플링)
		c.recordElapsedTime(totalSentCount)
	}
}

// recordElapsedTime 송신 타이밍 기록
func (c *SocketClient) recordElapsedTime(sentCount int) {
	if sentCount == 1 ||
		(sentCount >= 10 && sentCount < 100 && sentCount%10 == 0) ||
		(sentCount >= 100 && sentCount < 1000 && sentCount%100 == 0) ||
		(sentCount%1000 == 0) {
		if c.nElapsedIx < len(c.timeElapsedLog) {
			c.timeElapsedLog[c.nElapsedIx] = timeElapsed{
				elapsedTime:   time.Since(c.startTime),
				nSendCntTotal: sentCount,
			}
			c.nElapsedIx++
		} else {
			fmt.Printf("[WARNING] ElapsedIx(%d) < ELASPIX\n", c.nElapsedIx)
		}
	}
}

// logElapsedTimes 타임로그 저장
func (c *SocketClient) logElapsedTimes(logFile *os.File) {
	for _, t := range c.timeElapsedLog {
		if t.nSendCntTotal > 0 {
			fmt.Fprintf(logFile, "{ \"index\" : { \"_index\" : \"commspeed\", \"_type\" : \"record\" } }\n"+
				"{\"sync\":\"async\", \"packsize\":%d, \"packcnt\":%d, \"escount\":%d, \"estime\":%v}\n",
				c.packSize, c.packCount, t.nSendCntTotal, t.elapsedTime)
		}
	}
}
