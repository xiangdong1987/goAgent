package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"nodeAgent/fun"
	pb "nodeAgent/inter"
	"os"
)

const (
	port = ":50051"
)

//agentServer
type server struct {
	logger zerolog.Logger
	pb.UnimplementedAgentServer
}

//文件上传
func (s *server) Upload(stream pb.Agent_UploadServer) (err error) {
	fmt.Println("start upload")
	var file []byte
	var fileName string
	for {
		x, err := stream.Recv()
		//fmt.Println(x.Path)
		if err != nil {
			if err == io.EOF {
				goto END
			}
			err = errors.Wrapf(err, "failed unexpectadely while reading chunks from stream")
			return nil
		} else {
			fileName = x.Path
			file = append(file, x.Content...)
		}
	}
	fmt.Println("upload received")
END:
	//写入文件
	fd, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	buf := file
	fd.Write(buf)
	fd.Close()
	err = stream.SendAndClose(&pb.UploadStatus{
		Message: "Upload received with success",
		Code:    pb.UploadStatusCode_Ok,
	})
	if err != nil {
		err = errors.Wrapf(err, "failed to send status code")
		return
	}
	return
}

//配置下发
func (s *server) Config(ctx context.Context, in *pb.AgentRequest) (*pb.AgentReply, error) {
	log.Printf("Method: %v", in.GetMethod())
	log.Printf("Params: %v", in.GetParams())
	switch in.GetMethod() {
	case "save":
		fun.Save(in.GetParams())
	default:
		log.Println("No method")
	}
	return &pb.AgentReply{Message: "Method :" + in.GetMethod() + " Params:" + in.GetParams()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAgentServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
