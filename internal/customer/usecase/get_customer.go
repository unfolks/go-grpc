package usecase

import (
	"context"
	customer "hex-postgres-grpc/internal/customer/domain"
)

type ListCustomer struct {
	repo customer.Repository
}

func NewListCustomer(repo customer.Repository) *ListCustomer {
	return &ListCustomer{repo: repo}
}

func (g *ListCustomer) Execute(ctx context.Context) ([]customer.Customer, error) {
	return g.repo.FindAll(ctx)
}
