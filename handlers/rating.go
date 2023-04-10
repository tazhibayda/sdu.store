package handlers

import (
	"net/http"
	"sdu.store/server/model"
	"sdu.store/utils"
	"strconv"
)

func AddRating(writer http.ResponseWriter, request *http.Request) {
	user, _ := utils.Session(writer, request)

	productID, err := strconv.Atoi(request.FormValue("productID"))
	if err != nil {
		utils.BadRequest(writer, request, err)
		return
	}

	value, err := strconv.Atoi(request.FormValue("value"))
	if err != nil {
		utils.BadRequest(writer, request, err)
		return
	}
	rating := model.Rating{UserID: int(user.ID), ProductID: productID, Value: value}
	if err = rating.Create(); err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}
	http.Redirect(writer, request, "", http.StatusSeeOther)
}
