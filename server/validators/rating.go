package validators

import "sdu.store/server/model"

type RatingValidator struct {
	*model.Rating
	Validator
}

func (v *RatingValidator) Check() {
	if v.Value < 1 || v.Value > 5 {
		v.errors = append(v.errors, "Rating value should be in range from 1 till 5")
	}
	if _, err := model.GetProductByID(v.ProductID); err != nil {
		v.errors = append(v.errors, "Product doesn't exist")
	}
}
