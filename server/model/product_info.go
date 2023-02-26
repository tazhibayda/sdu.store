package model

import "time"

type ProductInfo struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
}
