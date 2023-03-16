package model

import (
	"fmt"
	"net/http"
	"sdu.store/server"
	"strconv"
	"strings"
	"time"
)

type Product struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	CategoryID int64     `json:"category_id"`
	Price      float64   `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
}

type ProductOutput struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if r.Method == "POST" {
		categoryID, _ := strconv.Atoi(r.FormValue("category"))
		fmt.Println(int64(categoryID))
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		product = Product{
			Name:       r.FormValue("name"),
			CategoryID: int64(categoryID),
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
	vars := strings.Split(r.URL.Path, "/")
	productID := vars[len(vars)-1]
	product := Product{}
	server.DB.Where("ID = ?", productID).Delete(&product)
	http.Redirect(w, r, "/Admin/products", http.StatusSeeOther)
}
