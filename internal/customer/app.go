package customer

import (
	"database/sql"
	"hex-postgres-grpc/internal/customer/adapters/http"
	"hex-postgres-grpc/internal/customer/adapters/postgres"
	"hex-postgres-grpc/internal/customer/usecase"
)

type Components struct {
	HTTPHandler *http.Handler
}

func Init(db *sql.DB) Components {
	repo := postgres.NewRepository(db)
	createUC := usecase.NewCreateCustomer(repo)
	listUC := usecase.NewListCustomer(repo)

	httpHandler := http.NewHandler(createUC, listUC)

	return Components{
		HTTPHandler: httpHandler,
	}
}
