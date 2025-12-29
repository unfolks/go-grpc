package product

import (
	"context"
)

type Service interface {
	CreateProduct(ctx context.Context, name string, price float64) (Product, error)
	GetProduct(ctx context.Context, id string) (Product, error)
	UpdateProduct(ctx context.Context, id string, name string, price float64) (Product, error)
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context) ([]Product, error)
}
