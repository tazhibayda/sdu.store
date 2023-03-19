package model

import (
	"net/http"
	"sdu.store/server"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"login"`
	Username string `json:"username"`
	Password string `json:"password"`
	Is_admin bool   `json:"is_admin"`
	Is_staff bool   `json:"is_staff"`
}

var Users []User

func GetUserByUsername(username string) (*User, error) {
	user := User{}
	server.DB.Where("username", username).Find(&user)

	if user.Username == "" {
		return nil, http.ErrAbortHandler
	}

	return &user, nil
}

func GetUserByID(id int64) (*User, error) {
	user := User{}
	server.DB.Where("id", id).Find(&user)

	if user.Username == "" {
		return nil, http.ErrAbortHandler
	}

	return &user, nil
}

func (user *User) IsAdmin() bool {
	return user.Is_admin
}

func (user *User) IsStaff() bool {
	return user.Is_staff
}

func (user *User) Delete() {
	server.DB.Where("ID=?", user.ID).Delete(&User{})
}

func (user *User) Update() {
	isStaff := user.Is_staff
	isAdmin := user.Is_admin
	server.DB.First(user)
	user.Is_admin = isAdmin
	user.Is_staff = isStaff || isAdmin
	server.DB.Save(user)
}
