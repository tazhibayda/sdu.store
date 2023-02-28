package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/server/validators"
)

func Login(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		fmt.Println("GET")
		t, _ := template.ParseFiles("templates/login.html")
		t.Execute(writer, nil)
	} else {

		Username := request.PostFormValue("username")
		Password := request.PostFormValue("password")

		user, err := validators.GetUserByUsername(Username)
		if err != nil {
			panic("User not Exists")
		}
		if user.Password == Password {
			// Need to add cookie
			//cookie := http.Cookie{Name: "session", Value: string(user)}
			//
			//http.SetCookie(writer, &cookie)
		}
	}
}

func Logout(writer http.ResponseWriter, request *http.Request) {

}

func SignUp(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		fmt.Println(err)
		t, _ := template.ParseFiles("templates/sign-up.html")
		t.Execute(writer, []string{"Server error"})
		return
		//danger method
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
			//danger method
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
