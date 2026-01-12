package domain

import "time"

type BaseEntity struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedBy string     `json:"created_by,omitempty"`
	UpdatedBy *string    `json:"updated_by,omitempty"`
	DeletedBy *string    `json:"deleted_by,omitempty"`
}
