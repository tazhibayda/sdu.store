package validators

import (
	"sdu.store/server/model"
)

type PurchaseValidator struct {
	Validator
	Color     string
	Size      string
	ProductID int
}

func NewPurchaseValidator(color, size string, productID int) PurchaseValidator {
	return PurchaseValidator{Size: size, Color: color, ProductID: productID, Validator: Validator{errors: []string{}}}
}

func (v *PurchaseValidator) Check() {
	item, err := model.GetNotSoldItemBySizeAndColorAndProductID(v.Color, v.Size, v.ProductID)

	if err != nil || item.ID == 0 {
		v.errors = append(v.errors, NotExistItem)
	}
}
