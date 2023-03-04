package model

import (
	"net/http"
	"sdu.store/server"
	"strings"
)

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	if r.Method == "POST" {
		name := r.FormValue("name")
		category = Category{Name: name}
	} else {
		http.Redirect(w, r, "/Admin/categories", http.StatusMethodNotAllowed)
	}
	server.DB.Create(&category)
	//json.NewEncoder(w).Encode(user)
	http.Redirect(w, r, "/Admin/categories", http.StatusSeeOther)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := strings.Split(r.URL.Path, "/")
	categoryID := vars[len(vars)-1]
	category := Category{}
	server.DB.Where("ID = ?", categoryID).Delete(&category)
	http.Redirect(w, r, "/Admin/categories", http.StatusSeeOther)
}

func ConfigCategories() {
	server.DB.Create(&Category{Name: "Hoodies"})
	server.DB.Create(&Category{Name: "Caps"})
	server.DB.Create(&Category{Name: "T-shirts"})
}
