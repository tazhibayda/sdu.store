package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/server/validators"
	"sdu.store/utils"
	"time"
)

func Login(writer http.ResponseWriter, request *http.Request) {

	if request.Method != "POST" {
		http.Redirect(writer, request, "/login-page", http.StatusMethodNotAllowed)
		return
	}

	Username := request.PostFormValue("username")
	Password := request.PostFormValue("password")
	user, err := model.GetUserByUsername(Username)
	if err != nil {
		utils.ExecuteTemplateWithoutNavbar(
			writer, []string{"Password or username is incorrect"}, "templates/sign-in.html",
		)
		return
	}

	if user.Password == Password || CheckPasswordHash(Password, user.Password) {
		DoLogin(writer, *user)
		http.Redirect(writer, request, "/", http.StatusFound)
		return
	}
	http.Redirect(writer, request, "/login-page", http.StatusBadRequest)

}

func DoLogin(writer http.ResponseWriter, user model.User) {
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
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(
		writer, &http.Cookie{
			Name:    "session_token",
			Value:   tokenString,
			Path:    "/",
			Expires: expirationTime,
		},
	)

	CurrentSession := model.Session{
		UserID:    user.ID,
		UUID:      tokenString,
		CreatedAt: time.Now(),
		DeletedAt: expirationTime,
		LastLogin: time.Now(),
		IP:        model.SetInet(),
	}

	var session model.Session

	session = CurrentSession

	server.DB.Create(&session)
}

func Logout(writer http.ResponseWriter, request *http.Request) {

	claims := utils.CheckCookie(writer, request)
	if claims == nil {
		http.Redirect(writer, request, "/index", http.StatusSeeOther)
	}

	var session model.Session
	server.DB.Last(&session)
	session.DeletedAt = time.Now()
	server.DB.Save(&session)

	c := &http.Cookie{
		Name:    "session_token",
		Expires: time.Now(),
	}

	http.SetCookie(writer, c)
	http.Redirect(writer, request, "/index", http.StatusSeeOther)

}

func SignUp(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if request.PostFormValue("password") == request.PostFormValue("repassword") {
		user := model.User{
			Username: request.PostFormValue("username"),
			Email:    request.PostFormValue("email"),
			Password: request.PostFormValue("password"),
		}
		v := validators.UserValidator{User: &user}
		if v.Check(); !v.IsValid() {
			t, _ := template.ParseFiles("templates/sign-up.html")
			t.Execute(writer, v.Errors())
			return
		}
		user.Password, _ = HashPassword(user.Password)
		server.DB.Create(&user)
		http.Redirect(writer, request, "/sign-in", 302)
	} else {
		t, _ := template.ParseFiles("templates/sign-up.html")
		t.Execute(writer, []string{"Two passwords doesn't match"})
		return
	}
}

func LoginPage(writer http.ResponseWriter, request *http.Request) {
	t, _ := template.ParseFiles("templates/sign-in.html")
	t.Execute(writer, nil)
}

func SignUpPage(writer http.ResponseWriter, request *http.Request) {
	t, _ := template.ParseFiles("templates/sign-up.html")
	t.Execute(writer, nil)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
