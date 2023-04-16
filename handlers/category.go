package handlers

import (
	"net/http"
	"strconv"

	"sdu.store/server/model"
	"sdu.store/utils"
)


func Category(writer http.ResponseWriter, request *http.Request){
	id, _ := strconv.Atoi(request.URL.Query().Get("id"))
	category, err := model.GetCategoryByID(id)

	if err != nil{
		utils.ServerErrorHandler(writer, request, err)
		return
	}

	utils.ExecuteTemplateWithoutNavbar(
		writer, request, category, "templates/base.html", "templates/category.html", "templates/public.navbar.html",
	)
}