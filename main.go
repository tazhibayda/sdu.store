package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sdu.store/handlers"
	"sdu.store/handlers/admin"
	"sdu.store/server"
	"sdu.store/server/model"
)

func main() {

	loadFiles()

	restart := flag.Bool("dbRestart", false, "Restarting database")
	flag.Parse()

	if *restart {
		fmt.Println("restart ")

		server.DB.AutoMigrate(model.Session{}, model.User{}, model.Userdata{})
		server.DB.AutoMigrate(
			model.Category{}, model.Delivery{}, model.Item{}, model.Image{}, model.Product{}, model.ProductInfo{},
			model.Supplier{}, model.DeliveryItem{},
		)
		model.ConfigCategories()
	}
	if _, err := template.ParseGlob("templates/*.html"); err != nil {
		panic(err)
	}
	mux := http.NewServeMux()

	files := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	mux.HandleFunc("/", handlers.Index)

	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/account", handlers.Account)
	mux.HandleFunc("/sign-up-page", handlers.SignUpPage)
	mux.HandleFunc("/sign-up", handlers.SignUp)
	mux.HandleFunc("/login-page", handlers.Login)

	mux.HandleFunc("/Admin", admin.AdminServe)
	mux.HandleFunc("/Admin/login-page", admin.AdminLoginPage)
	mux.HandleFunc("/Admin/login", admin.AdminLogin)
	mux.HandleFunc("/Admin/logout", admin.AdminLogout)

	mux.HandleFunc("/Admin/add-user", admin.CreateUser)
	mux.HandleFunc("/Admin/users", admin.AdminUsers)
	mux.HandleFunc("/Admin/user", admin.User)
	mux.HandleFunc("/Admin/categories", admin.Categories)
	mux.HandleFunc("/Admin/category", admin.Category)
	mux.HandleFunc("/Admin/add-category", admin.CreateCategory)

	mux.HandleFunc("/Admin/products", admin.AdminProducts)
	mux.HandleFunc("/Admin/product/create", admin.CreateProduct)
	mux.HandleFunc("/Admin/product/delete/", admin.DeleteProduct)

	mux.HandleFunc("/Admin/user/delete/", admin.DeleteUser)
	mux.HandleFunc("/Admin/session", admin.GetAllSessions)
	mux.HandleFunc("/Admin/userdata", admin.AdminUserdata)
	mux.HandleFunc("/Admin/userdata/create", admin.CreateUserdata)
	mux.HandleFunc("/Admin/userdata/delete/", admin.DeleteUserdata)

	files = http.FileServer(http.Dir("images"))
	mux.Handle("/images/", http.StripPrefix("/images/", files))
	mux.HandleFunc("/test/images/upload", model.UploadImage)
	mux.HandleFunc("/test/images", model.ShowImages)
	err := http.ListenAndServe(":9090", mux)
	if err != nil {
		log.Fatal(err.Error())
	}

}

func loadFiles() []string {
	files, err := ioutil.ReadDir("templates")

	if err != nil {
		panic(err)
	}

	html := make([]string, 0)

	for _, file := range files {
		html = append(html, file.Name())
	}
	return html
}
