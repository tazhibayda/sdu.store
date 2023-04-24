package validators

import (
	"sdu.store/server"
	"sdu.store/server/model"
)

type CategoryValidator struct {
	*model.Category
	Validator
}

func NewCategoryValidator(category *model.Category) CategoryValidator {
	return CategoryValidator{category, Validator{errors: []string{}}}
}

func (v *CategoryValidator) Check() {
	if v.Name == "" {
		v.errors = append(v.errors, InvalidFormatOfName)
	}
	var category model.Category
	if err := server.DB.Where("name = ? and id != ?", v.Name, v.ID).First(&category).Error; err == nil {
		v.errors = append(v.errors, RepeatedName)
	}
}
