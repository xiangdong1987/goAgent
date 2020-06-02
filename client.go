package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	pb "nodeAgent/inter"
	"nodeAgent/tools"
	"os"
	"time"
)

type Stats struct {
	StartedAt  time.Time
	FinishedAt time.Time
}

type ClientGRPC struct {
	logger    zerolog.Logger
	conn      *grpc.ClientConn
	client    pb.AgentClient
	chunkSize int
}

type ClientGRPCConfig struct {
	Address         string
	ChunkSize       int
	RootCertificate string
	Compress        bool
}

func main() {
	//uploadFile("/Users/xiangdd/go/src/goAgent/code_version/test.tar.gz", "/Users/xiangdd/go/src/goAgent/code_version/upload/test.tar.gz")
	params := make(map[string]interface{})
	params["path"] = "b.conf"
	params["content"] = "{\"b\":2}"
	config("save", params)
}

//配置上传
func config(method string, params map[string]interface{}) {
	config := &ClientGRPCConfig{Address: "localhost:50051", ChunkSize: 1024}
	client, err := NewClientGRPC(*config)
	if err != nil {
		log.Fatalf("New connnect fail: %v", err)
	}
	Stats, err := client.Config(context.Background(), method, params)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(Stats)
}

//文件上传
func uploadFile(file string, path string) {
	config := &ClientGRPCConfig{Address: "localhost:50051", ChunkSize: 1024}
	client, err := NewClientGRPC(*config)
	if err != nil {
		log.Fatalf("New connnect fail: %v", err)
	}
	Stats, err := client.UploadFile(context.Background(), file, path)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(Stats)
}

func NewClientGRPC(cfg ClientGRPCConfig) (c ClientGRPC, err error) {
	var (
		grpcOpts  = []grpc.DialOption{}
		grpcCreds credentials.TransportCredentials
	)
	if cfg.Address == "" {
		err = errors.Errorf("address must be specified")
		return
	}
	if cfg.Compress {
		grpcOpts = append(grpcOpts, grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	}
	if cfg.RootCertificate != "" {
		grpcCreds, err = credentials.NewClientTLSFromFile(cfg.RootCertificate, "localhost")
		if err != nil {
			err = errors.Wrapf(err,
				"failed to create grpc tls client via root-cert %s",
				cfg.RootCertificate)
			return
		}
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcCreds))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}
	switch {
	case cfg.ChunkSize == 0:
		err = errors.Errorf("ChunkSize must be specified")
		return
	case cfg.ChunkSize > (1 << 22):
		err = errors.Errorf("ChunkSize must be < than 4MB")
		return
	default:
		c.chunkSize = cfg.ChunkSize
	}
	c.logger = zerolog.New(os.Stdout).With().Str("from", "client").Logger()
	c.conn, err = grpc.Dial(cfg.Address, grpcOpts...)
	if err != nil {
		err = errors.Wrapf(err, "failed to start grpc connection with address %s", cfg.Address)
		return
	}
	c.client = pb.NewAgentClient(c.conn)
	return
}

//配置下发
func (c *ClientGRPC) Config(ctx context.Context, method string, params map[string]interface{}) (reply *pb.AgentReply, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err = c.client.Config(ctx, &pb.AgentRequest{Method: method, Params: tools.JsonEncode(params)})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	return reply, err
}

//文件上传
func (c *ClientGRPC) UploadFile(ctx context.Context, f string, path string) (stats Stats, err error) {
	var (
		writing = true
		buf     []byte
		n       int
		file    *os.File
		status  *pb.UploadStatus
	)

	file, err = os.Open(f)
	if err != nil {
		err = errors.Wrapf(err, "failed to open file %s", f)
		return
	}
	defer file.Close()
	stream, err := c.client.Upload(ctx)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to create upload stream for file %s",
			f)
		return
	}
	defer stream.CloseSend()
	stats.StartedAt = time.Now()
	buf = make([]byte, c.chunkSize)
	for writing {
		n, err = file.Read(buf)
		if err != nil {
			if err == io.EOF {
				writing = false
				err = nil
				continue
			}
			err = errors.Wrapf(err, "errored while copying from file to buf")
			return
		}
		err = stream.Send(&pb.Chunk{
			Path:    path,
			Content: buf[:n],
		})
		if err != nil {
			err = errors.Wrapf(err, "failed to send chunk via stream")
			return
		}
	}
	stats.FinishedAt = time.Now()
	status, err = stream.CloseAndRecv()
	if err != nil {
		err = errors.Wrapf(err, "failed to receive upstream status response")
		return
	}
	if status.Code != pb.UploadStatusCode_Ok {
		err = errors.Errorf(
			"upload failed - msg: %s",
			status.Message)
		return
	}
	return
}
