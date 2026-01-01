package usecase

import (
	"context"
	"time"

	"hex-postgres-grpc/internal/auth"
	product "hex-postgres-grpc/internal/product/domain"

	"github.com/google/uuid"
)

const SystemUserID = "00000000-0000-0000-0000-000000000000"

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

	sub, _ := auth.SubjectFromContext(ctx)
	createdBy := SystemUserID
	if sub.ID != "" {
		createdBy = sub.ID
	}

	id := uuid.NewString()
	p := product.Product{
		ID:        id,
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
		CreatedBy: createdBy,
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

	sub, _ := auth.SubjectFromContext(ctx)
	updatedBy := SystemUserID
	if sub.ID != "" {
		updatedBy = sub.ID
	}

	now := time.Now()
	p.Name = name
	p.Price = price
	p.UpdatedAt = &now
	p.UpdatedBy = &updatedBy

	if err := s.repo.Update(ctx, p); err != nil {
		return product.Product{}, err
	}

	return *p, nil
}

func (s *service) DeleteProduct(ctx context.Context, id string) error {
	sub, _ := auth.SubjectFromContext(ctx)
	deletedBy := SystemUserID
	if sub.ID != "" {
		deletedBy = sub.ID
	}
	return s.repo.Delete(ctx, id, deletedBy)
}

func (s *service) ListProducts(ctx context.Context) ([]product.Product, error) {
	return s.repo.FindAll(ctx)
}
