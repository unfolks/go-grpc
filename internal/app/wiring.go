package app

import (
	"database/sql"
	"fmt"
	"hex-postgres-grpc/internal/customer"
	"hex-postgres-grpc/internal/order"
	"hex-postgres-grpc/internal/product"

	_ "github.com/lib/pq"
)

type Application struct {
	DB       *sql.DB
	Order    order.Components
	Product  product.Components
	Customer customer.Components
}

func Init(cfg DBConfig) (*Application, error) {
	db, err := initDB(cfg)
	if err != nil {
		return nil, err
	}

	return &Application{
		DB:       db,
		Order:    order.Init(db),
		Product:  product.Init(db),
		Customer: customer.Init(db),
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
