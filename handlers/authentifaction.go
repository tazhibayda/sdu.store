package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/server/validators"
	"sdu.store/utils"
	"time"
)

func Login(writer http.ResponseWriter, request *http.Request) {
	Username := request.FormValue("username")
	Password := request.FormValue("password")
	user, err := model.GetUserByUsername(Username)
	if err != nil {
		utils.ExecuteTemplateWithoutNavbar(
			writer, request, []string{"Password or username is incorrect"}, "templates/sign-in.html",
		)
		return
	}

	if user.Password == Password || utils.CheckPasswordHash(Password, user.Password) {
		if err := DoLogin(writer, *user); err != nil {
			utils.ServerErrorHandler(writer, request, err)
			return
		}
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}
	utils.ExecuteTemplateWithoutNavbar(
		writer, request, []string{"Password or username is incorrect"}, "templates/sign-in.html",
	)
}

func DoLogin(writer http.ResponseWriter, user model.User) error {
	expirationTime := time.Now().Add(24 * 60 * time.Minute)
	usr := &model.Claims{
		User: &user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, usr)
	tokenString, err := token.SignedString(model.JwtKey)
	if err != nil {
		return err
	}
	http.SetCookie(
		writer, &http.Cookie{
			Name:    "session_token",
			Value:   tokenString,
			Path:    "/",
			Expires: expirationTime,
		},
	)
	return nil
}

func Logout(writer http.ResponseWriter, request *http.Request) {

	_, err := utils.CheckCookie(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	c := &http.Cookie{
		Name:    "session_token",
		Expires: time.Now(),
	}

	http.SetCookie(writer, c)
	http.Redirect(writer, request, "/", http.StatusSeeOther)

}

func SignUp(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	if request.PostFormValue("password") == request.PostFormValue("repassword") {
		user := model.User{
			Username: request.PostFormValue("username"),
			Email:    request.PostFormValue("email"),
			Password: request.PostFormValue("password"),
		}
		v := validators.NewUserValidator(&user)
		if v.Check(); !v.IsValid() {
			utils.ExecuteTemplateWithoutNavbar(writer, request, v.Errors(), "templates/sign-up.html")
			return
		}
		if err := server.DB.Create(&user).Error; err != nil {
			utils.ServerErrorHandler(writer, request, err)
			return
		}
		if user.ID == 1 {
			user.Is_admin = true
			user.Is_staff = true
			server.DB.Save(&user)
		}
		http.Redirect(writer, request, "/login", 302)
		return
	}

	utils.ExecuteTemplateWithoutNavbar(
		writer, request, []string{"Two passwords doesn't match"}, "templates/sign-up.html",
	)
	return

}

func LoginPage(writer http.ResponseWriter, request *http.Request) {
	_, err := utils.Session(writer, request)
	if err == nil {
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}
	utils.ExecuteTemplateWithoutNavbar(writer, request, nil, "templates/sign-in.html")

}

func SignUpPage(writer http.ResponseWriter, request *http.Request) {
	utils.ExecuteTemplateWithoutNavbar(writer, request, nil, "templates/sign-up.html")
}
