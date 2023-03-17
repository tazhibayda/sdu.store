package admin

import (
	"html/template"
	"net/http"
	"sdu.store/handlers"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/server/validators"
	"sdu.store/utils"
	"strconv"
	"strings"
	"time"
)

type UserTable struct {
	Users       []model.User
	Search      string
	StaffStatus int
	AdminStatus int
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user, err := utils.SessionStaff(w, r)
	if err != nil {
		http.Redirect(w, r, "/Admin/login-page", http.StatusTemporaryRedirect)
		return
	}

	if !user.Is_admin {
		http.Redirect(w, r, "Admin/login-page", http.StatusTemporaryRedirect)
		return
	}

	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")
		username := r.FormValue("username")
		isStaff := r.FormValue("staff") == "on"
		isAdmin := r.FormValue("admin") == "on"
		user := model.User{Email: email, Password: password, Username: username, Is_staff: isStaff, Is_admin: isAdmin}
		v := validators.UserValidator{User: &user}
		if v.Check(); !v.IsValid() {
			tm, _ := template.ParseFiles(
				"templates/admin/base.html", "templates/admin/navbar.html", "templates/admin/AdminAddUser.html",
			)
			tm.ExecuteTemplate(w, "base", v.Errors())
			return
		}
		user.Password, _ = handlers.HashPassword(user.Password)
		server.DB.Create(&user)
		http.Redirect(w, r, "/Admin/users", http.StatusSeeOther)
		return
	} else {
		tm, _ := template.ParseFiles(
			"templates/admin/base.html", "templates/admin/navbar.html", "templates/admin/AdminAddUser.html",
		)
		tm.ExecuteTemplate(w, "base", nil)
	}
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
		http.Redirect(w, r, "/Admin/login-page", http.StatusTemporaryRedirect)
		return
	}
	var users []model.User
	server.DB.Find(&users)
	hasFilter, userTable := HasFilterUserTable(r)
	if hasFilter {
		sortUserTable(users, &userTable)
	} else {
		userTable.Users = users
	}
	tm, err := template.ParseFiles(
		"templates/admin/base.html", "templates/admin/navbar.html", "templates/admin/AdminUsers.html",
	)
	err = tm.ExecuteTemplate(w, "base", userTable)
	if err != nil {
		return
	}
}

func User(writer http.ResponseWriter, request *http.Request) {
	user, err := utils.SessionStaff(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/Admin/login-page", http.StatusTemporaryRedirect)
		return
	}
	if !user.Is_admin {
		http.Redirect(writer, request, "/Admin/login-page", http.StatusTemporaryRedirect)
		return
	}

	if request.Method == "GET" {
		tm, _ := template.ParseFiles(
			"templates/admin/base.html", "templates/admin/navbar.html", "templates/admin/AdminUser.html",
		)
		err = tm.ExecuteTemplate(writer, "base", user)
		return
	}
	if request.Method == "POST" {
		if "1" == request.URL.Query().Get("delete") {
			user.Delete()
		} else {
			if "on" == request.FormValue("admin") && !user.Is_admin {
				user.Is_admin = true
			} else if "" == request.FormValue("admin") && user.Is_admin {
				user.Is_admin = false
			}
			if "on" == request.FormValue("staff") && !user.Is_staff {
				user.Is_staff = true
			} else if "" == request.FormValue("staff") && user.Is_staff {
				user.Is_staff = false
			}
			user.Update()
		}
		http.Redirect(writer, request, "/Admin/users", http.StatusTemporaryRedirect)
		return
	}

}

func CreateUserdata(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.SessionStaff(w, r); err != nil {
		http.Redirect(w, r, "/Admin/login-page", http.StatusTemporaryRedirect)
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

func HasFilterUserTable(request *http.Request) (hasFilter bool, filter UserTable) {
	if search := request.FormValue("search"); search != "" {
		hasFilter = true
		filter.Search = search
	}
	if staff := request.FormValue("staff"); staff != "" {
		hasFilter = true
		filter.StaffStatus, _ = strconv.Atoi(staff)
	}
	if admin := request.FormValue("admin"); admin != "" {
		hasFilter = true
		filter.AdminStatus, _ = strconv.Atoi(admin)
	}
	return
}

func sortUserTable(users []model.User, table *UserTable) {
	ans := []model.User{}
	for _, user := range users {
		if isValidUser(user, table) {
			ans = append(ans, user)
		}
	}
	table.Users = ans
}

func isValidUser(user model.User, table *UserTable) bool {
	if table.Search != "" {
		contains := strings.Contains(user.Username, table.Search)
		if !contains {
			return false
		}
	}
	if table.StaffStatus != 0 {
		if table.StaffStatus == 1 && !user.Is_staff {
			return false
		}
		if table.StaffStatus == 2 && user.Is_staff {
			return false
		}
	}
	if table.AdminStatus != 0 {
		if table.AdminStatus == 1 && !user.Is_admin {
			return false
		}
		if table.AdminStatus == 2 && user.Is_admin {
			return false
		}
	}
	return true
}
