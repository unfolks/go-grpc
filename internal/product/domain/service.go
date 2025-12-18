package product

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("product not found")
var ErrInvalidPrice = errors.New("invalid price")

type Service interface {
	CreateProduct(ctx context.Context, name string, price float64) (Product, error)
	GetProduct(ctx context.Context, id string) (Product, error)
	UpdateProduct(ctx context.Context, id string, name string, price float64) (Product, error)
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context) ([]Product, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateProduct(ctx context.Context, name string, price float64) (Product, error) {
	if price < 0 {
		return Product{}, ErrInvalidPrice
	}

	id := uuid.NewString()
	p := Product{
		ID:        id,
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Save(ctx, &p); err != nil {
		return Product{}, err
	}
	return p, nil
}

func (s *service) GetProduct(ctx context.Context, id string) (Product, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return Product{}, err
	}
	return *p, nil
}

func (s *service) UpdateProduct(ctx context.Context, id string, name string, price float64) (Product, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return Product{}, err
	}

	if price < 0 {
		return Product{}, ErrInvalidPrice
	}

	p.Name = name
	p.Price = price

	if err := s.repo.Update(ctx, p); err != nil {
		return Product{}, err
	}

	return *p, nil
}

func (s *service) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) ListProducts(ctx context.Context) ([]Product, error) {
	return s.repo.FindAll(ctx)
}
