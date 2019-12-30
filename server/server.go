package main

import (
	"context"
	"time"

	//"errors"
	"log"
	"net"

	pb "github.com/suzushin54/grpc-study"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	time.Sleep(3 * time.Second)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
	//st, _ := status.New(codes.Aborted, "aborted").WithDetails(&errdetails.RetryInfo{
	//	// request retry after 3 seconds
	//	RetryDelay:           &duration.Duration{
	//		Seconds:              3,
	//		Nanos:                0,
	//		//XXX_NoUnkeyedLiteral: struct{}{},
	//		//XXX_unrecognized:     nil,
	//		//XXX_sizecache:        0,
	//	},
	//	//XXX_NoUnkeyedLiteral: struct{}{},
	//	//XXX_unrecognized:     nil,
	//	//XXX_sizecache:        0,
	//})
	//return nil, st.Err()
	//return &pb.HelloReply{Message: "Hello " + in.Name}, nil
	// sample of error responses
	//return nil, status.New(codes.NotFound, "resource not found").Err()
}

func main() {
	addr := ":50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	log.Printf("gRPC server listening on " + addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
