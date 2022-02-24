package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
	streamPb "github.com/jaden7856/go-grpcUpload/streamProtoc"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
)

type timeElapsed struct {
	elapsedTime   time.Duration
	nSendCntTotal int
}

const chunkSize = 64 * 1024

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // 전송이 아직 살아있는지 확인하기 위해 서버에 ping 을 보냅니다.
	Timeout:             time.Second,      // 연결이 끊어진 것으로 간주하기 전에 ping ack 을 1초 동안 기다립니다.
	PermitWithoutStream: true,             // 활성 스트림 없이도 ping 보내기
}

type uploader struct {
	dir     string
	packCnt int
	loopCnt int
	client  streamPb.UploadFileServiceClient
	ctx     context.Context
	wg      sync.WaitGroup
	pool    *pb.Pool

	requests    chan string // 각 요청은 클라이언트가 액세스할 수 있는 클라이언트의 파일 경로입니다.
	DoneRequest chan string
	FailRequest chan string
}

func NewUploader(ctx context.Context, client streamPb.UploadFileServiceClient, dir string, packCnt int, loopCnt int) *uploader {
	d := &uploader{
		ctx:         ctx,
		client:      client,
		dir:         dir,
		packCnt:     packCnt,
		loopCnt:     loopCnt,
		requests:    make(chan string),
		DoneRequest: make(chan string),
		FailRequest: make(chan string),
	}
	// 한번에 몇개의 파일을 보낼지
	for i := 0; i < 5; i++ {
		d.wg.Add(1)
		go d.worker(i + 1)
	}

	d.pool, _ = pb.StartPool()
	return d
}

func (d *uploader) Do(filepath string) {
	// 파일 경로를 requests 에 전달
	d.requests <- filepath
}

func (d *uploader) Stop() {
	close(d.requests)
	d.wg.Wait()
	d.pool.RefreshRate = 500 * time.Millisecond
	d.pool.Stop()
}

//goland:noinspection ALL
func (d *uploader) worker(workerID int) {
	defer d.wg.Done()
	var (
		buf        []byte
		firstChunk bool

		startTime      time.Time
		srtTimeElapsed []timeElapsed
		nElapsedCnt    int = 0
		nElapsedIx     int = 0

		nSendCntTotal int = 0
		ix            int = 0
	)

	nElapsedCnt = d.packCnt*d.loopCnt/1000 + 30
	srtTimeElapsed = make([]timeElapsed, nElapsedCnt)

	// 시작 시간
	startTime = time.Now()

	// 파일 경로에서 파일들을 추출
	for request := range d.requests {
		file, errOpen := os.Open(request)
		if errOpen != nil {
			errOpen = errors.Wrapf(errOpen,
				"failed to open file %s",
				request)
			return
		}
		defer file.Close()

		//start uploader
		stream, err := d.client.Upload(d.ctx)
		if err != nil {
			err = errors.Wrapf(err,
				"failed to create upload stream for file %s",
				request)
			return
		}
		defer stream.CloseSend()

		// file info
		stat, errstat := file.Stat()
		if errstat != nil {
			err = errors.Wrapf(err,
				"Unable to get file size  %s",
				request)
			return
		}

		// client 측에서 server에 send할때 진행률 bar 표시
		bar := pb.New64(stat.Size()).Postfix(" " + filepath.Base(request))
		bar.Units = pb.U_BYTES
		d.pool.Add(bar)

		// 파일 크기만큼 바이트 슬라이스 생성
		buf = make([]byte, chunkSize)
		firstChunk = true
		for ix = 0; ix < d.packCnt; ix++ {
			// 파일의 내용을 읽어서 바이트 슬라이스에 저장
			f, errRead := file.Read(buf)
			if errRead != nil {
				if errRead == io.EOF {
					errRead = nil
					break
				}
				errRead = errors.Wrapf(err, "errored while copying from file to buf")
				return
			}

			if firstChunk {
				nSendCntTotal++
				// Send to Server
				err = stream.Send(&streamPb.UploadRequest{
					Content:  buf[:f],
					Filename: request,
				})
				firstChunk = false
			} else {
				nSendCntTotal++
				err = stream.Send(&streamPb.UploadRequest{
					Content: buf[:f],
				})
			}
			if err != nil {
				bar.FinishPrint("failed to send chunk via stream file : " + request)
				break
			}
			bar.Add64(int64(f))

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

		// 클라이언트에서 send를 완료하고 서버에서 다 받고나서 완료 메세지를 보내면
		// 그 메세지를 받아서 상태를 체크
		status, err := stream.CloseAndRecv()
		if err != nil { //retry needed
			fmt.Println("failed to receive upstream status response")
			bar.FinishPrint("Error uploading file : " + request + " Error :" + err.Error())
			bar.Reset(0)
			d.FailRequest <- request
			return
		}

		if status.Code != streamPb.UploadStatusCode_Ok { //retry needed
			fmt.Printf("Error uploading file : " + request + " :" + status.Message)
			bar.Reset(0)
			d.FailRequest <- request
			return
		}

		// 타임로그 작성
		if nElapsedIx < nElapsedCnt {
			srtTimeElapsed[nElapsedIx].nSendCntTotal = nSendCntTotal + 1
			srtTimeElapsed[nElapsedIx].elapsedTime = time.Since(startTime)
			nElapsedIx++
		}
		for ix = 0; ix < nElapsedCnt; ix++ {
			if srtTimeElapsed[ix].nSendCntTotal > 0 {
				fmt.Fprintf(file, "{ \"index\" : { \"_index\" : \"commspeed\", \"_type\" : \"record\", \"_id\" : \"%v\" } }\n{\"sender\":\"AsynGRPCst_%s\", \"packcnt\":%d, \"loopcnt\":%d, \"escount\":%d, \"estime\":%v}\n",
					time.Now().UnixNano(), d.packCnt, d.loopCnt, srtTimeElapsed[ix].nSendCntTotal, srtTimeElapsed[ix].elapsedTime.Nanoseconds())
			}
		}

		d.DoneRequest <- request
		bar.Finish()
	}
}

func UploadFile(ctx context.Context, client streamPb.UploadFileServiceClient, filePathList []string, dir string, packCnt int, loopCnt int) error {
	d := NewUploader(ctx, client, dir, packCnt, loopCnt)
	var errorUploadBulk error

	if dir != "" {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}
		defer d.Stop()

		go func() {
			for _, file := range files {
				// 파일인지 아닌지 확인
				if !file.IsDir() {
					d.Do(dir + "/" + file.Name())
				}
			}
		}()

		for _, file := range files {
			if !file.IsDir() {
				select {
				case <-d.DoneRequest:

				case req := <-d.FailRequest:
					fmt.Println("failed to  send " + req)
					errorUploadBulk = errors.Wrapf(errorUploadBulk, " Failed to send %s", req)

				}
			}
		}
		fmt.Println("All done ")

	} else {
		go func() {
			for _, file := range filePathList {
				d.Do(file)
			}
		}()
		defer d.Stop()

		for i := 0; i < len(filePathList); i++ {
			select {
			case <-d.DoneRequest:
			case req := <-d.FailRequest:
				fmt.Println("failed to send " + req)
				errorUploadBulk = errors.Wrapf(errorUploadBulk, " Failed to send %s", req)

			}
		}
	}

	return errorUploadBulk
}

