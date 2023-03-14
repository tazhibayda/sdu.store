package handlers

import (
	"html/template"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
)

func Index(writer http.ResponseWriter, request *http.Request) {
	tm, _ := template.ParseFiles("templates/index.gohtml")
	CallHeaderHtml(writer, request)
	tm.Execute(writer, nil)
}

func Account(writer http.ResponseWriter, request *http.Request) {
	CallHeaderHtml(writer, request)
	tm, _ := template.ParseFiles("templates/account.gohtml")

	claims := &Claims{}
	userdata := &model.Userdata{}
	type DataUser struct {
		User     *model.User
		Userdata *model.Userdata
	}
	du := &DataUser{}
	claims = CheckCookie(writer, request)

	if claims != nil {
		server.DB.Where("user_id = ?", claims.User.ID).Find(&userdata)
		du = &DataUser{User: claims.User, Userdata: userdata}
	} else {
		du = nil
	}

	err := tm.Execute(writer, du)
	if err != nil {
		return
	}
}
