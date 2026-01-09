package main

import (
	"hex-postgres-grpc/internal/app"
	categorypb "hex-postgres-grpc/proto/category"
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

// @title Hex Postgres gRPC API
// @version 1.0
// @description This is a sample server for a hexagonal architecture with Postgres and gRPC.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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
		a.AuthHandler.RegisterRoutes(mux) // Register auth routes
		a.Order.HTTPHandler.RegisterRoutes(mux)
		a.Product.HTTPHandler.RegisterRoutes(mux)
		a.Customer.HTTPHandler.RegisterRoutes(mux)
		a.Category.HTTPHandler.RegisterRoutes(mux)

		// Wrap mux with Auth middleware
		handler := a.Auth.HTTPMiddleware(mux)

		log.Println("HTTP listening :8080")
		log.Fatal(http.ListenAndServe(":8080", handler))
	}()

	grpcLis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("grpc listen %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(a.Auth.GRPCUnaryInterceptor),
	)
	orderpb.RegisterORderServiceServer(grpcServer, a.Order.GRPCServer)
	productpb.RegisterProductServiceServer(grpcServer, a.Product.GRPCServer)
	customerpb.RegisterCustomerServiceServer(grpcServer, a.Customer.GRPCServer)
	categorypb.RegisterCategoryServiceServer(grpcServer, a.Category.GRPCHandler)
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
