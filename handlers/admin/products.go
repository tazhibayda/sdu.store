package admin

import (
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/server/validators"
	"sdu.store/utils"
	"strconv"
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
	Errors   []string
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

func DeleteProduct(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		utils.BadRequest(writer, request, err)
		return
	}

	product, err := model.GetProductByID(id)
	if err != nil {
		utils.NotFound(writer, request, err)
		return
	}

	if err := product.Delete(); err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}

	http.Redirect(writer, request, "/Admin/products", http.StatusSeeOther)
}

func AddImage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}
	product, err := model.GetProductByID(id)
	if err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}
	filename, err := retrieveFileName(r)
	if err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}
	product.Images = append(product.Images, filename)
	if err = product.Update(); err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}
	http.Redirect(w, r, "/Admin/products", http.StatusSeeOther)
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	var product = &model.Product{}
	err := model.ParseProduct(product, r)
	if err != nil {
		utils.BadRequest(w, r, err)
		return
	}

	if err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}

	validator := validators.ProductValidator{Product: product}
	if validator.Check(); !validator.IsValid() {
		utils.ExecuteTemplateWithoutNavbar(
			w, r, validator.Errors(), "templates/admin/base.html", "templates/admin/navbar.html",
			"templates/admin/AdminAddProduct.html",
		)
	}

	if err = product.Create(); err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}
	http.Redirect(w, r, "/Admin/products", http.StatusSeeOther)
}

func ProductPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		utils.BadRequest(w, r, err)
	}

	p, err := model.GetProductByID(id)
	if err != nil {
		utils.NotFound(w, r, err)
		return
	}
	category, err := model.GetCategoryByID(p.CategoryID)
	if err != nil {
		utils.NotFound(w, r, err)
		return
	}
	product := ProductOutput{p, category.Name, nil}

	utils.ExecuteTemplateWithoutNavbar(
		w, r, product, "templates/admin/base.html", "templates/admin/navbar.html", "templates/admin/AdminProduct.html",
	)
}

func Product(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		utils.BadRequest(w, r, err)
		return
	}
	product.ID = uint(id)
	err = model.ParseProduct(&product, r)
	if err != nil {
		utils.BadRequest(w, r, err)
		return
	}
	p, err := model.GetProductByID(int(product.ID))
	if err != nil {
		utils.NotFound(w, r, err)
		return
	}
	product.Images = p.Images
	product.Colors = p.Colors
	product.Sizes = p.Sizes

	category, err := model.GetCategoryByID(product.CategoryID)
	if err != nil {
		utils.NotFound(w, r, err)
		return
	}

	validator := validators.ProductValidator{Product: &product}
	if validator.Check(); !validator.IsValid() {
		utils.ExecuteTemplateWithoutNavbar(
			w, r, ProductOutput{product, category.Name, validator.Errors()}, "templates/admin/base.html",
			"templates/admin/navbar.html",
			"templates/admin/AdminProduct.html",
		)
		return
	}

	err = product.Update()
	if err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}
	http.Redirect(w, r, "/Admin/products", http.StatusSeeOther)

}

func Products(w http.ResponseWriter, r *http.Request) {
	var productTable ProductTable

	if err := retrieveCategories(r, &productTable); err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}

	if err := filterProducts(r, &productTable); err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}

	utils.ExecuteTemplateWithoutNavbar(
		w, r, productTable, "templates/admin/base.html", "templates/admin/navbar.html",
		"templates/admin/AdminProducts.html",
	)

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
	var products []model.Product
	if err := db.Find(&products).Error; err != nil {
		return err
	}

	for _, product := range products {
		category, err := model.GetCategoryByID(product.CategoryID)
		if err != nil {
			return err
		}
		table.Products = append(table.Products, ProductOutput{product, category.Name, nil})
	}
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

func retrieveFileName(request *http.Request) (string, error) {
	err := request.ParseMultipartForm(2000000000) // grab the multipart form
	if err != nil {
		return "", err
	}

	file, header, err := request.FormFile("image")
	if err != nil {
		return "", err
	}

	defer file.Close()
	if err != nil {
		return "", err
	}

	filename := "images/" + header.Filename

	out, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)

	defer out.Close()
	if err != nil {
		return "", err
	}

	_, err = io.Copy(out, file)

	if err != nil {
		return "", err
	}

	return "/" + filename, nil
}
