package usecase

import (
	"context"
	"errors"
	"hex-postgres-grpc/internal/customer/domain"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) domain.Service {
	return &service{repo: repo}
}

func (s *service) CreateCustomer(ctx context.Context, name, email, address string) (*domain.Customer, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if address == "" {
		return nil, errors.New("address is required")
	}

	id := uuid.NewString()
	customer := domain.Customer{
		ID:        id,
		Name:      name,
		Email:     email,
		Address:   address,
		CreatedAt: time.Now(),
	}
	err := s.repo.Save(ctx, &customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (s *service) ListCustomers(ctx context.Context) ([]domain.Customer, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return s.repo.FindByID(ctx, id)
}
