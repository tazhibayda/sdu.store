package model

import "time"

type ProductInfo struct {
	ID        int64 `json:"id"`
	ProductID int64 `json:"product_id"`
	// forgot text will add in next update
	//Information string `json:"information"`
	CreatedAt time.Time `json:"created_at"`
}
