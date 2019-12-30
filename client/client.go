package main

import (
	"context"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"time"

	pb "github.com/suzushin54/grpc-study"
	"google.golang.org/grpc"
)

func main() {
	addr := "localhost:50051"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	name := os.Args[1]

	ctx := context.Background()
	// communicate by adding metadata when calling RPC
	md := metadata.Pairs("timestamp", time.Now().Format(time.Stamp))
	ctx = metadata.NewOutgoingContext(ctx, md)
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name}, grpc.Trailer(&md))
	//r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			log.Printf("gRPC Error (message: %s)", s.Message())
			for _, d := range s.Details() {
				switch info := d.(type) {
				case *errdetails.RetryInfo:
					log.Printf(" RetryInfo: %v", info)
				}
			}
			os.Exit(1)
		} else {
			log.Fatalf("could not greet: %v", err)
		}
	}
	log.Printf("Greeting: %s", r.Message)

}
