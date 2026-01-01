package domain

import "context"

type Service interface {
	CreateCustomer(ctx context.Context, name, email, address string) (*Customer, error)
	ListCustomers(ctx context.Context) ([]Customer, error)
	GetCustomer(ctx context.Context, id string) (*Customer, error)
}
