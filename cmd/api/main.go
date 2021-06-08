package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/flared/lokify/pkg/api"
)

var (
	port = flag.String("port", "8080", "Define on which port the server will run. Default: 8080")
)

func main() {
	flag.Parse()

	fmt.Printf("Start server on port %v\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), api.NewRouter()))
}
