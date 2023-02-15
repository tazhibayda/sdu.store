package main

import (
	"flag"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
)

var DB *gorm.DB = server.DB

func main() {
	restart := flag.Bool("dbRestart", false, "Restarting database")
	flag.Parse()
	if *restart {
		DB.AutoMigrate(model.Session{}, model.User{}, model.Userdata{})
	}

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
