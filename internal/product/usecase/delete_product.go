package usecase

import (
	"context"
	product "hex-postgres-grpc/internal/product/domain"
)

type DeleteProduct struct {
	repo product.Repository
}

func NewDeleteProduct(repo product.Repository) *DeleteProduct {
	return &DeleteProduct{repo: repo}
}

func (u *DeleteProduct) Execute(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
