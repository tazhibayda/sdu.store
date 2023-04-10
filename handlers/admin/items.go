package admin

import (
	"gorm.io/gorm"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/utils"
	"strconv"
)

type ItemTable struct {
	Items  []model.Item
	Search string
}

func Items(writer http.ResponseWriter, request *http.Request) {
	var itemsTable ItemTable

	if err := filterItems(request, &itemsTable); err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}

	utils.ExecuteTemplateWithoutNavbar(
		writer, request, itemsTable, "templates/admin/base.html", "templates/admin/navbar.html",
		"templates/admin/AdminItems.html",
	)
}

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
	productID, err := strconv.Atoi(request.FormValue("product-id"))
	if err != nil {
		utils.BadRequest(writer, request, err)
		return
	}

	var item = model.Item{Barcode: request.FormValue("barcode"), ProductID: uint(productID), Color: request.FormValue("color"), Size: request.FormValue("size")}
	if err := item.Create(); err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}
	http.Redirect(writer, request, "/Admin/products", http.StatusSeeOther)
}

func filterItems(request *http.Request, table *ItemTable) error {
	db := server.DB.Model(&model.Item{})
	if db.Error != nil {
		return db.Error
	}
	var err error
	if db, err = filterByBarcode(db, table, request); err != nil {
		return err
	}
	var items []model.Item
	if err := db.Find(&items).Error; err != nil {
		return err
	}

	table.Items = items

	return nil
}

func filterByBarcode(db *gorm.DB, table *ItemTable, request *http.Request) (*gorm.DB, error) {
	if request.FormValue("search") == "" {
		return db, nil
	}
	table.Search = request.FormValue("search")
	if err := db.Where("barcode iLike ?", "%"+table.Search+"%").Error; err != nil {
		return nil, err
	}
	return db, nil
}
