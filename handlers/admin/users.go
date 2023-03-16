package admin

import (
	"html/template"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/utils"
	"strconv"
	"strings"
	"time"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.SessionStaff(w, r); err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	var user model.User

	if r.Method == "POST" {
		login := r.FormValue("login")
		password := r.FormValue("password")
		username := r.FormValue("username")
		user = model.User{Email: login, Password: password, Username: username}
	}
	server.DB.Create(&user)
	//json.NewEncoder(w).Encode(user)
	http.Redirect(w, r, "/Admin/user", http.StatusSeeOther)

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.SessionStaff(w, r); err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
	}

	vars := strings.Split(r.URL.Path, "/")
	userID := vars[len(vars)-1]
	user := model.User{}
	server.DB.Where("ID = ?", userID).Delete(&user)
	http.Redirect(w, r, "/Admin/user", http.StatusSeeOther)
}

func AdminUserdata(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.SessionStaff(w, r); err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
	}
	tm, _ := template.ParseFiles("templates/Admin/AdminUserdata.gohtml")
	var userdata []model.Userdata
	server.DB.Find(&userdata)
	err := tm.Execute(w, userdata)
	if err != nil {
		return
	}
}

func AdminUsers(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.SessionStaff(w, r); err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	var user []model.User
	query := r.URL.Query()
	filters, presents := query["sort"]
	if !presents || len(filters) == 0 {
		server.DB.Find(&user)
	} else {
		if filters[0] == "" {

		}
	}
	tm, _ := template.ParseFiles("templates/Admin/AdminUser.gohtml")
	err := tm.Execute(w, user)
	if err != nil {
		return
	}
}

func CreateUserdata(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.SessionStaff(w, r); err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	uid, _ := strconv.ParseInt(r.FormValue("userid"), 10, 64)
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	phone := r.FormValue("phone")
	countrycode := r.FormValue("country_code")
	zip := r.FormValue("zip")
	birthday, _ := time.Parse("2006-01-02", r.FormValue("birthday"))

	userdata := model.Userdata{
		UserId:      uid,
		Firstname:   firstname,
		Lastname:    lastname,
		PhoneNumber: phone,
		CountryCode: countrycode,
		ZIPCode:     zip,
		Birthday:    birthday,
	}
	server.DB.Create(&userdata)
	http.Redirect(w, r, "/Admin/userdata", http.StatusSeeOther)

}

func DeleteUserdata(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.SessionStaff(w, r); err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusUnauthorized)
		return
	}

	vars := strings.Split(r.URL.Path, "/")
	userID := vars[len(vars)-1]
	var userdata model.Userdata
	server.DB.Where("User_ID = ?", userID).Delete(&userdata)
	http.Redirect(w, r, "/Admin/userdata", http.StatusSeeOther)

}
