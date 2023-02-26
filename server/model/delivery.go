package model

import "time"

type Delivery struct {
	ID          int64     `json:"id"`
	OrderID     int64     `json:"order_id"`
	Address     string    `json:"address"`
	PhoneNumber string    `json:"phone_number"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
