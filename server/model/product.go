package model

import "time"

type Product struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	CategoryID int64     `json:"category_id"`
	Price      float64   `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
}
