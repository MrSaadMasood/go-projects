package main

import (
	"fmt"
	"main/internal/database"
	_ "main/internal/env"
	"main/internal/server"
)

func main() {
	database.Connect()
	err := database.Setup(database.DB)
	if err != nil {
		panic(err)
	}
	server.Run()

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
}
