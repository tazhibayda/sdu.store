package validators

import "sdu.store/server/model"

type ItemValidator struct {
	Validator
	*model.Item
}

func NewItemValidator(item *model.Item) ItemValidator {
	return ItemValidator{Validator{errors: []string{}}, item}
}

func (v *ItemValidator) Check() {
	if _, err := model.GetItemByBarcode(v.Barcode); err == nil {
		v.errors = append(v.errors, AlreadyExistItemWithBarcode)
	}
	if len(v.Color) == 0 {
		v.errors = append(v.errors, ColorFormat)
	}
	if len(v.Size) == 0 {
		v.errors = append(v.errors, ColorFormat)
	}
}
