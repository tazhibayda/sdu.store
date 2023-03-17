package admin

import (
	"fmt"
	"html/template"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/utils"
	"strconv"
	"strings"
	"time"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.SessionStaff(w, r); err != nil {
		http.Redirect(w, r, "login", http.StatusUnauthorized)
		return
	}
	var product model.Product
	if r.Method == "POST" {
		categoryID, _ := strconv.Atoi(r.FormValue("category"))
		fmt.Println(int64(categoryID))
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		product = model.Product{
			Name:       r.FormValue("Name"),
			CategoryID: categoryID,
			Price:      price,
			CreatedAt:  time.Now(),
		}
		fmt.Println(product.CategoryID)
		server.DB.Create(&product)
	} else {
		http.Redirect(w, r, "/Admin/products", http.StatusMethodNotAllowed)
	}

	http.Redirect(w, r, "/Admin/products", http.StatusSeeOther)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.SessionStaff(w, r); err != nil {
		http.Redirect(w, r, "login", http.StatusUnauthorized)
		return
	}

	if _, err := utils.SessionStaff(w, r); err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
	}
	vars := strings.Split(r.URL.Path, "/")
	productID := vars[len(vars)-1]
	product := model.Product{}
	server.DB.Where("ID = ?", productID).Delete(&product)
	http.Redirect(w, r, "/Admin/products", http.StatusSeeOther)
}

func AdminProducts(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.SessionStaff(w, r); err != nil {
		http.Redirect(w, r, "login", http.StatusUnauthorized)
		return
	}

	tm, _ := template.ParseFiles("templates/Admin/Products.gohtml")
	var products []model.ProductOutput
	var productsDB []model.Product
	sort := "DESCENDING"
	query := r.URL.Query()
	filters, presents := query["sort"]
	if !presents || len(filters) == 0 {
		server.DB.Find(&productsDB)
	} else {
		if filters[0] == "ASCENDING" {
			server.DB.Order("price asc").Find(&productsDB)
			sort = "DESCENDING"
		} else if filters[0] == "DESCENDING" {
			server.DB.Order("price desc").Find(&productsDB)
			sort = "ASCENDING"
		} else {
			server.DB.Find(&productsDB)
		}
	}
	for _, product := range productsDB {
		curr := model.ProductOutput{ID: product.ID, Name: product.Name, Price: product.Price, CreatedAt: product.CreatedAt.Format("2006-02-01")}
		var category model.Category
		server.DB.Where("ID=?", product.CategoryID).Find(&category)
		products = append(products, curr)
	}
	output := struct {
		Categories []model.Category
		Products   []model.ProductOutput
		sort       string
	}{
		Products: products,
		sort:     sort,
	}
	server.DB.Find(&output.Categories)
	err := tm.Execute(w, output)
	if err != nil {
		return
	}
}
