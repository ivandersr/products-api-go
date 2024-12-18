package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ivandersr/products-api-go/configs"
	"github.com/ivandersr/products-api-go/internal/entity"
	"github.com/ivandersr/products-api-go/internal/infra/database"
	"github.com/ivandersr/products-api-go/internal/infra/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	conf := configs.LoadConfig(".")

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProductDB(db)
	productHandler := handlers.NewProductHandler(productDB)
	userDB := database.NewUserDB(db)
	userHandler := handlers.NewUserHandler(userDB, conf.TokenAuth, conf.JWTExpiresIn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/products", productHandler.CreateProduct)
	r.Get("/products/{id}", productHandler.GetProduct)
	r.Put("/products/{id}", productHandler.UpdateProduct)
	r.Delete("/products/{id}", productHandler.DeleteProduct)
	r.Get("/products", productHandler.GetProducts)

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/token", userHandler.GetJWT)

	http.ListenAndServe(":8000", r)
}
