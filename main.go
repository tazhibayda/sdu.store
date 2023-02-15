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
		err := DB.AutoMigrate(model.Session{}, model.User{}, model.Userdata{})
		if err != nil {
			return
		}
	}

	a := http.NewServeMux()

	// Request for postman
	a.HandleFunc("/request/users", model.GetUsers)
	a.HandleFunc("/request/user", model.GetUserByID)
	//

	// Admin Setting
	a.HandleFunc("/Admin", model.AdminServe)
	a.HandleFunc("/Admin/user/create", model.CreateUser)
	a.HandleFunc("/Admin/user", model.AdminUsers)
	a.HandleFunc("/Admin/user/delete/", model.DeleteUser)
	a.HandleFunc("/Admin/session", model.AdminServe)
	a.HandleFunc("/Admin/userdata", model.AdminUserdata)
	a.HandleFunc("/Admin/userdata/create", model.CreateUserdata)
	a.HandleFunc("/Admin/userdata/delete/", model.DeleteUserdata)
	//

	err := http.ListenAndServe(":9090", a)
	if err != nil {
		log.Fatal(err.Error())
	}

}
