package model

import (
	"gorm.io/gorm"
	"sdu.store/server"
)

type Rating struct {
	gorm.Model
	UserID    int
	ProductID int
	Value     int
}

func (this *Rating) Create() error {
	var exists bool
	if server.DB.Model(this).Select("count(*) > 0").Where(
		"product_id = ? and user_id = ?", this.ProductID, this.UserID,
	).Find(&exists); exists {
		return this.Update()
	}

	return server.DB.Create(this).Error
}

func (this *Rating) Update() error {
	return server.DB.Model(this).Where("ID = ?", this.ID).Update("Value", this.Value).Error
}

func (this *Rating) Delete() error {
	return server.DB.Delete(this).Error
}
