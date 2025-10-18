// Package env acts as central place for env access
package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var MaxConcurrency int

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	maxConcurrency := os.Getenv("MAX_CONCURRENCY")
	if maxConcurrency == "" {
		MaxConcurrency = 30
	} else {
		MaxConcurrency, err = strconv.Atoi(maxConcurrency)
		if err != nil {
			panic(err)
		}
	}
}
