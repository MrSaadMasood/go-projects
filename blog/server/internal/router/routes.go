package router

import (
	"net/http"
)

func Routes() *http.ServeMux {
	authMux := authRoutes()
	userMux := userRoutes()

	mux := http.NewServeMux()

	mux.Handle("/auth", http.StripPrefix("/auth", authMux))
	mux.Handle("/user", http.StripPrefix("/user", userMux))

	return mux

}
