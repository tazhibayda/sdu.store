package admin

import (
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/server/validators"
	"sdu.store/utils"
	"strconv"
	"strings"
)

type CategoryTable struct {
	Categories []model.Category
	Search     string
}

func CategoryPage(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(request.URL.Query().Get("id"))
	category, err := model.GetCategoryByID(id)

	if err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}

	utils.ExecuteTemplateWithoutNavbar(
		writer, request, category, "templates/admin/base.html", "templates/admin/navbar.html",
		"templates/admin/AdminCategory.html",
	)
	return
}

func CategoryDelete(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(request.URL.Query().Get("id"))
	category, err := model.GetCategoryByID(id)

	if err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}

	if err := category.Delete(); err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}

	http.Redirect(writer, request, "/Admin/categories", http.StatusSeeOther)
}

func Category(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(request.URL.Query().Get("id"))
	category, err := model.GetCategoryByID(id)

	if err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}

	category.Name = request.FormValue("name")
	validator := validators.CategoryValidator{Category: &category}
	if validator.Check(); !validator.IsValid() {
		utils.ExecuteTemplateWithoutNavbar(
			writer, request, validator.Errors(), "templates/admin/base.html", "templates/admin/navbar.html",
			"templates/admin/AdminAddCategory.html",
		)
		return
	}

	if err := category.Update(); err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}
	http.Redirect(writer, request, "/Admin/categories", http.StatusSeeOther)
	return
}

func Categories(w http.ResponseWriter, r *http.Request) {
	categories, err := model.GetAllCategory()
	if err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}

	hasFilter, categoryTable := HasFilterCategoryTable(r)
	if hasFilter {
		sortCategoryTable(categories, &categoryTable)
	} else {
		categoryTable.Categories = categories
	}
	utils.ExecuteTemplateWithoutNavbar(
		w, r, categoryTable, "templates/admin/base.html", "templates/admin/navbar.html",
		"templates/admin/AdminCategories.html",
	)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	category = model.Category{Name: r.FormValue("name")}

	validator := validators.CategoryValidator{Category: &category}

	if validator.Check(); !validator.IsValid() {
		utils.ExecuteTemplateWithoutNavbar(
			w, r, validator.Errors(), "templates/admin/base.html", "templates/admin/navbar.html",
			"templates/admin/AdminAddCategory.html",
		)
		return
	}

	if err := server.DB.Create(&category).Error; err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}

	http.Redirect(w, r, "/Admin/categories", http.StatusSeeOther)
}

func CreateCategoryPage(writer http.ResponseWriter, request *http.Request) {
	utils.ExecuteTemplateWithoutNavbar(
		writer, request, nil, "templates/admin/base.html", "templates/admin/navbar.html",
		"templates/admin/AdminAddCategory.html",
	)
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
