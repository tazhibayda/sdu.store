package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"sdu.store/server/model"
)

func Index(writer http.ResponseWriter, request *http.Request) {
	tm, _ := template.ParseFiles("templates/index.gohtml")
	t, _ := template.ParseFiles("templates/layouts/header.gohtml")

	if err := t.Execute(writer, nil); err != nil {
		panic(err)
	}
	cookie, err := request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			writer.WriteHeader(http.StatusUnauthorized)
			//return
		}
		writer.WriteHeader(http.StatusBadRequest)
		//return
	}
	if cookie == nil {
		writer.WriteHeader(http.StatusUnauthorized)
		//return
	} else {
		userSession, exists := model.CurrentSession[cookie.Value]
		if !exists {
			// If the session token is not present in session map, return an unauthorized error

		}
		fmt.Println(userSession.UserID)
		// If the session is present, but has expired, we can delete the session, and return
		// an unauthorized status
	}
	err = tm.Execute(writer, nil)
	if err != nil {
		return
	}
}
