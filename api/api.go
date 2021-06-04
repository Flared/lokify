package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/flared/lokify/pkg/loki"
	"github.com/gorilla/mux"
)

type Config struct {
	LokiBaseUrl   string `json:"loki_base_url"`
	LokifyBaseUrl string `json:"lokify_base_url"`
	BuildDir      string `json:"build_dir"`
}

type app struct {
	config     *Config
	lokiClient *loki.LokiClient
}

func AppConfig(appConfigUrl string) (*Config, error) {
	jsonFile, err := os.Open(appConfigUrl)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	var appConfig Config
	json.Unmarshal(bytes, &appConfig)

	return &appConfig, nil
}

func (a app) test(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "test")
}

func (a app) status(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "ok")
}

func (a app) query(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	query, ok := vars["query"]
	if !ok {
		query = "{container_name=\"firework-api\"} | json"
	}

	resp, errQuery := a.lokiClient.Query(query)
	if errQuery != nil {
		log.Printf("Loky client error, %v", errQuery)
		rw.WriteHeader(500)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(resp)
}

func (a app) queryRange(rw http.ResponseWriter, r *http.Request) {
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

	resp, errQuery := a.lokiClient.QueryRange(query, start, end)
	if errQuery != nil {
		log.Printf("Loky client error, %v", errQuery)
		rw.WriteHeader(500)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(resp)
}

func (a app) labels(rw http.ResponseWriter, r *http.Request) {
	resp, errQuery := a.lokiClient.Labels()
	if errQuery != nil {
		log.Printf("Loky client error, %v", errQuery)
		rw.WriteHeader(500)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(resp)
}

func (a app) labelValues(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	label, ok := vars["label"]
	if !ok {
		log.Print("label is missing")
		rw.WriteHeader(400)
		return
	}

	resp, errQuery := a.lokiClient.LabelValues(label)
	if errQuery != nil {
		log.Printf("Loky client error, %v", errQuery)
		rw.WriteHeader(500)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(resp)
}

func (a app) labelsValues(rw http.ResponseWriter, r *http.Request) {
	labelsResp, errQuery := a.lokiClient.Labels()
	if errQuery != nil {
		log.Printf("Loky client error, %v", errQuery)
		rw.WriteHeader(500)
		return
	}

	labelsValues := make(map[string][]string)
	for _, label := range labelsResp.Data {
		labelValuesResp, errLabelValues := a.lokiClient.LabelValues(label)
		if errLabelValues != nil {
			log.Printf("Loky client error, %v", errLabelValues)
			rw.WriteHeader(500)
			return
		}

		labelsValues[label] = labelValuesResp.Data
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(map[string]map[string][]string{
		"label": labelsValues,
	})
}

func (a app) index(rw http.ResponseWriter, r *http.Request) {
	pwd, _ := os.Getwd()
	t, err := template.ParseFiles(pwd + a.config.BuildDir)
	if err != nil {
		log.Printf("Template parse files error, %v", err)
		rw.WriteHeader(500)
		return
	}

	t.Execute(rw, a.config)
}

func enableCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(rw, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(rw, r)
	})
}

func RunServer(config *Config) {
	router := mux.NewRouter()
	a := app{
		config:     config,
		lokiClient: loki.NewClient(config.LokiBaseUrl),
	}

	router.Use(enableCorsMiddleware)
	router.Use(loggingMiddleware)

	router.HandleFunc("/", a.index).Methods("GET")
	router.HandleFunc("/api/status", a.status).Methods("GET")
	router.HandleFunc("/api/test", a.test).Methods("GET")

	router.HandleFunc("/api/query", a.query).Methods("GET")
	router.HandleFunc("/api/query", a.query).Methods("GET").Queries("query", "{query}")
	router.HandleFunc("/api/query_range", a.queryRange).Methods("GET").
		Queries("query", "{query}").
		Queries("start", "{start:[0-9]+}").
		Queries("end", "{end:[0-9]+}")

	router.HandleFunc("/api/labels", a.labels).Methods("GET")
	router.HandleFunc("/api/labels/values", a.labelsValues).Methods("GET")
	router.HandleFunc("/api/labels/{label}", a.labelValues).Methods("GET")

	fmt.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
