package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"html/template"
	"io"
	"net/http"
	"os"
	"sdu.store/server"
	"sdu.store/server/model"
)

const (
	ProductPhotoLocation  = "private/product"
	CategoryPhotoLocation = "private/category"
)

type OutputData struct {
	User model.User
	Data interface{}
}

func Session(writer http.ResponseWriter, request *http.Request) (session model.Session, err error) {
	cookie, err := request.Cookie("session_token")
	if err == nil {
		session = model.Session{UUID: cookie.Value}
		if ok := session.Check(); !ok {
			err = errors.New("Invalid session")
		}
	}
	return
}

func SessionStaff(writer http.ResponseWriter, request *http.Request) (session *model.User, err error) {
	cookie := CheckCookie(writer, request)
	if cookie == nil {
		return nil, http.ErrNoCookie
	}
	user := cookie.User
	server.DB.Find(user)
	if !user.IsStaff() {
		err = errors.New("Invalid staff session")
	}

	return user, err
}

func SessionAdmin(writer http.ResponseWriter, request *http.Request) (user *model.User, err error) {
	cookie := CheckCookie(writer, request)
	if cookie == nil {
		return nil, http.ErrNoCookie
	}
	user = cookie.User
	server.DB.Find(user)
	if !user.IsAdmin() {
		err = errors.New("Invalid admin session")
	}
	return
}

func CallHeader(writer http.ResponseWriter, request *http.Request) {
	user := CheckCookie(writer, request)
	if user != nil {
		tm, _ := template.ParseFiles("templates/base.html", "templates/private.navbar.html")
		tm.ExecuteTemplate(writer, "base", user)
	} else {
		tm, _ := template.ParseFiles("templates/base.html", "templates/public.navbar.html")
		tm.ExecuteTemplate(writer, "base", nil)
	}
}

func CheckCookie(writer http.ResponseWriter, request *http.Request) *model.Claims {
	claims := &model.Claims{}

	cookie, err := request.Cookie("session_token")
	if err != nil {
		return nil
	}
	if cookie != nil {

		key := cookie.Value

		token, err := jwt.ParseWithClaims(
			key, claims, func(token *jwt.Token) (interface{}, error) {
				return model.JwtKey, nil
			},
		)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				writer.WriteHeader(http.StatusUnauthorized)
				return nil
			}
			writer.WriteHeader(http.StatusBadRequest)
			return nil
		}
		if !token.Valid {
			writer.WriteHeader(http.StatusUnauthorized)
			return nil
		}
	}

	return claims
}

func ExecuteTemplateWithNavbar(w http.ResponseWriter, data OutputData, files ...string) {
	t, _ := template.ParseFiles(files...)
	t.Execute(w, data)
}

func ExecuteTemplateWithoutNavbar(w http.ResponseWriter, data interface{}, files ...string) {
	t, _ := template.ParseFiles(files...)
	t.Execute(w, data)
}

func pasteFile(request *http.Request, location string) (filename string, err error) {
	in, header, err := request.FormFile("photo")
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(location + "/photo/" + header.Filename)

	if err != nil {
		return
	}
	defer out.Close()
	io.Copy(out, in)
	return out.Name(), nil
}
