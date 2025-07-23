package response

import "net/http"

func Error(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
}
