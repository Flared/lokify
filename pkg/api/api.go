package api

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/flared/lokify/pkg/middleware"
	"github.com/gorilla/mux"
)

type context struct {
	client  *http.Client
	baseUrl string
}

func NewContext(client *http.Client, baseUrl string) *context {
	return &context{
		client:  client,
		baseUrl: baseUrl,
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func queryHandler(ctx *context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		query := vars["query"]
		limit := "100"

		params := strings.Join(
			[]string{
				"query=" + url.QueryEscape(query),
				"limit=" + url.QueryEscape(limit),
			},
			"&",
		)

		resp, err := ctx.client.Get(ctx.baseUrl + "/loki/api/v1/query?" + params)
		if err != nil {
			log.Printf("Loki Query error, %v", err)
			w.WriteHeader(500)
			return
		}

		defer resp.Body.Close()

		w.WriteHeader(resp.StatusCode)
		if _, err := ioutil.ReadAll(io.TeeReader(resp.Body, w)); err != nil {
			log.Printf("Loki Query error, %v", err)
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

func NewRouter(ctx *context) *mux.Router {
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
