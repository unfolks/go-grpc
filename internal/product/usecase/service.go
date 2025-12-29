package usecase

import (
	"context"
	"time"

	product "hex-postgres-grpc/internal/product/domain"

	"github.com/google/uuid"
)

type service struct {
	repo product.Repository
}

func NewService(repo product.Repository) product.Service {
	return &service{repo: repo}
}

func (s *service) CreateProduct(ctx context.Context, name string, price float64) (product.Product, error) {
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
	if err := s.repo.Save(ctx, &p); err != nil {
		return product.Product{}, err
	}
	return p, nil
}

func (s *service) GetProduct(ctx context.Context, id string) (product.Product, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return product.Product{}, err
	}
	return *p, nil
}

func (s *service) UpdateProduct(ctx context.Context, id string, name string, price float64) (product.Product, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return product.Product{}, err
	}

	if price < 0 {
		return product.Product{}, product.ErrInvalidPrice
	}

	p.Name = name
	p.Price = price

	if err := s.repo.Update(ctx, p); err != nil {
		return product.Product{}, err
	}

	return *p, nil
}

func (s *service) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) ListProducts(ctx context.Context) ([]product.Product, error) {
	return s.repo.FindAll(ctx)
}
