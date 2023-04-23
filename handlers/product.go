package handlers

import (
	"gorm.io/gorm"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/utils"
	"strconv"
)

type ProductTable struct {
	Products   []model.Product
	Search     string
	MaxPrice   int
	MinPrice   int
	MinRating  int
	Categories []struct {
		Category   model.Category
		IsSelected bool
	}
}

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

func Products(writer http.ResponseWriter, request *http.Request) {
	productTable := &ProductTable{}
	if err := retrieveCategories(request, productTable); err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}
	filterProducts(request, productTable)
	utils.ExecuteTemplateWithNavbar(writer, request, productTable, "templates/products.html")
}

func retrieveCategories(request *http.Request, table *ProductTable) error {
	categories, err := model.GetAllCategory()
	if err != nil {
		return err
	}
	hasInput := false
	for _, category := range categories {

		isSelected := false
		if request.FormValue("category"+strconv.Itoa(category.ID)) != "" {
			isSelected = true
		}
		table.Categories = append(
			table.Categories, struct {
				Category   model.Category
				IsSelected bool
			}{Category: category, IsSelected: isSelected},
		)
		hasInput = hasInput || isSelected
	}
	if hasInput {
		return nil
	}
	setDefaultCategory(table)
	return nil
}

func setDefaultCategory(table *ProductTable) {
	for i, _ := range table.Categories {
		table.Categories[i].IsSelected = true
	}
}

func filterProducts(request *http.Request, table *ProductTable) error {
	db := server.DB.Model(&model.Product{})
	if db.Error != nil {
		return db.Error
	}
	var err error
	if db, err = filterProductsByPrice(request, table, db); err != nil {
		return err
	}
	if db, err = filterByName(db, table, request); err != nil {
		return err
	}
	if db, err = filterProductsByCategory(db, table); err != nil {
		return err
	}
	if db, err = filterProductsByRating(request, table, db); err != nil {
		return err
	}
	var products []model.Product
	if err := db.Find(&products).Error; err != nil {
		return err
	}
	table.Products = products
	return nil
}

func filterByName(db *gorm.DB, table *ProductTable, request *http.Request) (*gorm.DB, error) {
	if request.FormValue("search") == "" {
		return db, nil
	}
	table.Search = request.FormValue("search")
	if err := db.Where("name iLike ?", "%"+table.Search+"%").Error; err != nil {
		return nil, err
	}
	return db, nil
}

func filterProductsByCategory(db *gorm.DB, table *ProductTable) (*gorm.DB, error) {
	categoriesID := []int{}
	for _, category := range table.Categories {
		if category.IsSelected {
			categoriesID = append(categoriesID, category.Category.ID)
		}
	}
	if err := db.Where("category_id in ?", categoriesID).Error; err != nil {
		return nil, err
	}
	return db, nil
}

func filterProductsByPrice(request *http.Request, table *ProductTable, db *gorm.DB) (*gorm.DB, error) {
	var err error
	table.MaxPrice, err = strconv.Atoi(request.FormValue("max-price"))
	if err != nil {
		table.MaxPrice = 10000000
	}
	table.MinPrice, err = strconv.Atoi(request.FormValue("min-price"))

	if err != nil {
		table.MinPrice = 0
	}

	db = db.Where("price <= ?", table.MaxPrice)
	db = db.Where("price >= ?", table.MinPrice)
	return db, nil
}

func filterProductsByRating(request *http.Request, table *ProductTable, db *gorm.DB) (*gorm.DB, error) {
	var err error
	table.MinRating, err = strconv.Atoi(request.FormValue("rating"))
	if err != nil {
		table.MinRating = 0
	}
	db = db.Where("rating >= ?", table.MinRating)
	return db, nil
}
