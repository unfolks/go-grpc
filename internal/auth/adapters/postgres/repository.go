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
