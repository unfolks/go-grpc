package usecase

import (
	"context"
	"time"

	product "hex-postgres-grpc/internal/product/domain"

	"github.com/google/uuid"
)

type CreateProduct struct {
	repo product.Repository
}

func NewCreateProduct(repo product.Repository) *CreateProduct {
	return &CreateProduct{repo: repo}
}

func (u *CreateProduct) Execute(ctx context.Context, name string, price float64) (product.Product, error) {
	if price < 0 {
		return product.Product{}, product.ErrInvalidPrice
	}

	id := uuid.NewString()
	p := product.Product{
		ID:        id,
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
	}
	if err := u.repo.Save(ctx, &p); err != nil {
		return product.Product{}, err
	}
	return p, nil
}
