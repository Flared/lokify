package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/flared/lokify/pkg/api/appctx"
	"github.com/flared/lokify/pkg/api/middleware"
	"github.com/gorilla/mux"
)

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func queryHandler(ctx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		resp, err := ctx.Loki.Query(vars["query"])
		if err != nil {
			log.Printf("Loki Query error, %v", err)
			w.WriteHeader(500)
			return
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("json Encode error, %v", err)
			w.WriteHeader(500)
			return
		}
	}

	handler := http.HandlerFunc(handlerFunc)
	return middleware.JSONHeaders(handler)
}

func index(w http.ResponseWriter, r *http.Request) {
	s := make(map[string]string)
	s["base_url"] = "http://localhost:8080"

	distDir := "../ui/build/index.html"
	if t, err := template.ParseFiles(distDir); err != nil {
		log.Printf("Template parse files error, %v", err)
		w.WriteHeader(500)
		return
	} else {
		t.Execute(w, s)
	}
}

func NewRouter(ctx *appctx.Context) *mux.Router {
	router := mux.NewRouter()

	// middlewares
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(middleware.EnableCors)
	router.Use(middleware.Logging)

	// base
	router.HandleFunc("/", index).Methods(http.MethodGet)
	router.HandleFunc("/api/status", status).Methods(http.MethodGet)

	// loki proxy
	router.Path("/api/query").Handler(queryHandler(ctx)).Methods(http.MethodGet).Queries("query", "{query}")

	return router
}
