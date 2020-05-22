package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "nodeAgent/inter"
	"nodeAgent/tools"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAgentClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	params := make(map[string]interface{})
	params["path"] = "/data/config/a.conf"
	params["content"] = "{\"a\":1}"
	r, err := c.Config(ctx, &pb.AgentRequest{Method: "save", Params: tools.JsonEncode(params)})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