//goland:noinspection ALL
func uploadCommand() cli.Command {
	return cli.Command{
		Name:  "upload-test",
		Usage: "Uploads files to server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "a",
				Value: "localhost:port",
				Usage: "Server address",
			},
			&cli.StringFlag{
				Name:  "d",
				Value: "",
				Usage: "base directory",
			},
			&cli.IntFlag{
				Name:  "c",
				Value: 1000000,
				Usage: "packet count",
			},
			&cli.IntFlag{
				Name:  "l",
				Value: 1,
				Usage: "loop count",
			},
		},
		Action: func(c *cli.Context) error {
			var options = []grpc.DialOption{
				grpc.WithInsecure(),
				// 연결이 작동될 때까지 Dial 호출자가 차단. 이것이 없으면 Dial은 즉시 반환되고 서버 연결은 백그라운드에서 발생합니다.
				grpc.WithBlock(),
				// 클라이언트 연결이 마지막으로 발생한 연결 오류와 context.DeadlineExceeded 오류를 모두 포함하는 문자열을 반환
				grpc.WithReturnConnectionError(),
				// 클라이언트와 서버가 연결이 되어있는지 지속적으로 확인 watchdog
				grpc.WithKeepaliveParams(kacp),
			}
			addr := c.String("a")
			conn, err := grpc.Dial(addr, options...)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
					fmt.Println("faild conn close")
				}
			}(conn)

			resp, err := healthpb.NewHealthClient(conn).Check(context.Background(),
				&healthpb.HealthCheckRequest{
					Service: "test",
				})

			if err != nil {
				log.Printf("can't connect grpc server: %v", err)
			} else {
				log.Printf("status: %s", resp.GetStatus().String())
			}

			return UploadFile(
				context.Background(),
				streamPb.NewUploadFileServiceClient(conn),
				[]string{},
				c.String("d"),
				c.Int("c"),
				c.Int("l"),
			)
		},
	}
}