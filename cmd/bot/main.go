package main

import (
	"log"
	"multitrack-bot/internal/adapters"
	"multitrack-bot/internal/bot"
	"multitrack-bot/internal/config"
	"multitrack-bot/internal/core"
)

func main() {
	// load config
	cfg := config.Load()

	if cfg.BotToken == "" {
		log.Fatal("BOT_TOKEN is required")
	}

	// initialize adapters manager
	adapterManager := adapters.NewAdapterManager(cfg)

	// initialize tracking service
	trackingService := core.NewTrackingService(adapterManager)

	// create and run tg bot
	bot, err := bot.NewBot(cfg.BotToken, trackingService)
	if err != nil {
		log.Fatal("Failed to create a bot:", err)
	}

	log.Println("Bot is starting...")
	bot.Start()
}
