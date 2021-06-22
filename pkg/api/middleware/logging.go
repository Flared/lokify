package middleware

import (
	"log"
	"net/http"
)

func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		h.ServeHTTP(w, r)
	})
}
