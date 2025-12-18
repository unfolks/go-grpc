package usecase

import (
	"context"
	product "hex-postgres-grpc/internal/product/domain"
)

type GetProduct struct {
	repo product.Repository
}

func NewGetProduct(repo product.Repository) *GetProduct {
	return &GetProduct{repo: repo}
}

func (u *GetProduct) Execute(ctx context.Context, id string) (product.Product, error) {
	p, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return product.Product{}, err
	}
	return *p, nil
}
