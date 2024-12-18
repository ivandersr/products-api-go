package database

import (
	"github.com/ivandersr/products-api-go/internal/entity"
	"gorm.io/gorm"
)

type Product struct {
	DB *gorm.DB
}

func NewProductDB(db *gorm.DB) *Product {
	return &Product{DB: db}
}

func (p *Product) Create(product *entity.Product) error {
	return p.DB.Create(product).Error
}

func (p *Product) FindByID(id string) (*entity.Product, error) {
	var product entity.Product
	err := p.DB.First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *Product) FindAll(page, limit int, sort string) (*PaginatedResponse, error) {
	var products []entity.Product
	var err error
	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	if page != 0 && limit != 0 {
		offset := (page - 1) * limit
		err = p.DB.Limit(limit).Offset(offset).Order("created_at " + sort).Find(&products).Error
	} else {
		err = p.DB.Order("created_at " + sort).Find(&products).Error
	}

	if err != nil {
		return nil, err
	}

	response := &PaginatedResponse{
		Data:  products,
		Page:  page,
		Limit: limit,
	}

	return response, nil
}

func (p *Product) Update(product *entity.Product) error {
	_, err := p.FindByID(product.ID.String())
	if err != nil {
		return err
	}
	return p.DB.Save(product).Error
}

func (p *Product) Delete(id string) error {
	product, err := p.FindByID(id)
	if err != nil {
		return err
	}
	return p.DB.Delete(product).Error
}
