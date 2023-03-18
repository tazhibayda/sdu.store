package admin

import (
	"html/template"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/utils"
)

func GetAllSessions(w http.ResponseWriter, r *http.Request) {
	user, err := utils.SessionStaff(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
	}
	CheckAdmin(user, w, r)

	tm, _ := template.ParseFiles("templates/Admin/AdminSession.gohtml")
	var sessions []model.Session
	server.DB.Find(&sessions)
	if err := tm.Execute(w, sessions); err != nil {
		panic(err)
	}
}
