package order

import "time"

type Order struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
