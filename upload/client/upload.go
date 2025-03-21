package main

import (
	"context"
	"fmt"
	"github.com/jaden7856/go-grpcUpload/upload/streamProtoc"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const chunkSize = 64 * 1024

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // 전송이 아직 살아있는지 확인하기 위해 서버에 ping 을 보냅니다.
	Timeout:             time.Second,      // 연결이 끊어진 것으로 간주하기 전에 ping ack 을 1초 동안 기다립니다.
	PermitWithoutStream: true,             // 활성 스트림 없이도 ping 보내기
}

type uploader struct {
	dir    string
	client streamProtoc.UploadFileServiceClient
	ctx    context.Context
	wg     sync.WaitGroup
	pool   *pb.Pool

	requests    chan string // 각 요청은 클라이언트가 액세스할 수 있는 클라이언트의 파일 경로입니다.
	DoneRequest chan string
	FailRequest chan string
}

func NewUploader(ctx context.Context, client streamProtoc.UploadFileServiceClient, dir string) *uploader {
	d := &uploader{
		ctx:         ctx,
		client:      client,
		dir:         dir,
		requests:    make(chan string),
		DoneRequest: make(chan string),
		FailRequest: make(chan string),
	}
	// 한번에 몇개의 파일을 보낼지
	for i := 0; i < 5; i++ {
		d.wg.Add(1)
		go func(i int) {
			if err := d.worker(i + 1); err != nil {
				_ = errors.Wrapf(err, "worker methods is not running")
			}
		}(i)
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
	if err := d.pool.Stop(); err != nil {
		_ = errors.Wrapf(err, "pool was failed stoping")
	}
}

//goland:noinspection ALL
func (d *uploader) worker(_ int) error {
	var (
		buf        []byte
		firstChunk bool
	)
	defer d.wg.Done()

	// 파일 경로에서 파일들을 추출
	for request := range d.requests {
		file, err := os.Open(request)
		if err != nil {
			return errors.Wrapf(err, "failed to open file %s", request)
		}

		// start uploader
		stream, err := d.client.Upload(d.ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to create upload stream for file %s", request)
		}

		// file info
		stat, err := file.Stat()
		if err != nil {
			return errors.Wrapf(err, "unable to get file size  %s", request)
		}

		// client 측에서 server에 send할때 진행률 bar 표시
		bar := pb.New64(stat.Size()).Postfix(" " + filepath.Base(request))
		bar.Units = pb.U_BYTES
		d.pool.Add(bar)

		// 파일 크기만큼 바이트 슬라이스 생성
		buf = make([]byte, chunkSize)
		firstChunk = true
		for {
			// 파일의 내용을 읽어서 바이트 슬라이스에 저장
			f, err := file.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return errors.Wrapf(err, "errored while copying from file to buf")
			}

			if firstChunk {
				// Send to Server
				err = stream.Send(&streamProtoc.UploadRequest{
					Content:  buf[:f],
					Filename: request,
				})
				firstChunk = false
			} else {
				err = stream.Send(&streamProtoc.UploadRequest{
					Content: buf[:f],
				})
			}
			if err != nil {
				bar.FinishPrint("failed to send chunk via stream file : " + request)
				break
			}
			bar.Add64(int64(f))
		}

		// 클라이언트에서 send를 완료하고 서버에서 다 받고나서 완료 메세지를 보내면
		// 그 메세지를 받아서 상태를 체크
		if status, err := stream.CloseAndRecv(); err != nil || status.Code != streamProtoc.UploadStatusCode_Ok {
			fmt.Printf("Error uploading file : " + request + " :" + status.Message)
			bar.FinishPrint("Error uploading file : " + request + " Error :" + err.Error())
			bar.Reset(0)
			d.FailRequest <- request
			return nil
		}

		file.Close()
		d.DoneRequest <- request
		bar.Finish()
	}
	return nil
}

func UploadFile(ctx context.Context, client streamProtoc.UploadFileServiceClient, filePathList []string, dir string) error {
	var errUpload error

	d := NewUploader(ctx, client, dir)

	if dir != "" {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return errors.Wrapf(err, "can't read th dir")
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
					errUpload = errors.Wrapf(errUpload, " Failed to send %s", req)
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
				errUpload = errors.Wrapf(errUpload, " Failed to send %s", req)
			}
		}
	}
	return errUpload
}

//goland:noinspection ALL
func uploadCommand() cli.Command {
	return cli.Command{
		Name:  "upload",
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
				if err := conn.Close(); err != nil {
					fmt.Println("faild conn close")
				}
			}(conn)

			return UploadFile(
				context.Background(),
				streamProtoc.NewUploadFileServiceClient(conn),
				[]string{},
				c.String("d"),
			)
		},
	}
}
