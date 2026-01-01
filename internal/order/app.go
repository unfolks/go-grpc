package order

import (
	"database/sql"
	"hex-postgres-grpc/internal/auth"
	"hex-postgres-grpc/internal/order/adapters/grpc"
	"hex-postgres-grpc/internal/order/adapters/http"
	"hex-postgres-grpc/internal/order/adapters/postgres"
	orderdomain "hex-postgres-grpc/internal/order/domain"
)

type Components struct {
	Service     orderdomain.Service
	HTTPHandler *http.Handler
	GRPCServer  *grpc.Server
}

func Init(db *sql.DB, authSvc auth.Service) Components {
	repo := postgres.NewOrderRepoPG(db)
	svc := orderdomain.NewService(repo)
	httpHandler := http.NewHandler(svc, authSvc)
	grpcServer := grpc.NewOrderGRPCServer(svc) // TODO: Add auth interceptor if needed

	return Components{
		Service:     svc,
		HTTPHandler: httpHandler,
		GRPCServer:  grpcServer,
	}
}
