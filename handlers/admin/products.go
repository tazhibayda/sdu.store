package admin

import (
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/utils"
	"strconv"
	"strings"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	if r.Method == "POST" {
		categoryID, _ := strconv.Atoi(r.FormValue("category"))
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		product = model.Product{
			Name:       r.FormValue("Name"),
			CategoryID: uint(categoryID),
			Price:      price,
		}
		server.DB.Create(&product)
	} else {
		http.Redirect(w, r, "/Admin/products", http.StatusMethodNotAllowed)
		return
	}

	http.Redirect(w, r, "/Admin/products", http.StatusSeeOther)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := strings.Split(r.URL.Path, "/")
	productID := vars[len(vars)-1]
	product := model.Product{}
	server.DB.Where("ID = ?", productID).Delete(&product)
	http.Redirect(w, r, "/Admin/products", http.StatusSeeOther)
}

func Products(w http.ResponseWriter, r *http.Request) {
	_, err := utils.SessionStaff(w, r)
	if err != nil {
		http.Redirect(w, r, "Admin/login-page", http.StatusUnauthorized)
		return
	}
}
