package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	zmq "github.com/pebbe/zmq4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "grpc_colleen/protos/longLived" // import the generated protobuf package
)

type server struct {
	pb.UnimplementedLonglivedServer
}

func (s *server) Subscribe(req *pb.Request, stream pb.Longlived_SubscribeServer) error {

	//create context to containerize sockets
	zctx, _ := zmq.NewContext()
	zmqServer, _ := zctx.NewSocket(zmq.REQ)
	zmqServer.Connect("tcp://localhost:5555")
	zmqServer.Send("100000", 0)

	mySub := createSubSocket(zctx)

	time.Sleep(1 * time.Second)

	for {

		msg, _ := mySub.Recv(0)

		select {
		case <-stream.Context().Done():
			return nil
		default:
			err := stream.Send(&pb.Response{
				Data: msg,
			})
			if err != nil {
				return err
			}
		}
	}
}

func genRandomNum() string {
	x := rand.Int31()
	hexX := fmt.Sprintf("%x", x)
	return hexX
}

func createSubSocket(zctx *zmq.Context) *zmq.Socket {

	sub, _ := zctx.NewSocket(zmq.SUB)
	sub.Connect("tcp://localhost:5556")
	sub.SetSubscribe((""))
	fmt.Println("\nStarting to transmit data from the zeroMQ Pub/Sub socket:")
	return sub

}

func main() {
	// Start the gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLonglivedServer(s, &server{})
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
