package utils

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"
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

func SessionStaff(writer http.ResponseWriter, request *http.Request) (session model.Session, err error) {
	cookie, err := request.Cookie("session_token")
	if err == nil {
		session = model.Session{UUID: cookie.Value}
		if ok := session.CheckStaff(); !ok {
			err = errors.New("Invalid admin session")
		}
	}
	return
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
