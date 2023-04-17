package model

import (
	"gorm.io/gorm"
	"sdu.store/server"
)

type Comment struct {
	gorm.Model
	ID        int
	ProductID int
	UserID    int
	Text      string `gorm:"text"`
}

func (this *Comment) Create() (err error) {
	defer func() {
		if err != nil {
			return
		}
		var count int64
		if err = server.DB.Raw(
			"SELECT Count(*) FROM PRODUCTS, COMMENTS WHERE COMMENTS.PRODUCT_ID = PRODUCTS.ID AND PRODUCTS.ID = ?",
			this.ProductID,
		).Scan(&count).Error; err != nil {
			return
		}
		err = server.DB.Model(&Product{}).Where("ID = ?", this.ProductID).Update(
			"amount_comments", count,
		).Error
	}()
	err = server.DB.Create(this).Error
	return
}

func (this *Comment) Update() error {
	return server.DB.Model(this).Where("PRODUCT_ID = ? AND USER_ID = ?", this.ProductID, this.UserID).Update(
		"Text", this.Text,
	).Error
}

func (this *Comment) Delete() error {
	return server.DB.Delete(this).Error
}
