package product

import (
	"database/sql"
	"hex-postgres-grpc/internal/auth"
	"hex-postgres-grpc/internal/product/adapters/grpc"
	"hex-postgres-grpc/internal/product/adapters/http"
	"hex-postgres-grpc/internal/product/adapters/postgres"
	"hex-postgres-grpc/internal/product/usecase"
)

type Components struct {
	HTTPHandler *http.Handler
	GRPCServer  *grpc.Server
}

func Init(db *sql.DB, authSvc auth.Service) Components {
	repo := postgres.NewProductRepoPG(db)
	service := usecase.NewService(repo)

	httpHandler := http.NewHandler(service, authSvc)
	grpcServer := grpc.NewProductGRPCServer(service) // TODO: Add auth interceptor if needed

	return Components{
		HTTPHandler: httpHandler,
		GRPCServer:  grpcServer,
	}
}
