package admin

import (
	"html/template"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/utils"
	"strings"
)

func AdminCategories(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.SessionStaff(w, r); err != nil {
		http.Redirect(w, r, "login", http.StatusUnauthorized)
		return
	}

	tm, _ := template.ParseFiles("templates/Admin/Categories.gohtml")
	var categories []model.Category
	server.DB.Find(&categories)
	err := tm.Execute(w, categories)
	if err != nil {
		return
	}
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.SessionStaff(w, r); err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	var category model.Category
	if r.Method == "POST" {
		name := r.FormValue("name")
		category = model.Category{Name: name}
	} else {
		http.Redirect(w, r, "/Admin/categories", http.StatusMethodNotAllowed)
	}
	server.DB.Create(&category)
	//json.NewEncoder(w).Encode(user)
	http.Redirect(w, r, "/Admin/categories", http.StatusSeeOther)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.SessionStaff(w, r); err != nil {
		http.Redirect(w, r, "login", http.StatusUnauthorized)
		return
	}

	vars := strings.Split(r.URL.Path, "/")
	categoryID := vars[len(vars)-1]
	category := model.Category{}
	server.DB.Where("ID = ?", categoryID).Delete(&category)
	http.Redirect(w, r, "/Admin/categories", http.StatusSeeOther)
}
