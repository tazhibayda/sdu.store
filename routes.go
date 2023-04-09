package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"sdu.store/handlers"
	"sdu.store/handlers/admin"
	"sdu.store/middlewares"
)

func routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(handlers.NotFoundHandler)
	router.MethodNotAllowed = http.HandlerFunc(handlers.NotAllowedMethod)

	standardMiddleware := alice.New(middlewares.EnableCORS, middlewares.LoggingRequest)
	adminMiddleware := alice.New(middlewares.AdminLoggingMiddleware)
	staffMiddleware := alice.New(middlewares.StaffLoggingMiddleware)

	router.ServeFiles("/static/*filepath", http.Dir("static"))

	router.HandlerFunc(http.MethodGet, "/", handlers.Index)

	router.HandlerFunc(http.MethodPost, "/login", handlers.Login)
	router.HandlerFunc(http.MethodGet, "/login", handlers.LoginPage)
	router.HandlerFunc(http.MethodGet, "/logout", handlers.Logout)
	router.HandlerFunc(http.MethodGet, "/account", handlers.Account)
	router.HandlerFunc(http.MethodGet, "/sign-up", handlers.SignUpPage)
	router.HandlerFunc(http.MethodPost, "/sign-up", handlers.SignUp)

	router.Handler(http.MethodGet, "/Admin", staffMiddleware.ThenFunc(admin.AdminServe))
	router.HandlerFunc(http.MethodGet, "/Admin/login", admin.AdminLoginPage)
	router.HandlerFunc(http.MethodPost, "/Admin/login", admin.AdminLogin)
	router.Handler(http.MethodGet, "/Admin/logout", staffMiddleware.ThenFunc(admin.AdminLogout))
	router.Handler(http.MethodPost, "/Admin/add-user", adminMiddleware.ThenFunc(admin.CreateUser))
	router.Handler(http.MethodGet, "/Admin/add-user", adminMiddleware.ThenFunc(admin.AddUserPage))
	router.Handler(http.MethodGet, "/Admin/users", adminMiddleware.ThenFunc(admin.AdminUsers))
	router.Handler(http.MethodPost, "/Admin/user", adminMiddleware.ThenFunc(admin.User))
	router.Handler(http.MethodGet, "/Admin/user", adminMiddleware.ThenFunc(admin.UserPage))
	router.Handler(http.MethodPost, "/Admin/delete-user", adminMiddleware.ThenFunc(admin.UserDelete))

	router.Handler(http.MethodGet, "/Admin/categories", staffMiddleware.ThenFunc(admin.Categories))
	router.Handler(http.MethodGet, "/Admin/category", staffMiddleware.ThenFunc(admin.CategoryPage))
	router.Handler(http.MethodPost, "/Admin/category", staffMiddleware.ThenFunc(admin.Category))
	router.Handler(http.MethodPost, "/Admin/delete-category", staffMiddleware.ThenFunc(admin.CategoryDelete))
	router.Handler(http.MethodPost, "/Admin/add-category", staffMiddleware.ThenFunc(admin.CreateCategory))
	router.Handler(http.MethodGet, "/Admin/add-category", staffMiddleware.ThenFunc(admin.CreateCategoryPage))

	router.Handler(http.MethodGet, "/Admin/add-product", staffMiddleware.ThenFunc(admin.AddProductPage))
	router.Handler(http.MethodPost, "/Admin/add-product", staffMiddleware.ThenFunc(admin.AddProduct))
	router.Handler(http.MethodGet, "/Admin/products", staffMiddleware.ThenFunc(admin.Products))
	router.Handler(http.MethodGet, "/Admin/product", staffMiddleware.ThenFunc(admin.ProductPage))
	router.Handler(http.MethodPost, "/Admin/product", staffMiddleware.ThenFunc(admin.Product))
	router.Handler(http.MethodPost, "/Admin/delete-product", staffMiddleware.ThenFunc(admin.DeleteProduct))

	router.Handler(http.MethodGet, "/Admin/add-item", staffMiddleware.ThenFunc(admin.AddItemPage))
	router.Handler(http.MethodPost, "/Admin/add-item", staffMiddleware.ThenFunc(admin.AddItem))
	router.Handler(http.MethodGet, "/Admin/items", staffMiddleware.ThenFunc(admin.Items))

	//
	//router.HandleFunc("/Admin/products", admin.Products)
	//router.HandleFunc("/Admin/add-product", admin.CreateProduct)
	//router.HandleFunc("/Admin/product/delete/", admin.DeleteProduct)
	//
	//router.HandleFunc("/Admin/session", admin.GetAllSessions)
	//router.HandleFunc("/Admin/userdata", admin.AdminUserdata)
	//router.HandleFunc("/Admin/userdata/create", admin.CreateUserdata)
	//router.HandleFunc("/Admin/userdata/delete/", admin.DeleteUserdata)

	//files = http.FileServer(http.Dir("images"))
	//router.Handle("/images/", http.StripPrefix("/images/", files))
	return standardMiddleware.Then(router)
}
