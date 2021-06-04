package main

import (
	"log"
	"net/http"

	"github.com/flared/lokify/pkg/api"
)

func main() {
	router := api.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
