package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "test")
}

func enableCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func main() {
	router := mux.NewRouter()

	router.Use(enableCorsMiddleware)

	router.HandleFunc("/api/status", status)
	router.HandleFunc("/api/test", test).Methods("GET")

	http.ListenAndServe(":8080", router)
}
