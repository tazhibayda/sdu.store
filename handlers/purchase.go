package handlers

import (
	json2 "encoding/json"
	"net/http"
	"sdu.store/server/model"
	"sdu.store/server/validators"
	"sdu.store/utils"
	"strconv"
)

func Purchase(writer http.ResponseWriter, request *http.Request) {
	user, _ := utils.Session(writer, request)

	size := request.FormValue("size")

	color := request.FormValue("color")

	productID, err := strconv.Atoi(request.FormValue("product_id"))

	if err != nil {
		utils.BadRequest(writer, request, err)
		return
	}

	validator := validators.NewPurchaseValidator(color, size, productID)
	if validator.Check(); !validator.IsValid() {
		json, _ := json2.Marshal(validator.Errors())
		writer.Write(json)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	model.MakePurchase(productID, color, size, int(user.ID))
	json, _ := json2.Marshal(
		[]string{"Successfully bought!!!"})
	writer.Write(json)
	writer.WriteHeader(http.StatusOK)
}
