package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"

	streamPb "github.com/jaden7856/go-grpcUpload/streamProtoc"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var kaep = keepalive.EnforcementPolicy{
	// 클라이언트가 5초에 한 번 이상 ping 을 보내는 경우 연결을 종료합니다.
	MinTime: 5 * time.Second,
	// 활성 스트림이 없는 경우에도 ping 허용
	PermitWithoutStream: true,
}

var kasp = keepalive.ServerParameters{
	// 클라이언트가 60초 동안 유휴 상태이면 GOAWAY 를 보냅니다.
	MaxConnectionIdle: 15 * time.Second,
	// 연결이 여전히 활성 상태인지 확인하기 위해 5초 동안 클라이언트가
	// 유휴 상태인 경우 클라이언트를 ping 합니다.
	Time: 5 * time.Second,
	// 연결이 끊어졌다고 가정하기 전에 ping ack 를 1초 동안 기다립니다.
	Timeout: 1 * time.Second,
}

var (
	fp       *os.File
	fileData *streamPb.UploadRequest
	filename string
)

type ServerGRPC struct {
	streamPb.UnimplementedUploadFileServiceServer
	server *grpc.Server

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
		err = errors.Errorf("address must be specified")
		return
	}

	s.Address = cfg.Address
	// 지정한 경로에 파일이 있는지 확인
	if _, err = os.Stat(cfg.DestDir); err != nil {
		fmt.Println("Directory doesn't exist")
		return
	}

	s.destDir = cfg.DestDir

	return
}

func (s *ServerGRPC) Listen() (err error) {
	lis, err := net.Listen("tcp", s.Address)
	if err != nil {
		err = errors.Wrapf(err, "failed to listen on  %s", s.Address)
		return
	}

	// gRPC 서버 생성
	s.server = grpc.NewServer(
		// Server keepalive
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
	)

	// gRPC 서버 등록
	streamPb.RegisterUploadFileServiceServer(s.server, s)

	if err = s.server.Serve(lis); err != nil {
		err = errors.Wrapf(err, "errored listening for grpc connections")
		return
	}

	return
}

// writeToFp는 파일 포인터와 바이트 배열을 가져와서 파일에 바이트 배열을 씁니다.
// 포인터가 nil 이거나 파일 쓰기 오류인 경우 오류를 반환합니다.
func writeToFp(fp *os.File, data []byte) error {
	w := 0
	n := len(data)
	for {
		nw, err := fp.Write(data[w:])
		if err != nil {
			return errors.Wrapf(err, "failed write data")
		}
		w += nw
		if nw >= n {
			return nil
		}
	}
}

//goland:noinspection ALL
func (s *ServerGRPC) Upload(stream streamPb.UploadFileService_UploadServer) (err error) {
	firstChunk := true

	for {
		fileData, err = stream.Recv()
		if err != nil {
			// 더이상 들어오는 파일이없으면 for문 break
			// err == io.EOF 와 같은 의미 -> 클린소스
			if errors.Is(err, io.EOF) {
				break
			}
			err = errors.Wrapf(err, "failed while reading from stream")
			return
		}

		filename = fileData.Filename

		if firstChunk {
			// create file
			if filename != "" {
				fp, err = os.Create(path.Join(s.destDir, filepath.Base(filename)))
				if err != nil {
					fmt.Println("Unable to create file : " + filename)
					stream.SendAndClose(&streamPb.UploadResponse{
						Message: "Unable to create file : " + filename,
						Code:    streamPb.UploadStatusCode_Failed,
					})
					return
				}
				fp.Close()
			} else if filename == "" {
				fmt.Println("FileName not provided in first chunk : " + filename)
				stream.SendAndClose(&streamPb.UploadResponse{
					Message: "FileName not provided in first chunk : " + filename,
					Code:    streamPb.UploadStatusCode_Failed,
				})
				return
			}
			firstChunk = false
		}

		// 생성된 파일에 전송된 data들을 write하고, 에러 발생시 전송 후 종료
		if err = writeToFp(fp, fileData.Content); err != nil {
			fmt.Println("Unable to write chunk of filename : " + filename + " " + err.Error())
			stream.SendAndClose(&streamPb.UploadResponse{
				Message: "Unable to write chunk of filename : " + filename,
				Code:    streamPb.UploadStatusCode_Failed,
			})
			return
		}
	}

	// 정상적으로 다 파일을 받고나서 클아이언트에 정상적으로 완료가 되었다는 메세지를 전달
	if err = stream.SendAndClose(&streamPb.UploadResponse{
		Message: "Upload received with success",
		Code:    streamPb.UploadStatusCode_Ok,
	}); err != nil {
		err = errors.Wrapf(err, "failed to send status code")
		return
	}

	fmt.Println("Successfully received :" + filename + " in " + s.destDir)
	return
}

func (s *ServerGRPC) Close() {
	if s.server != nil {
		s.server.Stop()
	}
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
				Usage: "Destination directory Default is /tmp",
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
			if err = server.Listen(); err != nil {
				return errors.Wrapf(err, "can't connect with client")
			}

			defer server.Close()
			return nil
		},
	}
}
