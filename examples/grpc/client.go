package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/mike955/zrpc/examples/grpc/api/example"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	grpcAddress = "0.0.0.0:5180"
	defaultName = "world"
)

func main() {
	conn, err := grpc.DialContext(context.Background(), grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewExampleClient(conn)
	md := metadata.Pairs("traceId", "trace-id-value", "key2", "value2")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	response, err := c.Hello(ctx, &pb.HelloRequest{
		Data: "grpc",
	})
	fmt.Println(response.Data)
}
