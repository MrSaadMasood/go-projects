package router

import "net/http"

func healthRoutes(mux *http.ServeMux) {

	healthCheck := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
	}

	mux.HandleFunc("GET /health", healthCheck)
}
