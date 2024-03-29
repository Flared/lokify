package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/flared/lokify/pkg/api"
)

var (
	port        = flag.String("port", "8080", "Define on which port the server will run. Default: 8080")
	lokiBaseUrl = flag.String("loki-base-url", "http://loki:3100", "Define loki base url. Defautl: http://loki:3100")
)

func main() {
	flag.Parse()

	ctx := api.NewContext(&http.Client{}, *lokiBaseUrl)
	router := api.NewRouter(ctx)
	addr := fmt.Sprintf(":%v", *port)

	fmt.Printf("Start server on port %v\n", *port)
	log.Fatal(http.ListenAndServe(addr, router))
}
