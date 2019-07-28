package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/yboyacigil/grpc-rest-echo-server/pb"
	"google.golang.org/grpc"
)

// EchoServer ...
type EchoServer struct {
	wg sync.WaitGroup
}

// New creates a new EchoServer
func New() *EchoServer {
	return &EchoServer{}
}

// Start starts server
func (e *EchoServer) Start() {
	e.wg.Add(1)
	go func() {
		log.Fatal(e.startGRPC())
		e.wg.Done()
	}()
	e.wg.Add(1)
	go func() {
		log.Fatal(e.startREST())
		e.wg.Done()
	}()
}

// WaitStop waits for grpc and rest endpoints
func (e *EchoServer) WaitStop() {
	e.wg.Wait()
}

func (e *EchoServer) startGRPC() error {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterEchoServiceServer(grpcServer, e)
	grpcServer.Serve(lis)
	return nil
}

func (e *EchoServer) startREST() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterEchoServiceHandlerFromEndpoint(ctx, mux, ":8080", opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":8090", mux)
}

// Echo echoes message
func (e *EchoServer) Echo(ctx context.Context, r *pb.EchoMessage) (*pb.EchoMessage, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	return &pb.EchoMessage{
		Value: fmt.Sprintf("Echo: %s", r.Value),
	}, nil
}
