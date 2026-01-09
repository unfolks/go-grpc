package postgres

import (
	"context"
	"database/sql"

	"hex-postgres-grpc/internal/category/domain"
)

type CategoryRepoPG struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *CategoryRepoPG {
	return &CategoryRepoPG{db: db}
}
func (c *CategoryRepoPG) Save(ctx context.Context, p *domain.Category) error {
	const q = `INSERT INTO category (id, name, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := c.db.ExecContext(ctx, q, p.ID, p.Name, p.CreatedAt, p.UpdatedAt, p.CreatedBy, p.UpdatedBy)
	return err
}

func (c *CategoryRepoPG) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	const q = `SELECT id, name, created_at, updated_at, created_by, updated_by FROM category WHERE id = $1 AND deleted_at IS NULL`
	var cat domain.Category
	err := c.db.QueryRowContext(ctx, q, id).Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt, &cat.CreatedBy, &cat.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (c *CategoryRepoPG) Update(ctx context.Context, p *domain.Category) error {
	const q = `UPDATE category SET name = $1, updated_at = $2, updated_by = $3 WHERE id = $4 AND deleted_at IS NULL`
	_, err := c.db.ExecContext(ctx, q, p.Name, p.UpdatedAt, p.UpdatedBy, p.ID)
	return err
}

func (c *CategoryRepoPG) Delete(ctx context.Context, id string) error {
	const q = `UPDATE category SET deleted_at = NOW() WHERE id = $1`
	_, err := c.db.ExecContext(ctx, q, id)
	return err
}

func (c *CategoryRepoPG) FindAll(ctx context.Context) ([]*domain.Category, error) {
	const q = `SELECT id, name, created_at, updated_at, created_by, updated_by FROM category WHERE deleted_at IS NULL`
	rows, err := c.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var cat domain.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt, &cat.CreatedBy, &cat.UpdatedBy); err != nil {
			return nil, err
		}
		categories = append(categories, &cat)
	}
	return categories, nil
}
