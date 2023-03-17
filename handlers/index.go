package handlers

import (
	"html/template"
	"net/http"
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

func Account(writer http.ResponseWriter, request *http.Request) {
	claims := &Claims{}
	claims = CheckCookie(writer, request)

	if claims == nil {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		writer.Write([]byte("<script>alert('Please login')</script>"))
		return
	}

	CallHeaderHtml(writer, request)
	tm, _ := template.ParseFiles("templates/account.gohtml")

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

	err := tm.Execute(writer, du)
	if err != nil {
		return
	}
}
