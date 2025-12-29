package product

import (
	"database/sql"
	"hex-postgres-grpc/internal/product/adapters/grpc"
	"hex-postgres-grpc/internal/product/adapters/http"
	"hex-postgres-grpc/internal/product/adapters/postgres"
	"hex-postgres-grpc/internal/product/usecase"
)

type Components struct {
	HTTPHandler *http.Handler
	GRPCServer  *grpc.Server
}

func Init(db *sql.DB) Components {
	repo := postgres.NewProductRepoPG(db)
	service := usecase.NewService(repo)

	httpHandler := http.NewHandler(service)
	grpcServer := grpc.NewProductGRPCServer(service)

	return Components{
		HTTPHandler: httpHandler,
		GRPCServer:  grpcServer,
	}
}
