package domain

import "context"

type Service interface {
	CreateCategory(ctx context.Context, name, userID string) (*Category, error)
	GetCategory(ctx context.Context, id string) (*Category, error)
	UpdateCategory(ctx context.Context, id, name, userID string) (*Category, error)
	DeleteCategory(ctx context.Context, id, userID string) error
	ListCategories(ctx context.Context) ([]*Category, error)
}
