package handlers

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
	"sdu.store/server/validators"
	"time"
)

type Claims struct {
	User *model.User
	jwt.RegisteredClaims
}

var jwtKey = []byte("sdu.store.token")

func Login(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {

		CallHeaderHtml(writer, request)

		t, _ := template.ParseFiles("templates/login.gohtml")
		if err := t.Execute(writer, nil); err != nil {
			panic(err)
		}

	} else {

		Username := request.PostFormValue("username")
		Password := request.PostFormValue("password")
		user, err := validators.GetUserByUsername(Username)
		if err != nil {
			panic("User not Exists")
		}
		if user.Password == Password || CheckPasswordHash(Password, user.Password) {

			doLogin(writer, *user)

			http.Redirect(writer, request, "/index", http.StatusSeeOther)
		}
	}
}

func doLogin(writer http.ResponseWriter, user model.User) {

	expirationTime := time.Now().Add(24 * 60 * time.Minute)
	usr := &Claims{
		User: &user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, usr)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:    "session_token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	CurrentSession := model.Session{
		UserID:    user.ID,
		UUID:      tokenString,
		CreatedAt: time.Now(),
		DeletedAt: expirationTime,
		LastLogin: time.Now(),
		IP:        model.SetInet(),
	}

	var session model.Session

	session = CurrentSession

	server.DB.Create(&session)
}

func Logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			writer.WriteHeader(http.StatusUnauthorized)
		} else {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	var session model.Session
	server.DB.Last(&session)
	session.DeletedAt = time.Now()
	server.DB.Save(&session)

	cookie = &http.Cookie{
		Name:    "session_token",
		Expires: time.Now(),
	}

	http.SetCookie(writer, cookie)
	http.Redirect(writer, request, "/index", http.StatusSeeOther)
}

func SignUp(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		fmt.Println(err)
		t, _ := template.ParseFiles("templates/sign-up.html")
		err := t.Execute(writer, []string{"Server error"})
		if err != nil {
			return
		}
		return
	}
	if request.Method == "POST" {
		psw, err := HashPassword(request.PostFormValue("password"))
		if err != nil {
			panic(err)
		}
		user := model.User{
			Username: request.PostFormValue("username"),
			Email:    request.PostFormValue("email"),
			Password: psw,
		}
		v := validators.UserValidator{User: &user}
		if v.Check(); !v.IsValid() {
			t, _ := template.ParseFiles("templates/sign-up.html")
			fmt.Println(v.Errors())
			err := t.Execute(writer, v.Errors())
			if err != nil {
				return
			}
			return
		}
		if err := server.DB.Create(&user); err != nil {
			t, _ := template.ParseFiles("templates/sign-up.html")
			err := t.Execute(writer, []string{"User Created"})
			if err != nil {
				return
			}
			return
		}
		http.Redirect(writer, request, "/sign-in", 302)
	} else {
		t, _ := template.ParseFiles("templates/sign-up.html")
		err := t.Execute(writer, []string{"Two passwords doesn't match"})
		if err != nil {
			return
		}
		return
	}
}

func LoginPage(writer http.ResponseWriter, request *http.Request) {

}

func SignUpPage(writer http.ResponseWriter, request *http.Request) {
	t, _ := template.ParseFiles("templates/sign-up.html")
	err := t.Execute(writer, nil)
	if err != nil {
		return
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CheckCookie(writer http.ResponseWriter, request *http.Request) *Claims {
	claims := &Claims{}

	cookie, err := request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			writer.WriteHeader(http.StatusUnauthorized)
			return nil
		} else {
			writer.WriteHeader(http.StatusBadRequest)
			return nil
		}
	}
	if cookie != nil {

		key := cookie.Value

		token, err := jwt.ParseWithClaims(key, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				writer.WriteHeader(http.StatusUnauthorized)
				return nil
			}
			writer.WriteHeader(http.StatusBadRequest)
			return nil
		}
		if !token.Valid {
			writer.WriteHeader(http.StatusUnauthorized)
			return nil
		}
	}

	return claims
}

func CallHeaderHtml(writer http.ResponseWriter, request *http.Request) {

	t, _ := template.ParseFiles("templates/layouts/header.gohtml")
	user := &model.User{}

	claims := CheckCookie(writer, request)
	if claims != nil {
		user = claims.User
	}

	if err := t.Execute(writer, user.Email); err != nil {
		panic(err)
	}
}
