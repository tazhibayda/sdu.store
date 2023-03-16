package handlers

import (
	uuid2 "github.com/google/uuid"
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
	Username := request.PostFormValue("username")
	Password := request.PostFormValue("password")
	user, err := model.GetUserByUsername(Username)
	if err != nil {
		utils.ExecuteTemplateWithoutNavbar(
			writer, []string{"Password or username is incorrect"}, "templates/sign-in.html",
		)
		return
	}
	if CheckPasswordHash(Password, user.Password) {
		doLogin(writer, *user)
		http.Redirect(writer, request, "/", http.StatusFound)
		return
	}
	http.Redirect(writer, request, "/login-page", http.StatusBadRequest)
}

func doLogin(writer http.ResponseWriter, user model.User) {

	uuid := uuid2.NewString()

	sTime := 60 * 60

	http.SetCookie(
		writer, &http.Cookie{
			Name:   "session_token",
			Value:  uuid,
			MaxAge: sTime,
		},
	)
	CurrentSession := model.Session{
		UserID:    user.ID,
		UUID:      uuid,
		CreatedAt: time.Now(),
		DeletedAt: time.Now().Add(time.Duration(sTime)),
		LastLogin: time.Now(),
		IP:        model.SetInet(),
	}
	server.DB.Create(&CurrentSession)
}

func Logout(writer http.ResponseWriter, request *http.Request) {
	session, err := utils.Session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}
	user, _ := model.GetUserByID(session.UserID)
	user.DeleteSessions()
	http.Redirect(writer, request, "/", http.StatusFound)
}

func SignUp(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if request.PostFormValue("password") != request.PostFormValue("repassword") {
		t, _ := template.ParseFiles("templates/sign-up.html")
		t.Execute(writer, []string{"Two passwords doesn't match"})
		return
	}
	if err != nil {
		t, _ := template.ParseFiles("templates/sign-up.html")
		t.Execute(writer, []string{"Server error"})
		return
	}
	if request.Method == "POST" {
		if err != nil {
			panic(err)
		}
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
		http.Redirect(writer, request, "/login-page", 302)
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
