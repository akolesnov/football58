package main

import (
	"context"
	"log"
	"net/http"

	"github.com/akolesnov/football58/backend/internal/config"
	"github.com/akolesnov/football58/backend/internal/db"
	httpapi "github.com/akolesnov/football58/backend/internal/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	postgres, err := db.OpenPostgres(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer postgres.Close()

	healthHandler := httpapi.NewHealthHandler(postgres)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler.Health)
	mux.HandleFunc("/health/db", healthHandler.Database)

	if err := http.ListenAndServe(cfg.HTTPAddr, mux); err != nil {
		log.Fatal(err)
	}
}
