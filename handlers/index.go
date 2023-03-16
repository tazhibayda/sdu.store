package handlers

import (
	"html/template"
	"net/http"
	"sdu.store/server/model"
	"sdu.store/utils"
)

func Index(writer http.ResponseWriter, request *http.Request) {
	session, err := utils.Session(writer, request)
	if err != nil {
		tm, _ := template.ParseFiles("templates/base.html", "templates/index.html", "templates/public.navbar.html")
		tm.ExecuteTemplate(writer, "base", nil)
		return
	}
	user, _ := model.GetUserByID(session.UserID)
	tm, _ := template.ParseFiles("templates/base.html", "templates/index.html", "templates/private.navbar.html")
	tm.ExecuteTemplate(writer, "base", user)
}
