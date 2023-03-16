package admin

import (
	uuid2 "github.com/google/uuid"
	"html/template"
	"net/http"
	"sdu.store/handlers"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/utils"
	"time"
)

func AdminLoginPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/admin/sign-in.html")
	t.Execute(w, nil)
}

func AdminLogout(writer http.ResponseWriter, request *http.Request) {
	session, err := utils.SessionStaff(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/Admin", http.StatusTemporaryRedirect)
		return
	}
	user, _ := model.GetUserByID(session.UserID)
	user.DeleteSessions()
	http.Redirect(writer, request, "/Admin", http.StatusFound)
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

	if handlers.CheckPasswordHash(Password, user.Password) {
		if !user.Is_staff {
			utils.ExecuteTemplateWithoutNavbar(
				w, []string{"User doesn't have access to admin page"}, "templates/admin/sign-in.html",
			)
			return
		}
		doLogin(w, *user)
		http.Redirect(w, r, "/Admin", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusBadRequest)
}

func AdminServe(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.SessionStaff(w, r); err != nil {
		http.Redirect(w, r, "/Admin/login-page", http.StatusTemporaryRedirect)
		return
	}
	tm, _ := template.ParseFiles(
		"templates/admin/base.html", "templates/admin/index.html", "templates/admin/navbar.html",
	)
	tm.ExecuteTemplate(w, "base", nil)
	return
}

func doLogin(writer http.ResponseWriter, user model.User) {

	uuid := uuid2.NewString()

	sTime := 60 * 60

	http.SetCookie(
		writer, &http.Cookie{
			Name:   "session_token",
			Value:  uuid,
			MaxAge: sTime,
		},
	)
	CurrentSession := model.Session{
		UserID:    user.ID,
		UUID:      uuid,
		CreatedAt: time.Now(),
		DeletedAt: time.Now().Add(time.Duration(sTime)),
		LastLogin: time.Now(),
		IP:        model.SetInet(),
	}
	server.DB.Create(&CurrentSession)
}
