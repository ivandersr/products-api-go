package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/ivandersr/products-api-go/internal/dto"
	"github.com/ivandersr/products-api-go/internal/entity"
	"github.com/ivandersr/products-api-go/internal/infra/database"
	entityPkg "github.com/ivandersr/products-api-go/pkg/entity"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// Create product godoc
// @Summary 		 Create product
// @Description 	 Creates a new product
// @Tags 			 products
// @Accept 		 	 json
// @Produce		 	 json
// @Param 			 request  body 	   dto.CreateProductInput  true  "product request"
// @Success		 	 201
// @Failure		 	 500      {object} Error
// @Router 		 	 /products [post]
// @Security 		 ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	newProduct, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.ProductDB.Create(newProduct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetProduct godoc
// @Summary 		 Find a product
// @Description 	 Returns a product by its ID
// @Tags 			 products
// @Produce		 	 json
// @Param			 id	  	  path   	string  	true		"product ID"   Format(uuid)
// @Success		 	 200 	  {object}  entity.Product
// @Failure			 401
// @Failure			 404
// @Failure		 	 500      {object}  Error
// @Router 		 	 /products/{id} [get]
// @Security 		 ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct godoc
// @Summary			Updates a product
// @Description		Updates a product data by its ID
// @Tags			products
// @Accept			json
// @Produce			json
// @Param			id	  	  path   	string  	true		"product ID"   Format(uuid)
// @Param 			request   body 	    dto.CreateProductInput	true 		   "product request"
// @Success			200
// @Failure			401
// @Failure			404
// @Failure 		500		  {object}  Error
// @Router		    /products/{id} [put]
// @Security 		ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.ProductDB.Update(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteProduct godoc
// @Summary			Deletes a product
// @Description		Deletes a product by its ID
// @Tags			products
// @Produce			json
// @Param			id	  	  path   	string  	true		"product ID"   Format(uuid)
// @Success			204
// @Failure			401
// @Failure			404
// @Failure 		500		  {object}  Error
// @Router		    /products/{id} [delete]
// @Security 		ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetProducts godoc
// @Summary 		 List Porducts
// @Description 	 Returns a list of all products with optional pagination
// @Tags 			 products
// @Accept 		 	 json
// @Produce		 	 json
// @Param			 page	  query   	string   	   false      "page number"
// @Param			 limit	  query   	string   	   false      "items per page"
// @Success		 	 200 	  {object}  database.PaginatedResponse
// @Failure			 401
// @Failure		 	 500      {object}  Error
// @Router 		 	 /products [get]
// @Security 		 ApiKeyAuth
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 0
	}
	sort := r.URL.Query().Get("sort")
	products, err := h.ProductDB.FindAll(page, limit, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
