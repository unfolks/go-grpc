package product

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("product not found")
var ErrInvalidPrice = errors.New("invalid price")

type Product struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Price     float64    `json:"price"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	UpdatedBy *string    `json:"updated_by,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	DeletedBy *string    `json:"deleted_by,omitempty"`
}
