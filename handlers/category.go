package handlers

import (
	"net/http"
	"strconv"

	"sdu.store/server/model"
	"sdu.store/utils"
)

func Category(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		utils.BadRequest(writer, request, err)
		return
	}
	category, err := model.GetCategoryByID(id)

	if err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}

	utils.ExecuteTemplateWithoutNavbar(
		writer, request, category, "templates/base.html", "templates/category.html", "templates/public.navbar.html",
	)
}
