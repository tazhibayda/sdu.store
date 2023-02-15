package main

import (
	"gorm.io/gorm"
	"log"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
)

var DB *gorm.DB = server.DB

func main() {

	a := http.NewServeMux()

	a.HandleFunc("/Users", model.GetUsers)
	a.HandleFunc("/User", model.GetUserByID)
	a.HandleFunc("/Create", model.CreateUser)
	a.HandleFunc("/Admin", model.AdminServe)
	a.HandleFunc("/Admin/user", model.AdminUsers)
	a.HandleFunc("/Admin/user/delete/", model.DeleteUser)
	a.HandleFunc("/Admin/session", model.AdminServe)
	a.HandleFunc("/Admin/userdata", model.AdminUserdata)
	err := http.ListenAndServe(":9090", a)
	if err != nil {
		log.Fatal(err.Error())
	}

}
