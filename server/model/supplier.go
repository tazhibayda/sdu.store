package model

import "time"

type Supplier struct {
	UserId    int64     `json:"user_id"`
	ProductId int64     `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
}
