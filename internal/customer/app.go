package customer

import (
	"database/sql"
	"hex-postgres-grpc/internal/auth"
	"hex-postgres-grpc/internal/customer/adapters/grpc"
	"hex-postgres-grpc/internal/customer/adapters/http"
	"hex-postgres-grpc/internal/customer/adapters/postgres"
	"hex-postgres-grpc/internal/customer/usecase"
)

type Components struct {
	HTTPHandler *http.Handler
	GRPCServer  *grpc.Server
}

func Init(db *sql.DB, authSvc auth.Service) Components {
	repo := postgres.NewRepository(db)
	service := usecase.NewService(repo)

	httpHandler := http.NewHandler(service, authSvc)
	grpcServer := grpc.NewServer(service, authSvc)

	return Components{
		HTTPHandler: httpHandler,
		GRPCServer:  grpcServer,
	}
}
