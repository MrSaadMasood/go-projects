package main

import (
	"fmt"
	"main/internal/contexter"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	secret := "6d43f22ad88a9d67001c1c58763b665d"
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"email": "admin@gmail.com"})
	token, err := claims.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("token sing error", err.Error())
		return
	}

	fmt.Println("the token is", token)

	type ExtendedClaims struct {
		contexter.TokenUser
		jwt.RegisteredClaims
	}
	var extendedClaims ExtendedClaims

	t, err := jwt.ParseWithClaims(token, &extendedClaims, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !t.Valid {
		fmt.Println("t is invalid", t)
		return
	}

	fmt.Println("the decoded token is", extendedClaims.Email)
}
