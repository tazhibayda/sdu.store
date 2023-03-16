package model

import (
	"sdu.store/server"
)

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func ConfigCategories() {
	server.DB.Create(&Category{Name: "Hoodies"})
	server.DB.Create(&Category{Name: "Caps"})
	server.DB.Create(&Category{Name: "T-shirts"})
}
