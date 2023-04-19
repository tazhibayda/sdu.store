package handlers

import (
	"html/template"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/utils"
	//"sdu.store/utils"
)

func Index(writer http.ResponseWriter, request *http.Request) {
	utils.ExecuteTemplateWithNavbar(
		writer, request, nil, "templates/index.html",
	)
	return
}

func Account(writer http.ResponseWriter, request *http.Request) {
	claims, err := utils.CheckCookie(writer, request)

	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		writer.Write([]byte("<script>alert('Please login')</script>"))
		return
	}

	tm, _ := template.ParseFiles("templates/base.html", "templates/account.html", "templates/private.navbar.html")

	userdata := &model.Userdata{}
	type DataUser struct {
		User     *model.User
		Userdata *model.Userdata
	}
	du := &DataUser{}

	if claims != nil {
		server.DB.Where("user_id = ?", claims.User.ID).Find(&userdata)
		du = &DataUser{User: claims.User, Userdata: userdata}
	} else {
		du = nil
	}

	err = tm.ExecuteTemplate(writer, "base", du)
	if err != nil {
		return
	}
}
