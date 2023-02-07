package server

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB = ConnectDB()

func ConnectDB() *gorm.DB {

	_ = pq.Driver{}
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		Username, Password, Host, Port, Dbname)

	sqlDB, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to database! ")

	return db
}
