package utils

import (
	"errors"
	"fmt"
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

func ServerErrorHandler(writer http.ResponseWriter, request *http.Request, err error) {
	ErrorLogger(err.Error(), request)
	ErrorTemplate(writer, "Server Error", http.StatusInternalServerError, "templates/error.html")
}

func NotFound(writer http.ResponseWriter, request *http.Request, err error) {
	ErrorLogger(err.Error(), request)
	ErrorTemplate(writer, "Not Found", http.StatusNotFound, "templates/error.html")
}

func BadRequest(writer http.ResponseWriter, request *http.Request, err error) {
	ErrorLogger(err.Error(), request)
	ErrorTemplate(writer, "Bad Request", http.StatusBadRequest, "templates/error.html")
}

func Session(writer http.ResponseWriter, request *http.Request) (user *model.User, err error) {
	cookie, err := CheckCookie(writer, request)
	if err != nil {
		return nil, err
	}
	user = cookie.User
	if err = server.DB.Find(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func SessionStaff(writer http.ResponseWriter, request *http.Request) (user *model.User, err error) {
	cookie, err := CheckCookie(writer, request)
	if err != nil {
		return nil, err
	}
	user = cookie.User
	err = server.DB.Find(user).Error
	if !user.IsStaff() || err != nil {
		err = errors.New("Invalid staff access")
		return nil, err
	}
	return
}

func SessionAdmin(writer http.ResponseWriter, request *http.Request) (user *model.User, err error) {
	cookie, err := CheckCookie(writer, request)
	if err != nil {
		return nil, err
	}
	user = cookie.User
	err = server.DB.Find(user).Error
	if !user.IsAdmin() || err != nil {
		err = errors.New("Invalid admin access")
		return nil, err
	}
	return
}

func CheckCookie(writer http.ResponseWriter, request *http.Request) (*model.Claims, error) {
	claims := &model.Claims{}

	cookie, err := request.Cookie("session_token")
	if err != nil {
		return nil, err
	}
	if cookie == nil {
		return nil, InvalidTokenError
	}

	key := cookie.Value

	token, err := jwt.ParseWithClaims(
		key, claims, func(token *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		writer.WriteHeader(http.StatusUnauthorized)
		return nil, InvalidTokenError
	}

	var exists bool
	if server.DB.Model(&model.User{}).
		Select("count(*) > 0").
		Where("id = ?", claims.User.ID).
		Find(&exists); !exists {
		return nil, InvalidTokenError
	}
	return claims, nil
}

func ErrorTemplate(w http.ResponseWriter, err string, status int, files ...string) {
	t, _ := template.ParseFiles(files...)
	output := struct {
		Error  string
		Status int
	}{
		err,
		status,
	}
	t.Execute(w, output)
}

func ExecuteTemplateWithNavbar(
	w http.ResponseWriter, r *http.Request, data interface{}, user model.User, files ...string,
) {
	t, err := template.ParseFiles(files...)
	if err != nil {
		ServerErrorHandler(w, r, err)
		return
	}
	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		ServerErrorHandler(w, r, err)
	}
}

func ExecuteTemplateWithoutNavbar(w http.ResponseWriter, r *http.Request, data interface{}, files ...string) {
	t, err := template.ParseFiles(files...)
	if err != nil {
		ServerErrorHandler(w, r, err)
		return
	}
	if err = t.ExecuteTemplate(w, "base", data); err != nil {
		ServerErrorHandler(w, r, err)
	}
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

func ErrorLogger(err string, request *http.Request) {
	fmt.Fprintf(os.Stdout, "%s:\t %s: \t", request.URL, request.Method)
	fmt.Fprintf(os.Stdout, "%s \n", err)
}

func PrintInfo(message string) {
	fmt.Fprintln(os.Stdout, message)
}
