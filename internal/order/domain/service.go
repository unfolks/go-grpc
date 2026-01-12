package order

import (
	"context"
	"errors"
	domain_common "hex-postgres-grpc/internal/common/domain"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("order not found")
var ErrInvalidAmount = errors.New("invalid amount")

type Service interface {
	CreateOrder(ctx context.Context, amount float64) (Order, error)
	GetOrder(ctx context.Context, id string) (Order, error)
	UpdateOrder(ctx context.Context, id string, amount float64) (Order, error)
	DeleteOrder(ctx context.Context, id string) error
	ListOrders(ctx context.Context) ([]Order, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateOrder(ctx context.Context, amount float64) (Order, error) {
	if amount <= 0 {
		return Order{}, ErrInvalidAmount
	}

	id := uuid.NewString()
	o := Order{
		BaseEntity: domain_common.BaseEntity{
			ID:        id,
			CreatedAt: time.Now(),
		},
		Amount: amount,
	}
	if err := s.repo.Save(ctx, &o); err != nil {
		return Order{}, err
	}
	return o, nil

}
func (s *service) GetOrder(ctx context.Context, id string) (Order, error) {
	o, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return Order{}, err
	}
	return *o, nil
}

func (s *service) UpdateOrder(ctx context.Context, id string, amount float64) (Order, error) {
	o, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return Order{}, err
	}

	if amount <= 0 {
		return Order{}, ErrInvalidAmount
	}

	o.Amount = amount
	if err := s.repo.Update(ctx, o); err != nil {
		return Order{}, err
	}

	return *o, nil
}

func (s *service) DeleteOrder(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) ListOrders(ctx context.Context) ([]Order, error) {
	return s.repo.FindAll(ctx)
}
