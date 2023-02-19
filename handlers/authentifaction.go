package handlers

import (
	"html/template"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
)

func Login(writer http.ResponseWriter, request *http.Request) {

}

func Logout(writer http.ResponseWriter, request *http.Request) {

}

func SignUp(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		//danger method
	}
	if request.PostFormValue("password") == request.PostFormValue("repassword") {
		user := model.User{
			Username: request.PostFormValue("username"),
			Password: request.PostFormValue("password"),
		}
		if err := server.DB.Create(&user); err != nil {
			//danger method
		}
		http.Redirect(writer, request, "/sign-in", 302)
	} else {
		//danger method
	}
}

func LoginPage(writer http.ResponseWriter, request *http.Request) {

}

func SignUpPage(writer http.ResponseWriter, request *http.Request) {
	t, _ := template.ParseFiles("templates/sign-up.html")
	t.Execute(writer, nil)
}
