package app

import (
	"database/sql"
	"fmt"
	"hex-postgres-grpc/internal/auth"
	authpg "hex-postgres-grpc/internal/auth/adapters/postgres"
	"hex-postgres-grpc/internal/category"
	"hex-postgres-grpc/internal/customer"
	"hex-postgres-grpc/internal/order"
	"hex-postgres-grpc/internal/product"

	_ "github.com/lib/pq"
)

type Application struct {
	DB          *sql.DB
	Order       order.Components
	Product     product.Components
	Customer    customer.Components
	Category    category.Component
	Auth        auth.Service
	AuthHandler *auth.Handler
	AuthRepo    auth.UserRepository
}

func Init(cfg DBConfig) (*Application, error) {
	db, err := initDB(cfg)
	if err != nil {
		return nil, err
	}

	authRepo := authpg.NewRepository(db)
	authSvc := auth.NewService("super-secret-key", authRepo)
	authHandler := auth.NewHandler(authSvc)

	return &Application{
		DB:          db,
		Order:       order.Init(db, authSvc),
		Product:     product.Init(db, authSvc),
		Customer:    customer.Init(db, authSvc),
		Category:    category.Init(db, authSvc),
		Auth:        authSvc,
		AuthHandler: authHandler,
		AuthRepo:    authRepo,
	}, nil
}

func initDB(cfg DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.HostPort, cfg.Name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

type DBConfig struct {
	User, Password, HostPort, Name string
}
