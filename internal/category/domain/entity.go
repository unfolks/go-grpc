package domain

import (
	"hex-postgres-grpc/internal/common/domain"
)

type Category struct {
	domain.BaseEntity
	Name string `json:"name"`
}
