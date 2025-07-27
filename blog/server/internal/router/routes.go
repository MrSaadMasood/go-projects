package router

import (
	"net/http"
)

func Configure() *http.ServeMux {
	router := http.NewServeMux()
	healthRoutes(router)
	authRoutes(router)
	userRoutes(router)
	return router
}
