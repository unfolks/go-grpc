package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"hex-postgres-grpc/internal/auth"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) auth.UserRepository {
	return &repository{db: db}
}

func (r *repository) GetByUsername(ctx context.Context, username string) (*auth.User, error) {
	query := `SELECT id, username, password_hash, role, attributes FROM users WHERE username = $1`
	var user auth.User
	var attrJSON []byte
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &attrJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	if err := json.Unmarshal(attrJSON, &user.Attributes); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*auth.User, error) {
	query := `SELECT id, username, password_hash, role, attributes FROM users WHERE id = $1`
	var user auth.User
	var attrJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &attrJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	if err := json.Unmarshal(attrJSON, &user.Attributes); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) Create(ctx context.Context, user *auth.User) error {
	attrJSON, err := json.Marshal(user.Attributes)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (id, username, password_hash, role, attributes) VALUES ($1, $2, $3, $4, $5)`
	_, err = r.db.ExecContext(ctx, query, user.ID, user.Username, user.PasswordHash, user.Role, attrJSON)
	return err
}

func (r *repository) Update(ctx context.Context, user *auth.User) error {
	attrJSON, err := json.Marshal(user.Attributes)
	if err != nil {
		return err
	}

	query := `UPDATE users SET username = $1, password_hash = $2, role = $3, attributes = $4 WHERE id = $5`
	res, err := r.db.ExecContext(ctx, query, user.Username, user.PasswordHash, user.Role, attrJSON, user.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *repository) List(ctx context.Context) ([]*auth.User, error) {
	query := `SELECT id, username, password_hash, role, attributes FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*auth.User
	for rows.Next() {
		var user auth.User
		var attrJSON []byte
		if err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &attrJSON); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(attrJSON, &user.Attributes); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
