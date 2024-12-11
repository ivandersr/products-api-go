package database

import (
	"fmt"
	"testing"

	"github.com/ivandersr/products-api-go/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Product 01", 80)
	productDB := NewProduct(db)

	err = productDB.Create(product)
	assert.Nil(t, err)

	var productFound entity.Product
	err = db.First(&productFound, "id = ?", product.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, "Product 01", productFound.Name)
	assert.Equal(t, 80, productFound.Price)
	assert.NotNil(t, product.CreatedAt)
}

func TestFindProductByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Product 01", 80)
	productDB := NewProduct(db)

	db.Create(product)

	productFound, err := productDB.FindByID(product.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, product.Price, productFound.Price)
}

func TestFindAllProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	var products []entity.Product
	for i := 1; i <= 10; i++ {
		product, _ := entity.NewProduct(fmt.Sprintf("Product %d", i), i*10)
		products = append(products, *product)
	}
	db.Create(products)
	productDB := NewProduct(db)
	foundProducts, err := productDB.FindAll(0, 0, "")
	assert.Nil(t, err)
	assert.Len(t, foundProducts, len(products))
}

func TestFindAllProductsWithPagination(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	for i := 1; i <= 10; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), i*10)
		assert.NoError(t, err)
		db.Create(product)
	}
	productDB := NewProduct(db)
	foundProducts, err := productDB.FindAll(2, 4, "")
	assert.Nil(t, err)
	assert.Len(t, foundProducts, 4)
	assert.Equal(t, "Product 5", foundProducts[0].Name)
	assert.Equal(t, "Product 8", foundProducts[3].Name)
}

func TestUpdateProdcut(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Product 01", 80)
	productDB := NewProduct(db)

	db.Create(product)
	product.Name = "Updated Product 01"
	product.Price = 100
	var foundProduct entity.Product
	err = productDB.Update(product)
	assert.Nil(t, err)
	err = db.First(&foundProduct, "id = ?", product.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, "Updated Product 01", foundProduct.Name)
	assert.Equal(t, 100, foundProduct.Price)
}

func TestDeleteProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Product 01", 80)
	productDB := NewProduct(db)

	db.Create(product)

	err = productDB.Delete(product.ID.String())
	assert.NoError(t, err)

	_, err = productDB.FindByID(product.ID.String())
	assert.Error(t, err)
	assert.Equal(t, "record not found", err.Error())
}
