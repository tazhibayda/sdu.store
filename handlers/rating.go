package handlers

import (
	"net/http"
	"sdu.store/server/model"
	"sdu.store/server/validators"
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
	validator := validators.NewRatingValidator(&model.Rating{UserID: int(user.ID), ProductID: productID, Value: value})
	if validator.Check(); !validator.IsValid() {
		utils.BadRequest(writer, request, nil)
		return
	}

	if err = validator.Rating.Create(); err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}
	http.Redirect(writer, request, "", http.StatusSeeOther)
}
