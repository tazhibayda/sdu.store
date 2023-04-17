package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sdu.store/server"
	"sdu.store/server/model"
)

func main() {

	restart := flag.Bool("dbRestart", false, "Restarting database")
	flag.Parse()

	if *restart {
		fmt.Println("restart ")

		server.DB.AutoMigrate(model.User{}, model.Userdata{})
		server.DB.AutoMigrate(
			model.Category{}, model.Delivery{}, model.Product{}, model.Item{},
			model.Supplier{}, model.DeliveryItem{}, model.Rating{}, model.Comment{},
		)
		model.ConfigCategories()
	}

	err := http.ListenAndServe(":9090", routes())
	if err != nil {
		log.Fatal(err.Error())
	}

}
