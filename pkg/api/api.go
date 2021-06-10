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
	QueryRange(string, string, string) (*loki.QueryResponse, error)
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

func queryRangeHandler(ctx *AppContext) http.Handler {
	handlerFunc := func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		query, ok := vars["query"]
		if !ok {
			log.Print("query is missing")
			rw.WriteHeader(400)
			return
		}
		start, ok := vars["start"]
		if !ok {
			log.Print("start is missing")
			rw.WriteHeader(400)
			return
		}
		end, ok := vars["end"]
		if !ok {
			log.Print("end is missing")
			rw.WriteHeader(400)
			return
		}

		resp, err := ctx.loki.QueryRange(query, start, end)
		if err != nil {
			log.Printf("Loky client error, %v", err)
			rw.WriteHeader(500)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(200)
		json.NewEncoder(rw).Encode(resp)
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

	router.Path("/api/query").Handler(queryHandler(ctx)).Methods(http.MethodGet).Queries("query", "{query}")
	router.Path("/api/query_range").Handler(queryRangeHandler(ctx)).Methods("GET").
		Queries("query", "{query}", "start", "{start:[0-9]+}", "end", "{end:[0-9]+}")

	return router
}
