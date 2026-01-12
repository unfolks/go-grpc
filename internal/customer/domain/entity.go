package domain

import (
	domain_common "hex-postgres-grpc/internal/common/domain"
)

type Customer struct {
	domain_common.BaseEntity
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}
