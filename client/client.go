package main

import (
	"context"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"time"

	pb "github.com/suzushin54/grpc-study"
	"google.golang.org/grpc"
)

func main() {
	// Register resolver for load balancing
	resolver.Register(&exampleResolverBuilder{})
	addr := "testScheme:///example"

	// use TLS authentication
	//cred, err := credentials.NewClientTLSFromFile("server.crt", "")
	//if err != nil {
	//	log.Fatal(err)
	//}
	// rpcの処理の前後にInterceptorでログ出力処理を差し込む
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		//grpc.WithTransportCredentials(cred),
		grpc.WithUnaryInterceptor(unaryInterceptor),
		grpc.WithBalancerName("round_robin"))

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

func unaryInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("before call: %s, request: %+v", method, req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("after call: %s, response: %+v", method, reply)
	return err
}

type exampleResolverBuilder struct{}

func (*exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target: 		target,
		cc:				cc,
		addressStore: 	map[string][]string{
			"example": {"localhost:50051", "localhost:50052"},
		},
	}
	r.start()
	return r, nil
}

func (*exampleResolverBuilder) Scheme() string { return "testScheme" }

type exampleResolver struct {
	target			resolver.Target
	cc				resolver.ClientConn
	addressStore	map[string][]string
}

func(r *exampleResolver) start() {
	addressStrings := r.addressStore[r.target.Endpoint]
	address := make([]resolver.Address, len(addressStrings))
	for i, s := range addressStrings {
		address[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: address})
}
func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*exampleResolver) Close()									 {}

