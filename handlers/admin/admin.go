package admin

import (
	"html/template"
	"net/http"
	"sdu.store/handlers"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/utils"
	"time"
)

func AdminLoginPage(w http.ResponseWriter, r *http.Request) {
	user, _ := utils.SessionStaff(w, r)
	CheckAdmin(user, w, r)
	t, _ := template.ParseFiles("templates/admin/sign-in.html")
	t.Execute(w, nil)
}

func AdminLogout(writer http.ResponseWriter, request *http.Request) {

	claims := utils.CheckCookie(writer, request)
	if claims == nil {
		http.Redirect(writer, request, "/Admin/login-page", 200)
		return
	}

	var session model.Session
	server.DB.Last(&session)
	session.DeletedAt = time.Now()
	server.DB.Save(&session)

	c := &http.Cookie{
		Name:    "session_token",
		Path:    "/",
		Expires: time.Now(),
	}

	http.SetCookie(writer, c)
	http.Redirect(writer, request, "/Admin/login-page", http.StatusSeeOther)

}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	Username := r.PostFormValue("username")
	Password := r.PostFormValue("password")
	user, err := model.GetUserByUsername(Username)
	if err != nil {
		utils.ExecuteTemplateWithoutNavbar(
			w, []string{"Password or username is incorrect"}, "templates/admin/sign-in.html",
		)
		return
	}

	if handlers.CheckPasswordHash(Password, user.Password) || user.Password == Password {
		if !user.Is_staff {
			utils.ExecuteTemplateWithoutNavbar(
				w, []string{"User doesn't have access to admin page"}, "templates/admin/sign-in.html",
			)
			return
		}

		handlers.DoLogin(w, *user)
		http.Redirect(w, r, "/Admin", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusBadRequest)
}

func AdminServe(w http.ResponseWriter, r *http.Request) {
	user, err := utils.SessionStaff(w, r)

	tm, _ := template.ParseFiles(
		"templates/admin/base.html", "templates/admin/index.html", "templates/admin/navbar.html",
	)
	if err != nil {
		http.Redirect(w, r, "/Admin/login-page", 302)
	}
	CheckAdmin(user, w, r)
	tm.ExecuteTemplate(w, "base", nil)

}

func CheckAdmin(u *model.User, w http.ResponseWriter, r *http.Request) {

	if u == nil {
		http.Redirect(w, r, "/Admin/login-page", 302)
	}
	if !u.Is_admin || !u.Is_staff {
		http.Redirect(w, r, "/", 302)
	}

}
