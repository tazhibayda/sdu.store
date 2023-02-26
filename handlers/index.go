package handlers

import "net/http"

func Index(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "templates/index.html")
}
