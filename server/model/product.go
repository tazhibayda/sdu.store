package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"sdu.store/server"
	"strconv"
	"strings"
)

type Product struct {
	gorm.Model
	Name           string         `json:"name"`
	CategoryID     int            `json:"category"`
	Price          float64        `json:"price"`
	Images         pq.StringArray `gorm:"type:text[]" json:"images"`
	Sizes          pq.StringArray `gorm:"type:text[]" json:"sizes"`
	Colors         pq.StringArray `gorm:"type:text[]" json:"colors"`
	Description    string         `json:"description" input:"string"`
	Items          []Item
	Rating         float64 `gorm:"default:0"`
	AmountRatings  int64   `gorm:"default:0"`
	Comments       []Comment
	AmountComments int64
}

func (this *Product) Delete() error {
	return server.DB.Delete(this).Error
}
func (this *Product) Create() error {
	return server.DB.Create(this).Error
}

func (this *Product) Update() error {
	return server.DB.Save(this).Error
}

func GetProductByID(productID int) (Product, error) {
	var product Product
	err := server.DB.Where("ID=?", productID).Preload("Comments").Find(&product).Error
	return product, err
}

func ParseProduct(product *Product, request *http.Request) error {
	var err error
	// parsing name
	product.Name = request.PostFormValue("name")

	// parsing category id
	product.CategoryID, err = strconv.Atoi(request.PostFormValue("category"))
	if err != nil {
		return err
	}

	// parsing price
	product.Price, err = strconv.ParseFloat(request.PostFormValue("price"), 64)
	if err != nil {
		return err
	}

	// parsing images
	product.Images, err = retrieveFilesNameSlice(request, "images")
	if err != nil {
		return err
	}

	// parsing sizes
	product.Sizes = retrieveSliceFromRequest(request, "sizes")

	// parsing colors
	product.Colors = retrieveSliceFromRequest(request, "colors")

	// parsing description
	product.Description = request.PostFormValue("description")
	return nil
}

func retrieveFilesNameSlice(request *http.Request, name string) ([]string, error) {
	err := request.ParseMultipartForm(2000000000) // grab the multipart form
	if err != nil {
		return nil, err
	}

	formdata := request.MultipartForm // ok, no problem so far, read the Form data

	//get the *fileheaders
	files := formdata.File[name]

	ans := []string{}

	for i, _ := range files { // loop through the files one by one
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			return nil, err
		}

		filename := "images/" + files[i].Filename

		out, err := os.Create(filename)

		defer out.Close()
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(out, file) // file not files[i] !

		if err != nil {
			return nil, err
		}

		ans = append(ans, "/"+filename)

	}
	return ans, nil
}

func retrieveSliceFromRequest(request *http.Request, name string) []string {
	mp := make(map[string]bool)
	strs := strings.Split(request.PostFormValue(name), " ")
	for _, str := range strs {
		if str != "" {
			mp[str] = true
		}
	}
	ans := []string{}
	for key, _ := range mp {
		ans = append(ans, key)
	}
	return ans
}
