package admin

import (
	"net/http"
	"sdu.store/handlers"
	"sdu.store/server/model"
	"sdu.store/utils"
	"time"
)

func AdminLoginPage(w http.ResponseWriter, r *http.Request) {
	_, err := utils.SessionStaff(w, r)
	if err == nil {
		http.Redirect(w, r, "/Admin", http.StatusSeeOther)
		return
	}

	access := r.URL.Query().Get("access")

	if access != "" {
		utils.ExecuteTemplateWithoutNavbar(w, r, []string{"Need " + access + " access"}, "templates/admin/sign-in.html")
		return
	}
	utils.ExecuteTemplateWithoutNavbar(w, r, nil, "templates/admin/sign-in.html")
}

func AdminLogout(writer http.ResponseWriter, request *http.Request) {
	_, err := utils.CheckCookie(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	c := &http.Cookie{
		Name:    "session_token",
		Path:    "/",
		Expires: time.Now(),
	}

	http.SetCookie(writer, c)
	http.Redirect(writer, request, "/Admin/login", http.StatusSeeOther)

}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	Username := r.PostFormValue("username")
	Password := r.PostFormValue("password")
	user, err := model.GetUserByUsername(Username)
	if err != nil {
		utils.ExecuteTemplateWithoutNavbar(
			w, r, []string{"Password or username is incorrect"}, "templates/admin/sign-in.html",
		)
		return
	}

	if utils.CheckPasswordHash(Password, user.Password) || user.Password == Password {
		if !user.Is_staff {
			utils.ExecuteTemplateWithoutNavbar(
				w, r, []string{"User doesn't have access to admin page"}, "templates/admin/sign-in.html",
			)
			return
		}

		if err := handlers.DoLogin(w, *user); err != nil {
			utils.ServerErrorHandler(w, r, err)
			return
		}
		http.Redirect(w, r, "/Admin", http.StatusFound)
		return
	}

	utils.ExecuteTemplateWithoutNavbar(
		w, r, []string{"Password or username is not correct"}, "templates/admin/sign-in.html",
	)
}

func AdminServe(w http.ResponseWriter, r *http.Request) {
	access := r.URL.Query().Get("access")
	if access != "" {
		utils.ExecuteTemplateWithoutNavbar(
			w, r, []string{"Need " + access + " access"}, "templates/admin/base.html", "templates/admin/index.html",
			"templates/admin/navbar.html",
		)
		return
	}
	utils.ExecuteTemplateWithoutNavbar(
		w, r, nil, "templates/admin/base.html", "templates/admin/index.html", "templates/admin/navbar.html",
	)
}
