package handlers

import (
	"html/template"
	"net/http"
)

func Index(writer http.ResponseWriter, request *http.Request) {
	tm, _ := template.ParseFiles("templates/index.gohtml")
	CallHeaderHtml(writer, request)
	tm.Execute(writer, nil)
}
