package env

import (
	"fmt"
	"os"
)

var JwtSecret string

func init() {
	js := os.Getenv("JWT_SECRET")
	if js == "" {
		panic(fmt.Errorf("env for JWT_SECRET not provided"))
	}

	JwtSecret = js
}
