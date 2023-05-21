package model

import (
	"gorm.io/gorm"
	"sdu.store/server"
)

type Purchase struct {
	gorm.Model
	ItemID int
	UserID int
}

func MakePurchase(productID int, color string, size string, userID int) error {
	var item = &Item{}
	if err := server.DB.Model(item).Where("PRODUCT_ID = ? AND IS_SOLD = FALSE AND SIZE = ? AND COLOR = ?", productID, size, color).Find(item).Error; err != nil {
		return err
	}
	item.IsSold = true
	purchase := Purchase{ItemID: int(item.ID), UserID: userID}
	if err := server.DB.Create(&purchase).Error; err != nil {
		return err
	}
	return server.DB.Save(item).Error
}

func GetAllPurchases() []Purchase {
	var purchases []Purchase
	server.DB.Find(&purchases)
	return purchases
}

func GetAllPurchaseByUserID(userID int) []Item {
	var purchases []Purchase
	server.DB.Find(purchases, "USER_ID = ?", userID).Preload("Item")
	var items []Item
	for _, purchase := range purchases {
		var item = &Item{}
		server.DB.Find(item, "ID = ?", purchase.ItemID)
		items = append(items, *item)
	}
	return items
}
