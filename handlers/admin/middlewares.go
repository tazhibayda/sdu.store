package admin

import (
	"net/http"
	"sdu.store/utils"
)

func StaffLoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, err := utils.SessionStaff(writer, request)
		if err != nil {
			http.Redirect(writer, request, "/Admin/login", http.StatusTemporaryRedirect)
			return
		}
		next(writer, request)
	}
}

func AdminLoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, err := utils.SessionAdmin(writer, request)
		if err != nil {
			http.Redirect(writer, request, "/Admin/login", http.StatusTemporaryRedirect)
			return
		}
		next(writer, request)
	}
}
