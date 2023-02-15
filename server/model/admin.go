package model

import (
	"html/template"
	"net/http"
	"sdu.store/server"
)

var tmpl *template.Template

func AdminServe(w http.ResponseWriter, r *http.Request) {

	//tm := template.New("templates/Admin/Admin.gohtml")
	http.ServeFile(w, r, "templates/Admin/Admin.gohtml")
	//tm, _ = tm.Parse("{{.}}")
	//tm.Execute(w,t)

}

func AdminUserdata(w http.ResponseWriter, r *http.Request) {

	tm, _ := template.ParseFiles("templates/Admin/AdminUserdata.gohtml")
	userdata := []Userdata{}
	server.DB.Find(&userdata)
	err := tm.Execute(w, userdata)
	if err != nil {
		return
	}

}

func AdminUsers(w http.ResponseWriter, r *http.Request) {
	//files, err := template.ParseFiles("templates/Admin/*.gothml")
	//if err != nil {
	//	return
	//}

	tm, _ := template.ParseFiles("templates/Admin/AdminUser.gohtml")
	user := []User{}
	server.DB.Find(&user)
	err := tm.Execute(w, user)
	if err != nil {
		return
	}
}
