package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/akolesnov/football58/backend/internal/config"
	"github.com/akolesnov/football58/backend/internal/db"
)

func main() {
	cfg := config.Load()

	postgres, err := db.OpenPostgres(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer postgres.Close()

	http.HandleFunc("/", hello)

	if err := http.ListenAndServe(cfg.HTTPAddr, nil); err != nil {
		log.Fatal(err)
	}
}

func hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "hello world")
}
