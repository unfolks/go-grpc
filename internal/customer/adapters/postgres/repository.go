package postgres

import (
	"context"
	"database/sql"
	"hex-postgres-grpc/internal/customer/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, c *domain.Customer) error {
	query := "INSERT INTO customers (id, name, email, address, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.ExecContext(ctx, query, c.ID, c.Name, c.Email, c.Address, c.CreatedAt)

	// rows, _ := res.RowsAffected()
	// fmt.Println(rows)
	return err
}

func (r *Repository) FindByID(ctx context.Context, id string) (*domain.Customer, error) {
	query := "SELECT id, name, email, address, created_at FROM customers WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	var c domain.Customer
	err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Address, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // or specific not found error
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]domain.Customer, error) {
	query := "SELECT id, name, email, address, created_at FROM customers"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []domain.Customer{}
	for rows.Next() {
		var c domain.Customer
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Address, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}
