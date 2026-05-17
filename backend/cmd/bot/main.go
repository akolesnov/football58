package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/akolesnov/football58/backend/internal/config"
)

func main() {
	cfg, err := config.LoadBot()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("bot config loaded, api_base_url=%s, admin_count=%d", cfg.APIBaseURL, len(cfg.TelegramAdminIDs))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Print("bot stopped")
}
