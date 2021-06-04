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

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "test")
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func query(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	query, ok := vars["query"]
	if !ok {
		query = "{container_name=\"firework-api\"} | json"
	}

	client := loki.NewClient("http://localhost:3100")
	resp, errQuery := client.Query(query)
	if errQuery != nil {
		log.Printf("Loky client error, %v", errQuery)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func index(w http.ResponseWriter, r *http.Request) {
	s := make(map[string]string)
	s["base_url"] = "http://localhost:8080"

	distDir := "../ui/build/index.html"
	t, err := template.ParseFiles(distDir)
	if err != nil {
		log.Printf("Template parse files error, %v", err)
		w.WriteHeader(500)
		return
	}

	t.Execute(w, s)
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

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(enableCorsMiddleware)
	router.Use(loggingMiddleware)

	router.HandleFunc("/", index).Methods(http.MethodGet)
	router.HandleFunc("/api/status", status).Methods(http.MethodGet)
	router.HandleFunc("/api/test", test).Methods(http.MethodGet)

	router.HandleFunc("/api/query", query).Methods(http.MethodGet)
	router.HandleFunc("/api/query", query).Methods(http.MethodGet).Queries("query", "{query}")

	return router
}
