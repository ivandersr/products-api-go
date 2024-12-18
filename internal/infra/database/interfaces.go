package database

import "github.com/ivandersr/products-api-go/internal/entity"

type PaginatedResponse struct {
	Data  []entity.Product `json:"data"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}

type ProductInterface interface {
	Create(product *entity.Product) error
	FindAll(page, limit int, sort string) (*PaginatedResponse, error)
	FindByID(id string) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id string) error
}
