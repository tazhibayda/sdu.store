package model

import (
	"github.com/lib/pq"
	"net/http"
	"sdu.store/server"
	"strconv"
	"strings"
	"time"
)

type ProductOutput struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
}

type Product struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	CategoryID  int            `json:"category_id"`
	Price       float64        `json:"price"`
	Images      pq.StringArray `gorm:"type:text[]"`
	Sizes       pq.StringArray `gorm:"type:text[]"`
	Colors      pq.StringArray `gorm:"type:text[]"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if r.Method == "POST" {
		categoryID, _ := strconv.Atoi(r.FormValue("category"))
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		product = Product{
			Name:       r.FormValue("name"),
			CategoryID: categoryID,
			Price:      price,
			CreatedAt:  time.Now(),
		}
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
