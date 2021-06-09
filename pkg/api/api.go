package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/flared/lokify/pkg/loki"
	"github.com/gorilla/mux"
)

type AppContext struct {
	loki lokiClient
}

func NewAppContext(loki lokiClient) *AppContext {
	return &AppContext{
		loki: loki,
	}
}

type lokiClient interface {
	Query(string) (*loki.QueryResponse, error)
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "test")
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func queryHandler(appCtx *AppContext) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		query, ok := vars["query"]
		if !ok {
			query = "{container_name=\"firework-api\"} | json"
		}

		resp, err := appCtx.loki.Query(query)
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

		w.Header().Set("Content-Type", "application/json")
	}

	return http.HandlerFunc(handlerFunc)
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

func enableCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func NewRouter(ctx *AppContext) *mux.Router {
	router := mux.NewRouter()

	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(enableCorsMiddleware)
	router.Use(loggingMiddleware)

	router.HandleFunc("/", index).Methods(http.MethodGet)
	router.HandleFunc("/api/status", status).Methods(http.MethodGet)
	router.HandleFunc("/api/test", test).Methods(http.MethodGet)

	router.Path("/api/query").Handler(queryHandler(ctx)).Methods(http.MethodGet)
	router.Path("/api/query").Handler(queryHandler(ctx)).Methods(http.MethodGet).Queries("query", "{query}")

	return router
}
