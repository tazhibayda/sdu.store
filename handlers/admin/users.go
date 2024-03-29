package admin

import (
	"html/template"
	"net/http"
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
	email := r.FormValue("email")
	password := r.FormValue("password")
	username := r.FormValue("username")
	isStaff := r.FormValue("staff") == "on"
	isAdmin := r.FormValue("admin") == "on"
	user := model.User{Email: email, Password: password, Username: username, Is_staff: isStaff, Is_admin: isAdmin}
	v := validators.NewUserValidator(&user)
	if v.Check(); !v.IsValid() {
		utils.ExecuteTemplateWithoutNavbar(
			w, r, v.Errors(), "templates/admin/base.html", "templates/admin/navbar.html",
			"templates/admin/AdminAddUser.html",
		)
		return
	}
	if err := server.DB.Create(&user).Error; err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}
	http.Redirect(w, r, "/Admin/users", http.StatusSeeOther)
	return
}

func AddUserPage(writer http.ResponseWriter, request *http.Request) {
	utils.ExecuteTemplateWithoutNavbar(
		writer, request, nil, "templates/admin/base.html", "templates/admin/navbar.html",
		"templates/admin/AdminAddUser.html",
	)
}

func AdminUserdata(w http.ResponseWriter, r *http.Request) {
	_, err := utils.SessionAdmin(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	tm, err := template.ParseFiles("templates/Admin/AdminUserdata.gohtml")
	if err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}
	var userdata []model.Userdata
	server.DB.Find(&userdata)
	err = tm.Execute(w, userdata)
	if err != nil {
		return
	}
}

func AdminUsers(w http.ResponseWriter, r *http.Request) {
	var users []model.User
	if err := server.DB.Find(&users).Error; err != nil {
		utils.ServerErrorHandler(w, r, err)
		return
	}
	hasFilter, userTable, err := HasFilterUserTable(r)
	if err != nil {
		utils.BadRequest(w, r, err)
		return
	}
	if hasFilter {
		sortUserTable(users, &userTable)
	} else {
		userTable.Users = users
	}
	utils.ExecuteTemplateWithoutNavbar(
		w, r, userTable, "templates/admin/base.html", "templates/admin/navbar.html", "templates/admin/AdminUsers.html",
	)
}

func UserPage(writer http.ResponseWriter, request *http.Request) {

	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		utils.BadRequest(writer, request, err)
		return
	}

	user, err := model.GetUserByID(int64(id))
	if err != nil {
		utils.NotFound(writer, request, err)
		return
	}
	utils.ExecuteTemplateWithoutNavbar(
		writer, request, user, "templates/admin/base.html", "templates/admin/navbar.html",
		"templates/admin/AdminUser.html",
	)
	return
}

func UserDelete(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		utils.BadRequest(writer, request, err)
		return
	}
	user, err := model.GetUserByID(int64(id))

	if err != nil {
		utils.NotFound(writer, request, err)
		return
	}

	if err := user.Delete(); err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}

	http.Redirect(writer, request, "/Admin/users", http.StatusSeeOther)
}

func User(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		utils.BadRequest(writer, request, err)
		return
	}
	user, err := model.GetUserByID(int64(id))

	if err != nil {
		utils.NotFound(writer, request, err)
		return
	}

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
	if err := user.Update(); err != nil {
		utils.ServerErrorHandler(writer, request, err)
		return
	}

	http.Redirect(writer, request, "/Admin/users", http.StatusSeeOther)
	return

}

func CreateUserdata(w http.ResponseWriter, r *http.Request) {
	_, err := utils.SessionAdmin(w, r)
	if err != nil {
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
	_, err := utils.SessionStaff(w, r)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusUnauthorized)
		return
	}

	vars := strings.Split(r.URL.Path, "/")
	userID := vars[len(vars)-1]
	var userdata model.Userdata
	server.DB.Where("User_ID = ?", userID).Delete(&userdata)
	http.Redirect(w, r, "/Admin/userdata", http.StatusSeeOther)

}

func HasFilterUserTable(request *http.Request) (hasFilter bool, filter UserTable, err error) {
	if search := request.FormValue("search"); search != "" {
		hasFilter = true
		filter.Search = search
	}
	if staff := request.FormValue("staff"); staff != "" {
		hasFilter = true
		filter.StaffStatus, err = strconv.Atoi(staff)
		if err != nil {
			return
		}
	}
	if admin := request.FormValue("admin"); admin != "" {
		hasFilter = true
		filter.AdminStatus, err = strconv.Atoi(admin)
		if err != nil {
			return
		}
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
