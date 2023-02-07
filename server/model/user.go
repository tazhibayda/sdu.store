package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sdu.store/server"
)

type User struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var Users []User

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User

	if r.Method == "POST" {
		login := r.FormValue("login")
		password := r.FormValue("password")
		username := r.FormValue("username")
		user = User{Login: login, Password: password, Username: username}
	}
	server.DB.Create(&user)
	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := []User{}
	server.DB.Find(&users)
	fmt.Println(users)
	json.NewEncoder(w).Encode(users)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {

	user := User{}
	server.DB.Where("id", r.FormValue("user_id")).Find(&user)
	fmt.Println(user)
	json.NewEncoder(w).Encode(user)
}
