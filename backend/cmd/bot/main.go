package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/akolesnov/football58/backend/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.LoadBot()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatal(err)
	}

	me, err := bot.GetMe()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("telegram bot started, username=%s, api_base_url=%s, admin_count=%d", me.UserName, cfg.APIBaseURL, len(cfg.TelegramAdminIDs))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Print("bot stopped")
}
