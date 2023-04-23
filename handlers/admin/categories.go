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

type CategoryOutput struct {
	*validators.CategoryValidator
	ErrorsString []string
}

func CategoryPage(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		utils.BadRequest(writer, request, err)
		return
	}
	category, err := model.GetCategoryByID(id)

	if err != nil {
		utils.NotFound(writer, request, err)
		return
	}

	utils.ExecuteTemplateWithoutNavbar(
		writer, request, newCategoryOutput(&category), "templates/admin/base.html",
		"templates/admin/navbar.html",
		"templates/admin/AdminCategory.html",
	)
	return
}

func CategoryDelete(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		utils.BadRequest(writer, request, err)
		return
	}

	category, err := model.GetCategoryByID(id)

	if err != nil {
		utils.NotFound(writer, request, err)
		return
	}

	if err := category.Delete(); err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}

	http.Redirect(writer, request, "/Admin/categories", http.StatusSeeOther)
}

func Category(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		utils.BadRequest(writer, request, err)
		return
	}
	category, err := model.GetCategoryByID(id)

	if err != nil {
		utils.NotFound(writer, request, err)
		return
	}

	category.Name = request.FormValue("name")
	categoryOuput := newCategoryOutput(&category)
	if categoryOuput.Check(); !categoryOuput.IsValid() {
		categoryOuput.ErrorsString = categoryOuput.Errors()
		utils.ExecuteTemplateWithoutNavbar(
			writer, request, categoryOuput, "templates/admin/base.html", "templates/admin/navbar.html",
			"templates/admin/AdminCategory.html",
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

	categoryOutput := newCategoryOutput(&category)

	if categoryOutput.Check(); !categoryOutput.IsValid() {
		categoryOutput.ErrorsString = categoryOutput.Errors()
		utils.ExecuteTemplateWithoutNavbar(
			w, r, categoryOutput, "templates/admin/base.html", "templates/admin/navbar.html",
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

func newCategoryOutput(category *model.Category) CategoryOutput {
	validator := validators.CategoryValidator{Category: category}
	return CategoryOutput{CategoryValidator: &validator}
}
