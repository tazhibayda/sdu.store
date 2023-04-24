package validators

import "sdu.store/server/model"

type ProductValidator struct {
	*model.Product
	Validator
}

func NewProductValidator(product *model.Product) ProductValidator {
	return ProductValidator{product, Validator{errors: []string{}}}
}

func (v *ProductValidator) Check() {
	if v.Name == "" {
		v.errors = append(v.errors, InvalidFormatOfName)
	}
	if _, err := model.GetCategoryByID(v.CategoryID); err != nil {
		v.errors = append(v.errors, DoesNotExistCategory)
	}
	if v.Price < 0 {
		v.errors = append(v.errors, InvalidPrice)
	}
	if len(v.Images) == 0 {
		v.errors = append(v.errors, ImageFormat)
	}
	if len(v.Colors) == 0 {
		v.errors = append(v.errors, ColorFormat)
	}
	if v.Description == "" {
		v.errors = append(v.errors, DescriptionFormat)
	}
}
