package router

import "net/http"

func userRoutes() *http.ServeMux {

	mux := http.NewServeMux()
	return mux

}
