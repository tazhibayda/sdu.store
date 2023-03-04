package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sdu.store/handlers"
	"sdu.store/server"
	"sdu.store/server/model"
)

func main() {

	loadFiles()

	restart := flag.Bool("dbRestart", false, "Restarting database")
	flag.Parse()

	if *restart {
		fmt.Println("restart ")
		model.ConfigCategories()
		server.DB.AutoMigrate(model.Session{}, model.User{}, model.Userdata{})
		server.DB.AutoMigrate(
			model.Category{}, model.Delivery{}, model.Item{}, model.Image{}, model.Product{}, model.ProductInfo{},
			model.Supplier{}, model.DeliveryItem{},
		)
	}
	if _, err := template.ParseGlob("templates/*.gohtml"); err != nil {
		panic(err)
	}
	mux := http.NewServeMux()

	files := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	// Request for postman
	mux.HandleFunc("/request/users", model.GetUsers)
	mux.HandleFunc("/request/user", model.GetUserByID)

	mux.HandleFunc("/index", handlers.Index)

	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/sign-up-page", handlers.SignUpPage)
	mux.HandleFunc("/sign-up", handlers.SignUp)
	mux.HandleFunc("/login-page", handlers.LoginPage)

	mux.HandleFunc("/Admin", model.AdminServe)
	mux.HandleFunc("/Admin/user/create", model.CreateUser)
	mux.HandleFunc("/Admin/user", model.AdminUsers)
	mux.HandleFunc("/Admin/categories", model.AdminCategories)
	mux.HandleFunc("/Admin/category/create", model.CreateCategory)
	mux.HandleFunc("/Admin/category/delete/", model.DeleteCategory)

	mux.HandleFunc("/Admin/user/delete/", model.DeleteUser)
	mux.HandleFunc("/Admin/session", model.GetAllSessions)
	mux.HandleFunc("/Admin/userdata", model.AdminUserdata)
	mux.HandleFunc("/Admin/userdata/create", model.CreateUserdata)
	mux.HandleFunc("/Admin/userdata/delete/", model.DeleteUserdata)

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
