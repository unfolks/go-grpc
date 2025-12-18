package usecase

import (
	"context"
	product "hex-postgres-grpc/internal/product/domain"
)

type UpdateProduct struct {
	repo product.Repository
}

func NewUpdateProduct(repo product.Repository) *UpdateProduct {
	return &UpdateProduct{repo: repo}
}

func (u *UpdateProduct) Execute(ctx context.Context, id string, name string, price float64) (product.Product, error) {
	p, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return product.Product{}, err
	}

	if price < 0 {
		return product.Product{}, product.ErrInvalidPrice
	}

	p.Name = name
	p.Price = price

	if err := u.repo.Update(ctx, p); err != nil {
		return product.Product{}, err
	}

	return *p, nil
}
