package validators

import (
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"unicode"
)

type Validator interface {
	Check()
	IsValid() bool
	Errors() []string
}

type UserValidator struct {
	User   *model.User
	errors []string
}

func (v *UserValidator) Check() {

	if v.User.Username == "" {
		//v.errors = append(v.errors, InvalidFormatOfUsername)
		panic(InvalidFormatOfUsername)
	}
	if err := server.DB.Where("username", v.User.Username).Find(v.User); err == nil {
		v.errors = append(v.errors, ExistUsername)
	}
	if !ValidPassword(v.User.Password) {
		v.errors = append(v.errors, InvalidFormatOfPassword)
	}

}

func (v *UserValidator) IsValid() bool {
	return len(v.errors) == 0
}

func (v *UserValidator) Errors() []string {
	return v.errors
}

func ValidPassword(s string) bool {
	if len(s) == 0 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasSymbol := false
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
		if unicode.IsSymbol(ch) {
			hasSymbol = true
		}
	}
	return hasSymbol && hasUpper && hasLower && hasDigit
}

func GetUserByUsername(username string) (*model.User, error) {

	user := model.User{}
	server.DB.Where("username", username).Find(&user)
	//json.NewEncoder(w).Encode(user)

	if user.Username == "" {
		return nil, http.ErrAbortHandler
	}

	return &user, nil
}
