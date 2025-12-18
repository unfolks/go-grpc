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
	createUC := usecase.NewCreateProduct(repo)
	getUC := usecase.NewGetProduct(repo)
	listUC := usecase.NewListProducts(repo)
	updateUC := usecase.NewUpdateProduct(repo)
	deleteUC := usecase.NewDeleteProduct(repo)

	httpHandler := http.NewHandler(createUC, getUC, listUC, updateUC, deleteUC)
	grpcServer := grpc.NewProductGRPCServer(createUC, getUC, listUC, updateUC, deleteUC)

	return Components{
		HTTPHandler: httpHandler,
		GRPCServer:  grpcServer,
	}
}
