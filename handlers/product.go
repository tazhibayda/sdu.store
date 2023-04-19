package handlers

import (
	"net/http"
	"sdu.store/server/model"
	"sdu.store/utils"
	"strconv"
)

func Product(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.FormValue("id"))
	if err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}
	product, err := model.GetProductByID(id)
	if err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}

	utils.ExecuteTemplateWithNavbar(writer, request, product, "templates/product.html")
}
