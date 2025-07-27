package middlewares

import (
	"encoding/json"
	"net/http"
)

func Error(h http.Handler) http.Handler {
	handleError := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {

				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]any{
					"message": "Internal server error",
					"err":     err,
				})
			}
		}()

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handleError)
}
