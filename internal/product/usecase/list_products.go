package usecase

import (
	"context"
	product "hex-postgres-grpc/internal/product/domain"
)

type ListProducts struct {
	repo product.Repository
}

func NewListProducts(repo product.Repository) *ListProducts {
	return &ListProducts{repo: repo}
}

func (u *ListProducts) Execute(ctx context.Context) ([]product.Product, error) {
	return u.repo.FindAll(ctx)
}
