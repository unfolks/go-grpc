package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	product "hex-postgres-grpc/internal/product/domain"
)

type ProductRepoPG struct {
	db *sql.DB
}

func NewProductRepoPG(db *sql.DB) *ProductRepoPG {
	return &ProductRepoPG{db: db}
}

func (r *ProductRepoPG) Save(ctx context.Context, p *product.Product) error {
	const q = `INSERT INTO products (id, name, price, created_at, created_by) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, q, p.ID, p.Name, p.Price, p.CreatedAt, p.CreatedBy)
	return err
}

func (r *ProductRepoPG) FindByID(ctx context.Context, id string) (*product.Product, error) {
	const query = `SELECT id, name, price, created_at, created_by, updated_at, updated_by FROM products WHERE id = $1 AND deleted_at IS NULL`
	var p product.Product
	var created time.Time
	row := r.db.QueryRowContext(ctx, query, id)

	if err := row.Scan(&p.ID, &p.Name, &p.Price, &created, &p.CreatedBy, &p.UpdatedAt, &p.UpdatedBy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, product.ErrNotFound
		}
		return nil, err
	}
	p.CreatedAt = created
	return &p, nil
}

func (r *ProductRepoPG) Update(ctx context.Context, p *product.Product) error {
	const q = `UPDATE products SET name = $1, price = $2, updated_at = $3, updated_by = $4 WHERE id = $5 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, p.Name, p.Price, p.UpdatedAt, p.UpdatedBy, p.ID)
	return err
}

func (r *ProductRepoPG) Delete(ctx context.Context, id string, deletedBy string) error {
	const q = `UPDATE products SET deleted_at = $1, deleted_by = $2 WHERE id = $3 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, time.Now(), deletedBy, id)
	return err
}

func (r *ProductRepoPG) FindAllPaginated(ctx context.Context, limit, offset int) ([]product.Product, int, error) {
	const countQ = `SELECT COUNT(*) FROM products WHERE deleted_at IS NULL`
	var total int
	if err := r.db.QueryRowContext(ctx, countQ).Scan(&total); err != nil {
		return nil, 0, err
	}

	const q = `SELECT id, name, price, created_at, created_by, updated_at, updated_by FROM products WHERE deleted_at IS NULL LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := []product.Product{}
	for rows.Next() {
		var p product.Product
		var created time.Time
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &created, &p.CreatedBy, &p.UpdatedAt, &p.UpdatedBy); err != nil {
			return nil, 0, err
		}
		p.CreatedAt = created
		products = append(products, p)
	}
	return products, total, rows.Err()
}
