package main

import (
	"log"
	"os"

	"github.com/flared/lokify/api"
)

type Config struct {
	LokiBaseUrl   string `json:"loki_base_url"`
	LokifyBaseUrl string `json:"lokify_base_url"`
	BuildDir      string `json:"build_dir"`
}

func main() {
	appConfigUrl := os.Getenv("APP_CONFIG_URL")
	if appConfigUrl == "" {
		log.Fatalf("lokify application needs a APP_CONFIG_URL env var")
	}

	appConfig, err := api.AppConfig(appConfigUrl)
	if err != nil {
		log.Fatalf("Load config error, %v", err)
	}

	api.RunServer(appConfig)
}
