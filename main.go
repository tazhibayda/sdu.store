package main

import (
	"flag"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
)

var DB *gorm.DB = server.DB

func main() {
	restart := flag.Bool("dbRestart", false, "Restarting database")
	flag.Parse()
	if *restart {
		DB.AutoMigrate(model.Session{}, model.User{}, model.Userdata{})
	}

	a := http.NewServeMux()

	a.HandleFunc("/Users", model.GetUsers)
	a.HandleFunc("/User", model.GetUserByID)
	a.HandleFunc("/Create", model.CreateUser)
	err := http.ListenAndServe(":9090", a)
	if err != nil {
		log.Fatal(err.Error())
	}

}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "asd")
	w.Write([]byte("Hello"))
	err := r.Write(w)
	if err != nil {
		return
	}
}
func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}
