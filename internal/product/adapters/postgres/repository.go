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
	const q = `INSERT INTO products (id, name, price, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, q, p.ID, p.Name, p.Price, p.CreatedAt)
	return err
}

func (r *ProductRepoPG) FindByID(ctx context.Context, id string) (*product.Product, error) {
	const query = `SELECT id, name, price, created_at FROM products WHERE id = $1`
	var p product.Product
	var created time.Time
	row := r.db.QueryRowContext(ctx, query, id)

	if err := row.Scan(&p.ID, &p.Name, &p.Price, &created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, product.ErrNotFound
		}
		return nil, err
	}
	p.CreatedAt = created
	return &p, nil
}

func (r *ProductRepoPG) Update(ctx context.Context, p *product.Product) error {
	const q = `UPDATE products SET name = $1, price = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, q, p.Name, p.Price, p.ID)
	return err
}

func (r *ProductRepoPG) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM products WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

func (r *ProductRepoPG) FindAll(ctx context.Context) ([]product.Product, error) {
	const q = `SELECT id, name, price, created_at FROM products`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []product.Product
	for rows.Next() {
		var p product.Product
		var created time.Time
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &created); err != nil {
			return nil, err
		}
		p.CreatedAt = created
		products = append(products, p)
	}
	return products, rows.Err()
}
