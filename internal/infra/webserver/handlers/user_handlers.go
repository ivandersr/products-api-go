package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/ivandersr/products-api-go/internal/dto"
	"github.com/ivandersr/products-api-go/internal/entity"
	"github.com/ivandersr/products-api-go/internal/infra/database"
)

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int64)
	var user dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	foundUser, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !foundUser.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": foundUser.ID.String(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: tokenString,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	newUser, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.UserDB.Create(newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
