package server

import (
	"log"
	"main/internal/env"
	"main/internal/middlewares"
	"main/internal/router"
	"net/http"
)

func Run() {
	routes := router.Configure()

	port := env.Port
	if port == "" {
		port = ":5000"
	}

	log.Println("Starting the server on the port", port)
	server := &http.Server{
		Addr: port,
		Handler: middlewares.Error(
			middlewares.Logger(
				middlewares.JSONHeader(routes),
			),
		),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
