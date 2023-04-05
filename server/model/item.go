package model

import (
	"gorm.io/gorm"
	"sdu.store/server"
)

type Item struct {
	gorm.Model
	ID        int64  `json:"id"`
	ProductID uint   `json:"product_id"`
	Color     string `json:"color"`
	Size      string `json:"size"`
	Quantity  int64  `json:"quantity"`
}

func (this *Item) CreateOrUpdate() error {
	var item Item
	if err := server.DB.Where(
		"product_id = ? AND COLOR = ? AND SIZE = ?", this.ProductID, this.Color, this.Size,
	).First(&item).Error; err != nil {
		item.ProductID = this.ProductID
		item.Color = this.Color
		item.Size = this.Size
		item.Quantity = this.Quantity
		return item.Create()
	}
	item.Quantity += this.Quantity
	*this = item
	return item.Update()
}

//func (this *Item) DeleteOrUpdate() error{
//	var item Item
//	if err := server.DB.Where("PRODUCT_ID = ? AND COLOR = ? AND SIZE = ? AND QUANTITY 0",  this.ProductID, this.Color, this.Size).First(&item).Error; err != nil{
//
//	}
//}

func (item *Item) Create() error {
	return server.DB.Create(item).Error
}

func (item *Item) Update() error {
	return server.DB.Model(item).Updates(item).Error
}
