package order

import (
	"database/sql"
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

func Init(db *sql.DB) Components {
	repo := postgres.NewOrderRepoPG(db)
	svc := orderdomain.NewService(repo)
	httpHandler := http.NewHandler(svc)
	grpcServer := grpc.NewOrderGRPCServer(svc)

	return Components{
		Service:     svc,
		HTTPHandler: httpHandler,
		GRPCServer:  grpcServer,
	}
}
