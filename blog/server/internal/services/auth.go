package services

import (
	"encoding/json"
	"fmt"
	"main/internal/database"
	"main/internal/env"
	"main/pkg/response"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	type SignUpBody struct {
		Name     string `json:"name" validate:"required,min=3,max=50"`
		Password string `json:"password" validate:"required,min=10,max=30"`
		Email    string `json:"email" validate:"email,required"`
	}

	var signupBody SignUpBody
	err := json.NewDecoder(r.Body).Decode(&signupBody)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(&signupBody)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupBody.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createUserQuery := "INSERT INTO users (name, email, password, created_at) values ($1, $2, $3, $4)"
	user, err := database.DB.Query(createUserQuery, signupBody.Name, signupBody.Email, hashedPassword, time.Now().UTC())
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("the user created is", user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User Created Successfully"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	type LoginUser struct {
		Email    string `json:"email" validate:"email,required"`
		Password string `json:"password" validate:"required"`
	}

	var loginUserBody LoginUser
	err := json.NewDecoder(r.Body).Decode(&loginUserBody)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(&loginUserBody)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var email string
	var password string
	err = database.DB.QueryRow("SELECT email, password FROM users WHERE email = $1", loginUserBody.Email).Scan(&email, &password)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(loginUserBody.Password))
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"email": email})
	token, err := claims.SignedString([]byte(env.JwtSecret))
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data": map[string]string{
			"token": token,
		},
	})
}
