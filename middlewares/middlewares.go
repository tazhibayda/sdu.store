package middlewares

import (
	"fmt"
	"net/http"
	"sdu.store/utils"
)

func StaffLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			_, err := utils.SessionStaff(writer, request)
			if err != nil {
				utils.ErrorLogger(err.Error(), request)
				http.Redirect(writer, request, "/Admin/login?access=staff", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(writer, request)
		},
	)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			_, err := utils.Session(writer, request)
			if err != nil {
				utils.ErrorLogger(err.Error(), request)
				writer.WriteHeader(http.StatusNetworkAuthenticationRequired)
				http.Redirect(writer, request, "/login", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(writer, request)
		},
	)
}

func AdminLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			_, err := utils.SessionAdmin(writer, request)
			if err != nil {
				utils.ErrorLogger(err.Error(), request)
				http.Redirect(writer, request, "/Admin?access=admin", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(writer, request)
		},
	)
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Method")
			origin := r.Header.Get("Origin")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Credentials")
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Connection", "close")
					utils.ServerErrorHandler(w, r, fmt.Errorf("%s", err))
				}
			}()

			next.ServeHTTP(w, r)
		},
	)
}

func LoggingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			utils.PrintInfo(fmt.Sprintf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL))
			next.ServeHTTP(w, r)
		},
	)
}
