package model

import (
	"gorm.io/gorm"
	"sdu.store/server"
)

type Item struct {
	gorm.Model
	Barcode   string `json:"id" gorm:"primaryKey"`
	IsSold    bool   `json:"is-sold"`
	ProductID uint   `json:"product_id"`
	Color     string `json:"color"`
	Size      string `json:"size"`
}

//func (this *Item) DeleteOrUpdate() error{
//	var item Item
//	if err := server.DB.Where("PRODUCT_ID = ? AND COLOR = ? AND SIZE = ? AND QUANTITY 0",  this.ProductID, this.Color, this.Size).First(&item).Error; err != nil{
//
//	}
//}

func GetAllItems() (items []Item, err error) {
	err = server.DB.Find(&items).Error
	return
}

func GetNotSoldItemBySizeAndColorAndProductID(color, size string, productID int) (*Item, error) {
	var item = &Item{}
	err := server.DB.Find(item, "COLOR = ? AND SIZE = ? AND IS_SOLD = FALSE AND PRODUCT_ID = ?", color, size, productID).Error
	return item, err
}

func GetItemByBarcode(barcode string) (*Item, error) {
	item := &Item{}
	err := server.DB.Where("BARCODE = ?", barcode).Find(item).Error
	return item, err
}

func (item *Item) Create() error {
	return server.DB.Create(item).Error
}

func (item *Item) Update() error {
	return server.DB.Model(item).Updates(item).Error
}
