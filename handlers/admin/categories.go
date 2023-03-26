package admin

import (
	"html/template"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/utils"
	"strconv"
	"strings"
)

type CategoryTable struct {
	Categories []model.Category
	Search     string
}

func Category(writer http.ResponseWriter, request *http.Request) {
	_, err := utils.SessionStaff(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/Admin/login-page", http.StatusTemporaryRedirect)
		return
	}

	id, _ := strconv.Atoi(request.URL.Query().Get("id"))
	category := model.GetCategoryByID(id)
	if request.Method == "GET" {
		tm, _ := template.ParseFiles(
			"templates/admin/base.html", "templates/admin/navbar.html", "templates/admin/AdminCategory.html",
		)
		err = tm.ExecuteTemplate(writer, "base", category)
		return
	}
	if request.Method == "POST" {
		if "1" == request.URL.Query().Get("delete") {
			category.Delete()
		} else {
			name := request.FormValue("name")
			if name == "" {
				tm, _ := template.ParseFiles(
					"templates/admin/base.html", "templates/admin/navbar.html", "templates/admin/AdminCategory.html",
				)
				err = tm.ExecuteTemplate(writer, "base", category)
				return
			}
			category.Name = name
			category.Update()
		}
		http.Redirect(writer, request, "/Admin/categories", http.StatusTemporaryRedirect)
		return
	}
}

func Categories(w http.ResponseWriter, r *http.Request) {
	_, err := utils.SessionStaff(w, r)
	if err != nil {
		http.Redirect(w, r, "/Admin/login-page", http.StatusTemporaryRedirect)
		return
	}

	var categories []model.Category
	server.DB.Find(&categories)

	hasFilter, categoryTable := HasFilterCategoryTable(r)
	if hasFilter {
		sortCategoryTable(categories, &categoryTable)
	} else {
		categoryTable.Categories = categories
	}
	tm, err := template.ParseFiles(
		"templates/admin/base.html", "templates/admin/navbar.html", "templates/admin/AdminCategories.html",
	)
	err = tm.ExecuteTemplate(w, "base", categoryTable)
	if err != nil {
		return
	}
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	_, err := utils.SessionStaff(w, r)
	if err != nil {
		http.Redirect(w, r, "/Admin/login-page", http.StatusTemporaryRedirect)
		return
	}
	var category model.Category
	if r.Method == "POST" {
		name := r.FormValue("name")
		if name == "" {
			tm, _ := template.ParseFiles(
				"templates/admin/base.html", "templates/admin/navbar.html", "templates/admin/AdminAddCategory.html",
			)
			tm.ExecuteTemplate(w, "base", nil)
			return
		}
		category = model.Category{Name: name}
	} else {
		tm, _ := template.ParseFiles(
			"templates/admin/base.html", "templates/admin/navbar.html", "templates/admin/AdminAddCategory.html",
		)
		tm.ExecuteTemplate(w, "base", nil)
		return
	}
	server.DB.Create(&category)
	http.Redirect(w, r, "/Admin/categories", http.StatusFound)
}

func HasFilterCategoryTable(request *http.Request) (hasFilter bool, filter CategoryTable) {
	if search := request.FormValue("search"); search != "" {
		hasFilter = true
		filter.Search = search
	}
	return
}

func sortCategoryTable(categories []model.Category, table *CategoryTable) {
	ans := []model.Category{}
	for _, category := range categories {
		if isValidCategory(category, table) {
			ans = append(ans, category)
		}
	}
	table.Categories = ans
}

func isValidCategory(category model.Category, table *CategoryTable) bool {
	if table.Search != "" {
		contains := strings.Contains(category.Name, table.Search)
		return contains
	}
	return true
}
