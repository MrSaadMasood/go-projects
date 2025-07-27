package middlewares

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	requestLogger := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Since(start)
		log.Println(r.Method, r.URL.Path, end, w.Header().Get("status"))
	}
	return http.HandlerFunc(requestLogger)
}
