package model

import (
	"gorm.io/gorm"
	"sdu.store/server"
)

type Category struct {
	gorm.Model
	Name     string    `json:"name"`
	Products []Product `gorm:"foreignKey:ID"`
}

func ConfigCategories() {
	server.DB.Create(&Category{Name: "Hoodies"})
	server.DB.Create(&Category{Name: "Caps"})
	server.DB.Create(&Category{Name: "T-shirts"})
}

func GetCategoryByID(id int) (Category, error) {
	var category Category
	if err := server.DB.Where("ID=?", id).Find(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (category *Category) Delete() error {
	return server.DB.Where("ID=?", category.ID).Delete(&Category{}).Error
}

func (category *Category) Update() error {
	name := category.Name
	if err := server.DB.First(category).Error; err != nil {
		return err
	}
	category.Name = name
	return server.DB.Save(category).Error
}
