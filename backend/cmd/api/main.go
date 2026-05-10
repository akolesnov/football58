package main

import (
	"context"
	"log"
	"net/http"

	"github.com/akolesnov/football58/backend/internal/config"
	"github.com/akolesnov/football58/backend/internal/db"
	httpapi "github.com/akolesnov/football58/backend/internal/http"
	"github.com/akolesnov/football58/backend/internal/repository"
	"github.com/akolesnov/football58/backend/internal/service"
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
	venueRepository := repository.NewVenueRepository(postgres)
	venueService := service.NewVenueService(venueRepository)
	venueHandler := httpapi.NewVenueHandler(venueService)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler.Health)
	mux.HandleFunc("/health/db", healthHandler.Database)
	mux.HandleFunc("POST /venues", venueHandler.Create)
	mux.HandleFunc("GET /venues", venueHandler.List)
	mux.HandleFunc("GET /venues/{id}", venueHandler.GetByID)

	if err := http.ListenAndServe(cfg.HTTPAddr, mux); err != nil {
		log.Fatal(err)
	}
}
