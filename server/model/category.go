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

func GetCategoryByID(id int) Category {
	var category Category
	server.DB.Where("ID=?", id).Find(&category)
	return category
}

func (category *Category) Delete() {
	server.DB.Where("ID=?", category.ID).Delete(&Category{})
}

func (category *Category) Update() {
	name := category.Name
	server.DB.First(category)
	category.Name = name
	server.DB.Save(category)
}
