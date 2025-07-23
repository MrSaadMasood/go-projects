package main

import (
	"main/internal/database"
	"main/internal/server"
)

func main() {
	db := database.Connect()
	database.Setup(db)
	server.Run()
}
