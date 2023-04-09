package validators

import (
	"sdu.store/server"
	"sdu.store/server/model"
	"unicode"
)

type UserValidator struct {
	User *model.User
	Validator
}

func (v *UserValidator) Check() {

	if v.User.Username == "" {
		v.errors = append(v.errors, InvalidFormatOfName)
	}
	if server.DB.Where("username", v.User.Username).Find(v.User); v.User.ID != 0 {
		v.errors = append(v.errors, ExistUsername)
	}
	if v.User.Email == "" {
		v.errors = append(v.errors, InvalidFormatOfEmail)
	}
	if server.DB.Where("email", v.User.Email).Find(v.User); v.User.ID != 0 {
		v.errors = append(v.errors, ExistEmail)
	}
	if !ValidPassword(v.User.Password) {
		v.errors = append(v.errors, InvalidFormatOfPassword)
	}

}

func ValidPassword(s string) bool {
	if len(s) == 0 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	for _, ch := range s {
		if unicode.IsDigit(ch) {
			hasDigit = true
		}
		if unicode.IsLower(ch) {
			hasLower = true
		}
		if unicode.IsUpper(ch) {
			hasUpper = true
		}
	}
	return (hasUpper && hasLower && hasDigit)
}
