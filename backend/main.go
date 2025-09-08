package main

import (
	"PAYROLL_SYSTEM/backend/proto"
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcAddr := ":9090"
	httpAddr := ":8080"

	mongoURI := "mongodb+srv://payrolluser:payrolluser123@cluster0.mcefsvj.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	grpcServer := grpc.NewServer()
	srv := newServer(mongoURI)
	proto.RegisterPayrollServiceServer(grpcServer, srv)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}

	go func() {
		log.Println("gRPC server on", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("grpc failed: %v", err)
		}
	}()

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := proto.RegisterPayrollServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		log.Fatalf("gateway reg failed: %v", err)
	}

	log.Println("HTTP server on", httpAddr)
	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatalf("http failed: %v", err)
	}
}
