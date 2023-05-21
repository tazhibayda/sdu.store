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
		server.DB.Migrator().DropTable(
			model.Comment{}, model.Rating{}, model.DeliveryItem{}, model.Supplier{}, model.Item{}, model.Product{},
			model.Delivery{}, model.Category{}, model.Userdata{}, model.User{}, model.Purchase{},
		)
		server.DB.AutoMigrate(model.User{}, model.Userdata{},
			model.Category{}, model.Delivery{}, model.Product{}, model.Item{},
			model.Supplier{}, model.DeliveryItem{}, model.Rating{}, model.Comment{},
			model.Purchase{},
		)
		model.ConfigCategories()
	}

	err := http.ListenAndServe(":9090", routes())
	if err != nil {
		log.Fatal(err.Error())
	}

}
