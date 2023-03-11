package model

import (
	"html/template"
	"net/http"
	"sdu.store/server"
)

var _ *template.Template

func AdminServe(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/Admin/Admin.gohtml")
}

func AdminUserdata(w http.ResponseWriter, r *http.Request) {
	tm, _ := template.ParseFiles("templates/Admin/AdminUserdata.gohtml")
	var userdata []Userdata
	server.DB.Find(&userdata)
	err := tm.Execute(w, userdata)
	if err != nil {
		return
	}
}

func AdminUsers(w http.ResponseWriter, r *http.Request) {
	var user []User
	query := r.URL.Query()
	filters, presents := query["sort"]
	if !presents || len(filters) == 0 {
		server.DB.Find(&user)
	} else {
		if filters[0] == "" {

		}
	}
	tm, _ := template.ParseFiles("templates/Admin/AdminUser.gohtml")
	err := tm.Execute(w, user)
	if err != nil {
		return
	}
}

func AdminCategories(w http.ResponseWriter, r *http.Request) {
	tm, _ := template.ParseFiles("templates/Admin/Categories.gohtml")
	var categories []Category
	server.DB.Find(&categories)
	err := tm.Execute(w, categories)
	if err != nil {
		return
	}
}

func AdminProducts(w http.ResponseWriter, r *http.Request) {
	tm, _ := template.ParseFiles("templates/Admin/Products.gohtml")
	var products []ProductOutput
	var productsDB []Product

	query := r.URL.Query()
	filters, presents := query["sort"]
	if !presents || len(filters) == 0 {
		server.DB.Find(&productsDB)
	} else {
		if filters[0] == "ASCENDING" {
			server.DB.Order("price asc").Find(&productsDB)
		} else if filters[0] == "DESCENDING" {
			server.DB.Order("price desc").Find(&productsDB)
		} else {
			server.DB.Find(&productsDB)
		}
	}
	for _, product := range productsDB {
		curr := ProductOutput{ID: product.ID, Name: product.Name, Price: product.Price, CreatedAt: product.CreatedAt.Format("2006-02-01")}
		var category Category
		server.DB.Where("ID=?", product.CategoryID).Find(&category)
		products = append(products, curr)
	}
	output := struct {
		Categories []Category
		Products   []ProductOutput
	}{
		Products: products,
	}
	server.DB.Find(&output.Categories)
	err := tm.Execute(w, output)
	if err != nil {
		return
	}
}
