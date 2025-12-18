package domain

import "context"

type Repository interface {
	Save(ctx context.Context, customer *Customer) error
	FindByID(ctx context.Context, id string) (*Customer, error)
	FindAll(ctx context.Context) ([]Customer, error)
}
