package category

import (
	"database/sql"
	"hex-postgres-grpc/internal/auth"
	"hex-postgres-grpc/internal/category/adapter/grpc"
	"hex-postgres-grpc/internal/category/adapter/http"
	"hex-postgres-grpc/internal/category/adapter/postgres"
	"hex-postgres-grpc/internal/category/usecase"
)

type Component struct {
	HTTPHandler *http.Handler
	GRPCHandler *grpc.Server
}

func Init(db *sql.DB, authSvc auth.Service) Component {
	repo := postgres.NewRepository(db)
	service := usecase.NewService(repo)
	httpHandler := http.NewHandler(service, authSvc)
	grpcHandler := grpc.NewCategoryGRPCServer(service)

	return Component{
		HTTPHandler: httpHandler,
		GRPCHandler: grpcHandler,
	}
}
