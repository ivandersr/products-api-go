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

type Error struct {
	Message string `json:"message"`
}

func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

// GetJWT godoc
// @Summary 		 Generate JWT
// @Description 	 Generates a token to authenticate for requests
// @Tags 			 users
// @Accept 		 	 json
// @Produce		 	 json
// @Param 			 request  body 	   dto.GetJWTInput  true  "user request"
// @Success		 	 200	  {object} dto.GetJWTOutput
// @Failure		 	 500      {object} Error
// @Failure 		 401
// @Router 		 	 /users/token [post]
func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int)
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

	accessToken := dto.GetJWTOutput{AccessToken: tokenString}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

// Create user godoc
// @Summary 		 Create user
// @Description 	 Creates authenticatable user
// @Tags 			 users
// @Accept 		 	 json
// @Produce		 	 json
// @Param 			 request  body 	   dto.CreateUserInput  true  "user request"
// @Success		 	 201
// @Failure		 	 500      {object} Error
// @Router 		 	 /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	newUser, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	err = h.UserDB.Create(newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
