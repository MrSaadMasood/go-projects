package middlewares

import (
	"main/pkg/response"
	"net/http"
)

func JSONHeader(h http.Handler) http.Handler {
	applicationJSONHeader := func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("content-type")
		if header == "" {
			response.Error(w, "content type header not provided", http.StatusBadRequest)
			return
		}

		if header != "application/json" {
			response.Error(w, "invalid content type provided:"+header, http.StatusBadRequest)
			return
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(applicationJSONHeader)
}
