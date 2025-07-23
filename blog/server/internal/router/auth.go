package router

import (
	"main/internal/services"
	"net/http"
)

func authRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", services.Signup)
	mux.HandleFunc("POST /login", services.Login)
	return mux
}
