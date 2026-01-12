package order

import (
	domain_common "hex-postgres-grpc/internal/common/domain"
)

type Order struct {
	domain_common.BaseEntity
	Amount float64 `json:"amount"`
}
