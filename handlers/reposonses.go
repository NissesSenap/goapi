package handlers

import "net/http"

func postError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
