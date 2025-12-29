package main

import (
	"hex-postgres-grpc/internal/app"
	customerpb "hex-postgres-grpc/proto/customer"
	orderpb "hex-postgres-grpc/proto/order"
	productpb "hex-postgres-grpc/proto/product"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	cfg := app.DBConfig{
		User:     "admin",
		Password: "password123",
		HostPort: "localhost:5432",
		Name:     "ordersdb",
	}
	a, err := app.Init(cfg)
	if err != nil {
		log.Fatalf("init app: %v", err)
	}
	go func() {
		mux := http.NewServeMux()
		a.Order.HTTPHandler.RegisterRoutes(mux)
		a.Product.HTTPHandler.RegisterRoutes(mux)
		a.Customer.HTTPHandler.RegisterRoutes(mux)
		log.Println("HTTP listening :8080")
		log.Fatal(http.ListenAndServe(":8080", mux))
	}()

	grpcLis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("grpc listen %v", err)
	}
	grpcServer := grpc.NewServer()
	orderpb.RegisterORderServiceServer(grpcServer, a.Order.GRPCServer)
	productpb.RegisterProductServiceServer(grpcServer, a.Product.GRPCServer)
	customerpb.RegisterCustomerServiceServer(grpcServer, a.Customer.GRPCServer)
	go func() {
		log.Println("gRPC listening :50051")
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Println("shutting down gracefully")
	grpcServer.GracefulStop()
	_ = a.DB.Close()
}
