package usecase

import (
	"context"
	"hex-postgres-grpc/internal/category/domain"
	domain_common "hex-postgres-grpc/internal/common/domain"
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
		BaseEntity: domain_common.BaseEntity{
			ID:        id,
			CreatedAt: time.Now(),
			CreatedBy: userID,
		},
		Name: name,
	}
	now := time.Now()
	category.UpdatedAt = &now
	category.UpdatedBy = &userID
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
	now := time.Now()
	category.UpdatedAt = &now
	category.UpdatedBy = &userID

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
