package model

import (
	"time"
)

type Userdata struct {
	UserId      int64     `json:"user_id"`
	Firstname   string    `json:"first_name"`
	Lastname    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	CountryCode string    `json:"country_code"`
	ZIPCode     string    `json:"zip_code"`
	Birthday    time.Time `json:"birthday"`
}

//func SignUp(w http.ResponseWriter, r *http.Request) {
//
//	login := r.FormValue("login")
//	password := r.FormValue("password")
//	username := r.FormValue("username")
//
//	firstname := r.FormValue("firstname")
//	lastname := r.FormValue("lastname")
//	phone := r.FormValue("phone")
//	countrycode := r.FormValue("country_code")
//	zip := r.FormValue("zip")
//	birthday, _ := time.Parse("2006-01-02", r.FormValue("birthday"))
//
//	user := User{Email: login, Password: password, Username: username}
//
//	server.DB.Create(&user)
//
//	u := User{}
//
//	server.DB.Last(&u)
//
//	userdata := Userdata{
//		UserId:      u.ID,
//		Firstname:   firstname,
//		Lastname:    lastname,
//		PhoneNumber: phone,
//		CountryCode: countrycode,
//		ZIPCode:     zip,
//		Birthday:    birthday,
//	}
//	server.DB.Create(&userdata)
//}
