package main

import (
	"GoAPI/configs"
	"GoAPI/internal/entity"
	"GoAPI/internal/infra/database"
	"GoAPI/internal/infra/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDB := database.NewProductDB(db)
	productH := handlers.NewProductHandler(productDB)

	http.HandleFunc("/products", productH.CreateProduct)
	http.ListenAndServe(":8000", nil)
}
