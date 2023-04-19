package handlers

import (
	"net/http"
	"sdu.store/server/model"
	"sdu.store/utils"
	"strconv"
)

func AddComment(writer http.ResponseWriter, request *http.Request) {
	user, _ := utils.Session(writer, request)

	productID, err := strconv.Atoi(request.FormValue("productID"))
	if err != nil {
		utils.BadRequest(writer, request, err)
		return
	}

	text := request.FormValue("text")
	if err != nil {
		utils.BadRequest(writer, request, err)
		return
	}
	comment := model.Comment{UserID: int(user.ID), ProductID: productID, Text: text}
	if err = comment.Create(); err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}
	http.Redirect(writer, request, "/product?id="+strconv.Itoa(productID), http.StatusSeeOther)
}
