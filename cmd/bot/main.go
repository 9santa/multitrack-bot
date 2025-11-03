package main

import (
	"log"
	"multitrack-bot/internal/adapters"
	"multitrack-bot/internal/bot"
	"multitrack-bot/internal/config"
	"multitrack-bot/internal/core"
	"multitrack-bot/internal/gateway"
	"os"
	"os/signal"
	"syscall"
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

	// create gateway server
	gatewayServer := gateway.NewGateway(trackingService)

	// run bot and gateway in goroutines
	go func() {
		log.Println("Bot is starting...")
		bot.Start()
	}()

	go func() {
		log.Printf("Gateway server is starting on port %s...", cfg.ServerPort)
		if err := gatewayServer.Run(cfg.ServerPort); err != nil {
			log.Fatalf("Failed to start gateway server: %v", err)
		}
	}()

	// wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")

}
