package handlers

import (
	"net/http"
	"sdu.store/utils"
)

func NotFoundHandler(writer http.ResponseWriter, request *http.Request) {
	err := "Not Found"
	utils.ErrorLogger(err, request)
	utils.ErrorTemplate(writer, err, http.StatusNotFound, "templates/error.html")
}

func NotAllowedMethod(writer http.ResponseWriter, request *http.Request) {
	err := "Not Allowed Method"
	utils.ErrorLogger(err, request)
	utils.ErrorTemplate(writer, err, http.StatusMethodNotAllowed, "templates/error.html")
}
