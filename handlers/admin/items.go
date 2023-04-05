package admin

import (
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/utils"
	"strconv"
)

func AddItemPage(writer http.ResponseWriter, request *http.Request) {
	var product model.Product
	if err := server.DB.First(&product, "id = ?", request.FormValue("id")).Error; err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}
	utils.ExecuteTemplateWithoutNavbar(
		writer, request, product, "templates/admin/base.html", "templates/admin/navbar.html",
		"templates/admin/AdminAddItem.html",
	)
}

func AddItem(writer http.ResponseWriter, request *http.Request) {
	quantity, err := strconv.Atoi(request.FormValue("quantity"))
	if err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}

	productID, err := strconv.Atoi(request.FormValue("product-id"))
	if err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}

	var item = model.Item{ProductID: uint(productID), Color: request.FormValue("color"), Size: request.FormValue("size"), Quantity: int64(quantity)}
	if err := item.CreateOrUpdate(); err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}
	http.Redirect(writer, request, "/Admin/products", http.StatusSeeOther)
}
