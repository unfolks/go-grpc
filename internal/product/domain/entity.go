package product

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("product not found")
var ErrInvalidPrice = errors.New("invalid price")

type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}
