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
	cfg, err := config.LoadAPI()
	if err != nil {
		log.Fatal(err)
	}

	postgres, err := db.OpenPostgres(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer postgres.Close()

	healthHandler := httpapi.NewHealthHandler(postgres)
	gameRepository := repository.NewGameRepository(postgres)
	gameMemberRepository := repository.NewGameMemberRepository(postgres)
	gameService := service.NewGameService(gameRepository, gameMemberRepository)
	gameHandler := httpapi.NewGameHandler(gameService)
	gameMemberService := service.NewGameMemberService(postgres)
	gameMemberHandler := httpapi.NewGameMemberHandler(gameMemberService)
	userRepository := repository.NewUserRepository(postgres)
	userService := service.NewUserService(userRepository)
	userHandler := httpapi.NewUserHandler(userService)
	venueRepository := repository.NewVenueRepository(postgres)
	venueService := service.NewVenueService(venueRepository)
	venueHandler := httpapi.NewVenueHandler(venueService)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler.Health)
	mux.HandleFunc("/health/db", healthHandler.Database)
	mux.HandleFunc("POST /games", gameHandler.Create)
	mux.HandleFunc("GET /games", gameHandler.List)
	mux.HandleFunc("GET /games/{id}", gameHandler.GetByID)
	mux.HandleFunc("PATCH /games/{id}/telegram-message", gameHandler.UpdateTelegramMessage)
	mux.HandleFunc("POST /games/{id}/close", gameHandler.Close)
	mux.HandleFunc("POST /games/{id}/cancel", gameHandler.Cancel)
	mux.HandleFunc("POST /games/{id}/finish", gameHandler.Finish)
	mux.HandleFunc("POST /games/{id}/members", gameMemberHandler.Join)
	mux.HandleFunc("POST /games/{id}/members/{member_id}/cancel", gameMemberHandler.Cancel)
	mux.HandleFunc("POST /games/{id}/me/cancel", gameMemberHandler.CancelOwnTelegramMembership)
	mux.HandleFunc("POST /users/telegram/upsert", userHandler.UpsertTelegram)
	mux.HandleFunc("GET /users/{id}", userHandler.GetByID)
	mux.HandleFunc("GET /users/by-telegram/{telegram_id}", userHandler.GetByTelegramID)
	mux.HandleFunc("POST /venues", venueHandler.Create)
	mux.HandleFunc("GET /venues", venueHandler.List)
	mux.HandleFunc("GET /venues/{id}", venueHandler.GetByID)

	if err := http.ListenAndServe(cfg.HTTPAddr, mux); err != nil {
		log.Fatal(err)
	}
}
