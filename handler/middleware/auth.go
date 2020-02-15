package middleware

import (
	"net/http"

	"fileStore/handler"
)

// HTTPInterceptor is a middleware, it will check request token is valid
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	// receive a http.HandlerFunc, return a new http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		username := r.Form.Get("username")
		token := r.Form.Get("token")

		if len(username) < 3 || !handler.IsTokenValid(token, username) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		h(w, r)
	})
}
