package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	order "hex-postgres-grpc/internal/order/domain"
)

type OrderRepoPG struct {
	db *sql.DB
}

func NewOrderRepoPG(db *sql.DB) *OrderRepoPG {
	return &OrderRepoPG{db: db}
}

func (r *OrderRepoPG) Save(ctx context.Context, o *order.Order) error {
	const q = `INSERT INTO orders (id, amount, created_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, q, o.ID, o.Amount, o.CreatedAt)
	return err
}

func (r *OrderRepoPG) FindByID(ctx context.Context, id string) (*order.Order, error) {
	const query = `SELECT id, amount, created_at FROM orders WHERE id = $1`
	var o order.Order
	var created time.Time
	row := r.db.QueryRowContext(ctx, query, id)

	if err := row.Scan(&o.ID, &o.Amount, &created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, order.ErrNotFound
		}
		return nil, err
	}
	o.CreatedAt = created
	return &o, nil
}

func (r *OrderRepoPG) Update(ctx context.Context, o *order.Order) error {
	const q = `UPDATE orders SET amount = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, q, o.Amount, o.ID)
	return err
}

func (r *OrderRepoPG) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM orders WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

func (r *OrderRepoPG) FindAll(ctx context.Context) ([]order.Order, error) {
	const q = `SELECT id, amount, created_at FROM orders`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []order.Order{}
	for rows.Next() {
		var o order.Order
		var created time.Time
		if err := rows.Scan(&o.ID, &o.Amount, &created); err != nil {
			return nil, err
		}
		o.CreatedAt = created
		orders = append(orders, o)
	}
	return orders, rows.Err()
}
