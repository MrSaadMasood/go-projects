package server

import (
	"main/internal/middlewares"
	"main/internal/router"
	"net/http"
)

func Run() {
	routes := router.Routes()
	server := http.Server{Addr: ":5000", Handler: middlewares.JSONHeader(routes)}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
