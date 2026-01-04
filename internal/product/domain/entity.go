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

type PaginatedData struct {
	Data             []Product `json:"data"`
	CurrentPage      int       `json:"current_page"`
	HaveNextPage     bool      `json:"have_next_page"`
	HavePreviousPage bool      `json:"have_previous_page"`
	Limit            int       `json:"limit"`
	TotalItem        int       `json:"total_item"`
	TotalPage        int       `json:"total_page"`
}

type PaginatedResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    PaginatedData `json:"data"`
}
