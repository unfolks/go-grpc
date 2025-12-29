package customer

import (
	"database/sql"
	"hex-postgres-grpc/internal/customer/adapters/grpc"
	"hex-postgres-grpc/internal/customer/adapters/http"
	"hex-postgres-grpc/internal/customer/adapters/postgres"
	"hex-postgres-grpc/internal/customer/usecase"
)

type Components struct {
	HTTPHandler *http.Handler
	GRPCServer  *grpc.Server
}

func Init(db *sql.DB) Components {
	repo := postgres.NewRepository(db)
	service := usecase.NewService(repo)

	httpHandler := http.NewHandler(service)
	grpcServer := grpc.NewServer(service)

	return Components{
		HTTPHandler: httpHandler,
		GRPCServer:  grpcServer,
	}
}
