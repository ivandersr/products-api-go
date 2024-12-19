package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/ivandersr/products-api-go/configs"
	_ "github.com/ivandersr/products-api-go/docs"
	"github.com/ivandersr/products-api-go/internal/entity"
	"github.com/ivandersr/products-api-go/internal/infra/database"
	"github.com/ivandersr/products-api-go/internal/infra/webserver/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title            		   Products API Go
// @version          		   1.0
// @description      		   Products API with Authentication
// @termsOfService   		   http://swagger.io/terms

// @contact.name     		   Ivander
// @contact.url      		   https://linkedin.com/in/ivandersr
// @contact.email    		   ivandersalvadorruiz@gmail.com

// @license.name     		   MIT
// @license.url      		   https://opensource.org/license/mit

// @host                       localhost:8000
// @BasePath                   /
// @securityDefinitions.apiKey ApiKeyAuth
// @in                         header
// @name                       Authorization
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

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", r)
}
