package usecase

import (
	"context"
	"hex-postgres-grpc/internal/category/domain"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) domain.Service {
	return &service{repo: repo}
}

func (s *service) CreateCategory(ctx context.Context, name, userID string) (*domain.Category, error) {
	id := uuid.NewString()
	category := domain.Category{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	err := s.repo.Save(ctx, &category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *service) GetCategory(ctx context.Context, id string) (*domain.Category, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) UpdateCategory(ctx context.Context, id, name, userID string) (*domain.Category, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	category.Name = name
	category.UpdatedAt = time.Now()
	category.UpdatedBy = userID

	err = s.repo.Update(ctx, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *service) DeleteCategory(ctx context.Context, id, userID string) error {
	// For now, doing hard delete as repo.Delete implies it.
	// If soft delete is needed, we would update DeletedAt/DeletedBy.
	return s.repo.Delete(ctx, id)
}

func (s *service) ListCategories(ctx context.Context) ([]*domain.Category, error) {
	return s.repo.FindAll(ctx)
}
