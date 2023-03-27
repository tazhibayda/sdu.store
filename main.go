package main

import (
	"flag"
	"fmt"
	mux2 "github.com/gorilla/mux"
	"log"
	"net/http"
	"sdu.store/handlers"
	"sdu.store/handlers/admin"
	"sdu.store/server"
	"sdu.store/server/model"
)

func main() {

	restart := flag.Bool("dbRestart", false, "Restarting database")
	flag.Parse()

	if *restart {
		fmt.Println("restart ")

		server.DB.AutoMigrate(model.Session{}, model.User{}, model.Userdata{})
		server.DB.AutoMigrate(
			model.Category{}, model.Delivery{}, model.Product{}, model.Item{},
			model.Supplier{}, model.DeliveryItem{},
		)
		model.ConfigCategories()
	}
	mux := mux2.NewRouter()

	files := http.FileServer(http.Dir("static"))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", files))

	mux.HandleFunc("/", handlers.Index)

	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/account", handlers.Account)
	mux.HandleFunc("/sign-up-page", handlers.SignUpPage)
	mux.HandleFunc("/sign-up", handlers.SignUp)
	mux.HandleFunc("/login-page", handlers.LoginPage)

	mux.HandleFunc("/Admin", admin.StaffLoggingMiddleware(admin.AdminServe))
	mux.HandleFunc("/Admin/login", admin.AdminLoginPage).Methods("GET")
	mux.HandleFunc("/Admin/login", admin.AdminLogin).Methods("POST")
	mux.HandleFunc("/Admin/logout", admin.StaffLoggingMiddleware(admin.AdminLogout))

	mux.HandleFunc("/Admin/add-user", admin.AdminLoggingMiddleware(admin.CreateUser)).Methods("POST")
	mux.HandleFunc("/Admin/add-user", admin.AdminLoggingMiddleware(admin.AddUserPage)).Methods("GET")
	mux.HandleFunc("/Admin/users", admin.AdminLoggingMiddleware(admin.AdminUsers))
	mux.HandleFunc("/Admin/user", admin.AdminLoggingMiddleware(admin.User)).Methods("POST")
	mux.HandleFunc("/Admin/user", admin.AdminLoggingMiddleware(admin.UserPage)).Methods("GET")

	mux.HandleFunc("/Admin/categories", admin.StaffLoggingMiddleware(admin.Categories))
	mux.HandleFunc("/Admin/category", admin.StaffLoggingMiddleware(admin.CategoryPage)).Methods("GET")
	mux.HandleFunc("/Admin/category", admin.StaffLoggingMiddleware(admin.Category)).Methods("POST")
	mux.HandleFunc("/Admin/add-category", admin.StaffLoggingMiddleware(admin.CreateCategory)).Methods("POST")
	mux.HandleFunc("/Admin/add-category", admin.StaffLoggingMiddleware(admin.CreateCategoryPage)).Methods("GET")

	mux.HandleFunc("/Admin/products", admin.Products)
	mux.HandleFunc("/Admin/add-product", admin.CreateProduct)
	mux.HandleFunc("/Admin/product/delete/", admin.DeleteProduct)

	mux.HandleFunc("/Admin/session", admin.GetAllSessions)
	mux.HandleFunc("/Admin/userdata", admin.AdminUserdata)
	mux.HandleFunc("/Admin/userdata/create", admin.CreateUserdata)
	mux.HandleFunc("/Admin/userdata/delete/", admin.DeleteUserdata)

	files = http.FileServer(http.Dir("images"))
	mux.Handle("/images/", http.StripPrefix("/images/", files))
	err := http.ListenAndServe(":9090", mux)
	if err != nil {
		log.Fatal(err.Error())
	}

}
