package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"

	streamPb "github.com/jaden7856/go-grpcUpload/streamProtoc"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type ServerGRPC struct {
	streamPb.UnimplementedUploadFileServiceServer
	server  *grpc.Server
	Address string

	destDir string
}

type ServerGRPCConfig struct {
	Address string
	DestDir string
}

func NewServerGRPC(cfg ServerGRPCConfig) (s ServerGRPC, err error) {
	// 주소 입력 여부 확인
	if cfg.Address == "" {
		err = errors.Errorf("Address must be specified")
		return
	}

	s.Address = cfg.Address
	// 지정한 경로에 파일이 있는지 확인
	if _, err = os.Stat(cfg.DestDir); err != nil {
		fmt.Println("Directory doesnt exist")
		return
	}

	s.destDir = cfg.DestDir

	return
}

func (s *ServerGRPC) Listen() (err error) {
	var (
		lis         net.Listener
		serviceName = "test"
	)

	lis, err = net.Listen("tcp", s.Address)
	if err != nil {
		err = errors.Wrapf(err, "failed to listen on  %d", s.Address)
		return
	}

	// gRPC 서버 생성
	s.server = grpc.NewServer()
	// healthCheck 서버 생성
	healthServer := health.NewServer()

	healthpb.RegisterHealthServer(s.server, healthServer)
	streamPb.RegisterUploadFileServiceServer(s.server, s)

	// 정상적으로 연결이 된 상태
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	if err = s.server.Serve(lis); err != nil {
		err = errors.Wrapf(err, "errored listening for grpc connections")
		return
	}

	return
}

// writeToFp는 파일 포인터와 바이트 배열을 가져와서 파일에 바이트 배열을 씁니다.
// 포인터가 nil이거나 파일 쓰기 오류인 경우 오류를 반환합니다.
func writeToFp(fp *os.File, data []byte) error {
	w := 0
	n := len(data)
	for {
		nw, err := fp.Write(data[w:])
		if err != nil {
			return err
		}
		w += nw
		if nw >= n {
			return nil
		}
	}

}

//goland:noinspection ALL
func (s *ServerGRPC) Upload(stream streamPb.UploadFileService_UploadServer) (err error) {
	log.Println("start new server")
	firstChunk := true

	var (
		fp       *os.File
		fileData *streamPb.UploadRequest
		filename string
	)

	for {
		fileData, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			err = errors.Wrapf(err, "failed while reading from stream")
			return
		}

		// first chunk에 file name이 포함되어 있으면
		if firstChunk {

			if fileData.Filename != "" { //create file
				fp, err = os.Create(path.Join(s.destDir, filepath.Base(fileData.Filename)))
				if err != nil {
					fmt.Println("Unable to create file : " + fileData.Filename)
					stream.SendAndClose(&streamPb.UploadResponse{
						Message: "Unable to create file : " + fileData.Filename,
						Code:    streamPb.UploadStatusCode_Failed,
					})
					return
				}
				defer fp.Close()

			} else {

				fmt.Println("FileName not provided in first chunk : " + fileData.Filename)
				stream.SendAndClose(&streamPb.UploadResponse{
					Message: "FileName not provided in first chunk : " + fileData.Filename,
					Code:    streamPb.UploadStatusCode_Failed,
				})

				return
			}

			filename = fileData.Filename
			firstChunk = false
		}

		err = writeToFp(fp, fileData.Content)
		if err != nil {

			fmt.Println("Unable to write chunk of filename : " + filename + " " + err.Error())
			stream.SendAndClose(&streamPb.UploadResponse{
				Message: "Unable to write chunk of filename : " + filename,
				Code:    streamPb.UploadStatusCode_Failed,
			})
			return
		}
	}

	// send to Upload
	fmt.Println("여기까진 오냐?")
	err = stream.SendAndClose(&streamPb.UploadResponse{
		Message: "Upload received with success",
		Code:    streamPb.UploadStatusCode_Ok,
	})

	if err != nil {
		err = errors.Wrapf(err, "failed to send status code")
		return
	}
	fmt.Println("Successfully received and stored the file :" + filename + " in " + s.destDir)
	return
}

func (s *ServerGRPC) Close() {
	if s.server != nil {
		s.server.Stop()
	}

	return
}

func ServerCommand() cli.Command {
	return cli.Command{
		Name:  "serve",
		Usage: "gRPC server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "a",
				Usage: "Address to listen",
				Value: "localhost:80",
			},
			&cli.StringFlag{
				Name:  "d",
				Usage: "Destrination directory Default is /tmp",
				Value: "/tmp",
			},
		},
		Action: func(c *cli.Context) error {
			grpcServer, err := NewServerGRPC(ServerGRPCConfig{
				Address: c.String("a"),
				DestDir: c.String("d"),
			})
			if err != nil {
				fmt.Println("error is creating server")

				return err
			}

			server := &grpcServer
			err = server.Listen()

			defer server.Close()
			return nil
		},
	}
}
