package main

import (
	"GoAPI/configs"
	_ "GoAPI/docs"
	"GoAPI/internal/entity"
	"GoAPI/internal/infra/database"
	"GoAPI/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

// @title Go API
// @version 1.0
// @description Product API
// @termsOfService http://swagger.io/terms/

// @contact.name Pedro William
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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
	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productH.CreateProduct)
		r.Get("/{id}", productH.GetProduct)
		r.Get("/", productH.GetProducts)
		r.Put("/{id}", productH.UpdateProduct)
		r.Delete("/{id}", productH.DeleteProduct)
	})
	r.Post("/users", userH.CreateUser)
	r.Post("/users/generateToken", userH.GetJWT)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", r)
}
