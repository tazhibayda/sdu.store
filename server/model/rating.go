package model

import (
	"fmt"
	"gorm.io/gorm"
	"sdu.store/server"
)

type Rating struct {
	gorm.Model
	UserID    int
	ProductID int
	Value     int
}

func (this *Rating) Create() (err error) {
	defer func() {
		if err != nil {
			return
		}
		var rating float64
		if err = server.DB.Raw(
			"SELECT AVG(VALUE) FROM PRODUCTS, RATINGS WHERE RATINGS.PRODUCT_ID = PRODUCTS.ID AND PRODUCTS.ID = ?",
			this.ProductID,
		).Scan(&rating).Error; err != nil {
			return
		}
		err = server.DB.Model(&Product{}).Where("ID = ?", this.ProductID).Update(
			"rating", fmt.Sprintf("%.2f", rating),
		).Error
		var count int64
		if err = server.DB.Raw(
			"SELECT Count(*) FROM PRODUCTS, RATINGS WHERE RATINGS.PRODUCT_ID = PRODUCTS.ID AND PRODUCTS.ID = ?",
			this.ProductID,
		).Scan(&count).Error; err != nil {
			return
		}
		err = server.DB.Model(&Product{}).Where("ID = ?", this.ProductID).Update(
			"amount_ratings", count,
		).Error
	}()
	var count int64
	if err = server.DB.Model(this).Where(
		"product_id = ? and user_id = ?", this.ProductID, this.UserID,
	).Count(&count).Error; count > 0 {
		err = this.Update()
		return
	}
	err = server.DB.Create(this).Error
	return
}

func (this *Rating) Update() error {
	return server.DB.Model(this).Where("PRODUCT_ID = ? AND USER_ID = ?", this.ProductID, this.UserID).Update(
		"Value", this.Value,
	).Error
}

func (this *Rating) Delete() error {
	return server.DB.Delete(this).Error
}
