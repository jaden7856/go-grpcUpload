package async

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	pb "github.com/jaden7856/go-grpcUpload/AsyngRPC/protobuf"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	wg sync.WaitGroup

	opts   []grpc.DialOption
	conn   *grpc.ClientConn
	stream pb.Ipcgrpc_SendDataClient

	address   string
	packSize  int
	packCount int
	logTime   string
	debugMode int

	startTime      time.Time
	timeElapsedLog []timeElapsed
	nElapsedIx     int
}

// NewGrpcClient 생성자
func NewGrpcClient(address, logTime string, packSize, packCount, debugMode int) *GrpcClient {
	elapsedCnt := packCount/1000 + 30
	return &GrpcClient{
		wg:             sync.WaitGroup{},
		opts:           []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(WINSIZE), grpc.MaxCallRecvMsgSize(WINSIZE))},
		address:        address,
		packSize:       packSize,
		packCount:      packCount,
		logTime:        logTime,
		debugMode:      debugMode,
		timeElapsedLog: make([]timeElapsed, elapsedCnt),
	}
}

// Run gRPC 서버 실행 함수
func (c *GrpcClient) Run() error {
	bsBufS := make([]byte, c.packSize)
	c.addBuffer(bsBufS)

	var req pb.IpcRequest
	req.Payload = bsBufS
	req.Size = int64(c.packSize)

	logFile, err := os.OpenFile(c.logTime, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("[FATAL] failed to open log file: %v", err)
	}

	defer func(logFile *os.File) {
		if err := logFile.Close(); err != nil {
			fmt.Println("[WARNING] Failed close file :", err)
		}
	}(logFile)

	if c.debugMode > 0 {
		fmt.Printf("Start (%s)\n", c.address)
	}

	// 시작 시간
	c.startTime = time.Now()

	// 서버 연결
	c.conn, err = grpc.NewClient(c.address, c.opts...)
	if err != nil {
		return fmt.Errorf("[FATAL] New Client error : %v\n", err)
	}
	defer c.conn.Close()

	client := pb.NewIpcgrpcClient(c.conn)

	// 스트림 생성
	c.stream, err = client.SendData(context.Background())
	if err != nil {
		return fmt.Errorf("[FATAL] SendDataCTX error : %v\n", err)
	}

	c.wg.Add(2)
	go c.receiveData()
	go c.sendData(&req)
	c.wg.Wait()

	c.logElapsedTimes(logFile)

	// 종료
	if c.debugMode > 0 {
		fmt.Println("End")
	}
	return nil
}

func (c *GrpcClient) addBuffer(bsBufS []byte) {
	for ix := 0; ix < c.packSize; ix++ {
		bsBufS[ix] = 'a'
	}
}

func (c *GrpcClient) receiveData() {
	defer c.wg.Done()
	var totalReadCount int

	for {
		res, err := c.stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if c.debugMode > 0 {
					fmt.Printf("Read EOF\n\n")
					break
				}
			}
			fmt.Printf("Read ERR (%v)\n\n", err)
			return
		}

		totalReadCount++
		switch c.debugMode {
		case 1:
			fmt.Printf("Read (RSz:%d) (RCT:%d)\n", len(res.Payload), totalReadCount)
		case 2:
			fmt.Printf("Read (RSz:%d) (RCT:%d)\n(%v)\n", len(res.Payload), totalReadCount, res.Payload)
		}

		// 빠져나가기
		/*
			if nReadCntTotal >= *pnPackCount {
				fmt.Printf("[W] Read break (%d)(%d)\n", nReadSzSum, *pnPackCount)
				break
			}
		*/
	}
}

// sendData 송신
func (c *GrpcClient) sendData(req *pb.IpcRequest) {
	defer c.wg.Done()
	var totalSentCount int

	for i := 0; i < c.packCount; i++ {
		if err := c.stream.Send(req); err != nil {
			fmt.Printf("[FATAL] Send err (%v)\n", err)
			break
		}

		totalSentCount++

		switch c.debugMode {
		case 1:
			fmt.Printf("SendData (SSz:%d) (SCT:%d)\n", len(req.Payload), totalSentCount)
		case 2:
			fmt.Printf("SendData (SSz:%d) (SCT:%d)\n(%v)\n", len(req.Payload), totalSentCount, req.Payload)
		}

		// 경과 시간 저장
		c.recordElapsedTime(totalSentCount)
	}
}

// recordElapsedTime 송신 타이밍 기록
func (c *GrpcClient) recordElapsedTime(sentCount int) {
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
func (c *GrpcClient) logElapsedTimes(logFile *os.File) {
	for _, t := range c.timeElapsedLog {
		if t.nSendCntTotal > 0 {
			fmt.Fprintf(logFile, "{ \"index\" : { \"_index\" : \"commspeed\", \"_type\" : \"record\" } }\n"+
				"{\"sync\":\"async\", \"packsize\":%d, \"packcnt\":%d, \"escount\":%d, \"estime\":%v}\n",
				c.packSize, c.packCount, t.nSendCntTotal, t.elapsedTime)
		}
	}
}
