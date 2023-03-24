package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ProductOutput struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Category  string   `json:"category"`
	Price     float64  `json:"price"`
	Images    []string `json:"images"`
	Sizes     []string `json:"sizes"`
	Colors    []string `json:"colors"`
	CreatedAt string   `json:"created_at"`
}

type Product struct {
	gorm.Model
	Name        string `json:"name"`
	CategoryID  uint
	Price       float64        `json:"price"`
	Images      pq.StringArray `gorm:"type:text[]" json:"images"`
	Sizes       pq.StringArray `gorm:"type:text[]" json:"sizes"`
	Colors      pq.StringArray `gorm:"type:text[]" json:"colors"`
	Description string         `json:"description"`
	Items       []Item         `gorm:"foreignKey:ID"`
}
