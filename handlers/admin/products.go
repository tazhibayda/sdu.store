package admin

import (
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/utils"
	"strconv"
	"strings"
)

type ProductTable struct {
	Products   []ProductOutput
	Search     string
	MaxPrice   int
	MinPrice   int
	Categories []struct {
		Category   model.Category
		IsSelected bool
	}
}

type ProductOutput struct {
	model.Product
	Category string
}

func AddProductPage(w http.ResponseWriter, r *http.Request) {
	categories, err := model.GetAllCategory()
	if err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}
	utils.ExecuteTemplateWithoutNavbar(
		w, r, categories, "templates/admin/base.html", "templates/admin/navbar.html",
		"templates/admin/AdminAddProduct.html",
	)
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	var product = &model.Product{}
	err := model.ParseProduct(product, r)
	if err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}
	if err = product.Create(); err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}
	http.Redirect(w, r, "/Admin/products", http.StatusSeeOther)
}

func Products(w http.ResponseWriter, r *http.Request) {
	var products []model.Product

	if err := server.DB.Find(&products).Error; err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}

	var productTable ProductTable

	for _, product := range products {
		category, err := model.GetCategoryByID(product.CategoryID)
		if err != nil {
			utils.ServerErrorHandler(w, r, err)
			return
		}
		productTable.Products = append(productTable.Products, ProductOutput{product, category.Name})
	}

	utils.ExecuteTemplateWithoutNavbar(
		w, r, productTable, "templates/admin/base.html", "templates/admin/navbar.html",
		"templates/admin/AdminProducts.html",
	)

}

//func filterByCategories(request *http.Request, table *ProductTable) error {
//
//}

func retrieveCategories(request *http.Request, table *ProductTable) error {
	categoriesID := strings.Split(request.FormValue("categories"), " ")
	categoriesIDMap := make(map[string]bool)
	for _, id := range categoriesID {
		categoriesIDMap[id] = true
	}
	categories, err := model.GetAllCategory()
	if err != nil {
		return err
	}
	for _, category := range categories {
		isSelected := false
		if categoriesIDMap[strconv.Itoa(category.ID)] {
			isSelected = true
		}
		table.Categories = append(
			table.Categories, struct {
				Category   model.Category
				IsSelected bool
			}{Category: category, IsSelected: isSelected},
		)
	}
	return nil

}
