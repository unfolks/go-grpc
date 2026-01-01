package product

import "context"

type Repository interface {
	Save(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id string) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string, deletedBy string) error
	FindAll(ctx context.Context) ([]Product, error)
}
