package model

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	ID        int64 `json:"id"`
	ProductID uint
	Color     string `json:"color"`
	Size      string `json:"size"`
	Quantity  int64  `json:"quantity"`
}
