package admin

import (
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/utils"
)

type PurchaseRow struct {
	ID int
	model.User
	model.Product
	model.Item
}

func GetAllPurchase(writer http.ResponseWriter, request *http.Request) {
	purchases := model.GetAllPurchases()

	purchasesTable := []PurchaseRow{}

	for _, purchase := range purchases {
		row := PurchaseRow{User: model.User{}, Product: model.Product{}, Item: model.Item{}}
		row.User.ID = uint(purchase.UserID)
		row.Item.ID = uint(purchase.ItemID)
		server.DB.First(&row.User)
		server.DB.First(&row.Item)
		row.Product.ID = row.Item.ProductID
		server.DB.First(&row.Product)
		purchasesTable = append(purchasesTable, row)
	}

	utils.ExecuteTemplateWithoutNavbar(
		writer, request, purchasesTable, "templates/admin/base.html",
		"templates/admin/navbar.html",
		"templates/admin/AdminPurchases.html",
	)
}
