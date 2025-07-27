package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	JwtSecret   string
	PostgresURL string
	Port        string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	js := os.Getenv("JWT_SECRET")
	if js == "" {
		panic(fmt.Errorf("env for JWT_SECRET not provided"))
	}
	JwtSecret = js

	pu := os.Getenv("POSTGRES_URL")
	if pu == "" {
		panic(fmt.Errorf("env for POSTGRES_URL not provided"))
	}

	PostgresURL = pu

	port := os.Getenv("PORT")
	Port = port
}
