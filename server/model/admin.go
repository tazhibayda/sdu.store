package model

import (
	"html/template"
	"net/http"
	"sdu.store/server"
)

var _ *template.Template

func AdminServe(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/Admin/Admin.gohtml")
}

func AdminUserdata(w http.ResponseWriter, r *http.Request) {
	tm, _ := template.ParseFiles("templates/Admin/AdminUserdata.gohtml")
	var userdata []Userdata
	server.DB.Find(&userdata)
	err := tm.Execute(w, userdata)
	if err != nil {
		return
	}
}

func AdminUsers(w http.ResponseWriter, r *http.Request) {
	tm, _ := template.ParseFiles("templates/Admin/AdminUser.gohtml")
	var user []User
	server.DB.Find(&user)
	err := tm.Execute(w, user)
	if err != nil {
		return
	}
}

func AdminCategories(w http.ResponseWriter, r *http.Request) {
	tm, _ := template.ParseFiles("templates/Admin/Categories.gohtml")
	var categories []Category
	server.DB.Find(&categories)
	err := tm.Execute(w, categories)
	if err != nil {
		return
	}
}
