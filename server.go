package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "nodeAgent/inter"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedAgentServer
}

func (s *server) Config(ctx context.Context, in *pb.AgentRequest) (*pb.AgentReply, error) {

	log.Printf("Method: %v", in.GetMethod())
	log.Printf("Params: %v", in.GetParams())
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
