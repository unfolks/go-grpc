package usecase

import (
	"context"
	"errors"
	"hex-postgres-grpc/internal/customer/domain"
	"time"

	"github.com/google/uuid"
)

type CreateCustomer struct {
	repo domain.Repository
}

func NewCreateCustomer(repo domain.Repository) *CreateCustomer {
	return &CreateCustomer{repo: repo}
}

func (u *CreateCustomer) Execute(ctx context.Context, name, email, address string) (*domain.Customer, error) {
	// validation
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
	err := u.repo.Save(ctx, &customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}
