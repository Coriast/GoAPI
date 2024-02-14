package main

import (
	"GoAPI/configs"
	"GoAPI/internal/entity"
	"GoAPI/internal/infra/database"
	"GoAPI/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	configs, err := configs.LoadConfig(".")
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

	userDB := database.NewUserDB(db)
	userH := handlers.NewUserHandler(userDB, configs.TokenAuth, configs.JWTExpiresIn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/products", productH.CreateProduct)
	r.Get("/products/{id}", productH.GetProduct)
	r.Get("/products", productH.GetProducts)
	r.Put("/products/{id}", productH.UpdateProduct)
	r.Delete("/products/{id}", productH.DeleteProduct)

	r.Post("/users", userH.CreateUser)
	r.Post("/users/generateToken", userH.GetJWT)

	http.ListenAndServe(":8000", r)
}
