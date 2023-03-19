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
	claim := utils.CheckCookie(writer, request)
	if claim != nil {
		tm, _ := template.ParseFiles("templates/base.html", "templates/index.html", "templates/private.navbar.html")
		tm.ExecuteTemplate(writer, "base", *claim.User)
	} else {
		tm, _ := template.ParseFiles("templates/base.html", "templates/index.html", "templates/public.navbar.html")
		tm.ExecuteTemplate(writer, "base", nil)
	}
}

func Account(writer http.ResponseWriter, request *http.Request) {
	claims := &model.Claims{}
	claims = utils.CheckCookie(writer, request)

	if claims == nil {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		writer.Write([]byte("<script>alert('Please login')</script>"))
		return
	}

	//CallHeaderHtml(writer, request)
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

	err := tm.ExecuteTemplate(writer, "base", du)
	if err != nil {
		return
	}
}
