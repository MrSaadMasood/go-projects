package router

import (
	"main/internal/services"
	"net/http"
)

func authRoutes(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("POST /auth/register", services.Signup)
	mux.HandleFunc("POST /auth/login", services.Login)
	return mux
}
