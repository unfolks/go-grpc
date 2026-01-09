package domain

import "context"

type Repository interface {
	Save(ctx context.Context, category *Category) error
	FindByID(ctx context.Context, id string) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*Category, error)
}
