package handlers

import (
	"fmt"
	uuid2 "github.com/google/uuid"
	"html/template"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/server/validators"
	"time"
)

func Login(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {

		CallHeaderHtml(writer, request)

		t, _ := template.ParseFiles("templates/login.gohtml")
		if err := t.Execute(writer, nil); err != nil {
			panic(err)
		}

	} else {

		Username := request.PostFormValue("username")
		Password := request.PostFormValue("password")
		user, err := validators.GetUserByUsername(Username)
		if err != nil {
			panic("User not Exists")
		}
		if user.Password == Password {

			doLogin(writer, *user)

			http.Redirect(writer, request, "/index", http.StatusSeeOther)
		}
	}
}

func doLogin(writer http.ResponseWriter, user model.User) {

	uuid := uuid2.NewString()

	sTime := 60 * 5

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

	var session model.Session

	session = CurrentSession

	server.DB.Create(&session)
}

func Logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			writer.WriteHeader(http.StatusUnauthorized)
		} else {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	var session model.Session
	server.DB.Last(&session)
	session.DeletedAt = time.Now()
	server.DB.Save(&session)

	cookie = &http.Cookie{
		Name:   "session_token",
		Value:  "",
		MaxAge: 1,
	}
	http.SetCookie(writer, cookie)
	http.Redirect(writer, request, "/index", http.StatusSeeOther)
}

func SignUp(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		fmt.Println(err)
		t, _ := template.ParseFiles("templates/sign-up.html")
		t.Execute(writer, []string{"Server error"})
		return
	}
	if request.PostFormValue("password") == request.PostFormValue("repassword") {
		user := model.User{
			Username: request.PostFormValue("username"),
			Password: request.PostFormValue("password"),
		}
		v := validators.UserValidator{User: &user}
		if v.Check(); !v.IsValid() {
			t, _ := template.ParseFiles("templates/sign-up.html")
			fmt.Println(v.Errors())
			t.Execute(writer, v.Errors())
			return
		}
		if err := server.DB.Create(&user); err != nil {
			t, _ := template.ParseFiles("templates/sign-up.html")
			t.Execute(writer, []string{"Server error"})
			return
		}
		fmt.Println("Qazx")
		http.Redirect(writer, request, "/sign-in", 302)
	} else {
		t, _ := template.ParseFiles("templates/sign-up.html")
		t.Execute(writer, []string{"Two passwords doesn't match"})
		return
	}
}

func LoginPage(writer http.ResponseWriter, request *http.Request) {

}

func SignUpPage(writer http.ResponseWriter, request *http.Request) {
	t, _ := template.ParseFiles("templates/sign-up.html")
	t.Execute(writer, nil)
}

func CallHeaderHtml(writer http.ResponseWriter, request *http.Request) {

	t, _ := template.ParseFiles("templates/layouts/header.gohtml")
	user := &model.User{}
	cookie, err := request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			writer.WriteHeader(http.StatusUnauthorized)
			//return
		} else {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if cookie != nil {
		var session *model.Session
		server.DB.Where("uuid = ?", cookie.Value).Find(&session)

		server.DB.Where("id = ?", session.UserID).Find(&user)
	}
	if err := t.Execute(writer, user.Email); err != nil {
		panic(err)
	}
}
