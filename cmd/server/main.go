package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
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
	userHandler := handlers.NewUserHandler(userDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.WithValue("jwt", conf.TokenAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", conf.JWTExpiresIn))
	r.Use(middleware.Recoverer) // Graceful panic absorption with stack trace log, keeps API online
	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(conf.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
		r.Get("/", productHandler.GetProducts)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Post("/token", userHandler.GetJWT)
	})

	http.ListenAndServe(":8000", r)
}
