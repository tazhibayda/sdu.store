package model

import (
	"sdu.store/server"
)

type Category struct {
	ID       int
	Name     string    `json:"name"`
	Products []Product `gorm:"constraint:OnDelete:CASCADE;"`
}

func ConfigCategories() {
	server.DB.Create(&Category{Name: "Shopperler"})
	server.DB.Create(&Category{Name: "Jeideler"})
	server.DB.Create(&Category{Name: "Burjeideler"})
	server.DB.Create(&Category{Name: "Qosymshalar"})
}

func GetAllCategory() (categories []Category, err error) {
	err = server.DB.Find(&categories).Error
	return
}

func GetCategoryByID(id int) (Category, error) {
	var category Category
	if err := server.DB.Where("ID=?", id).Preload("Products").Find(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (category *Category) Delete() error {
	return server.DB.Where("ID=?", category.ID).Unscoped().Delete(&Category{}).Error
}

func (category *Category) Update() error {
	return server.DB.Model(category).Update("name", category.Name).Error
}
