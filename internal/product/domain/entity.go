package product

import (
	"errors"
	domain_common "hex-postgres-grpc/internal/common/domain"
)

var ErrNotFound = errors.New("product not found")
var ErrInvalidPrice = errors.New("invalid price")

type Product struct {
	domain_common.BaseEntity
	Name  string  `json:"name"`
	Price float64 `json:"price"`
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
