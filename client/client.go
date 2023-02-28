package main

import (
	"context"
	"fmt"
	"log"

	pb "grpc_colleen/protos/longLived" // import the generated protobuf package

	"google.golang.org/grpc"
)

func main() {
	// Set up a gRPC connection to the server
	conn, err := grpc.Dial("192.168.0.144:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create a Longlived client and subscribe to the server
	client := pb.NewLonglivedClient(conn)
	stream, err := client.Subscribe(context.Background(), &pb.Request{})
	if err != nil {
		log.Fatalf("failed to subscribe: %v", err)
	}

	// Continuously print the data received from the server
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Fatalf("failed to receive data: %v", err)
		}
		fmt.Println(resp.Data)
	}
}
