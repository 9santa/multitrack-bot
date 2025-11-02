package bot

import (
	"fmt"
	"log"
	"multitrack-bot/internal/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api             *tgbotapi.BotAPI
	trackingService *core.TrackingService
	pendingCarrier  map[int64]string // chatID -> carrier name
}

func NewBot(token string, trackingService *core.TrackingService) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create a bot: %w", err)
	}

	bot := &Bot{
		api:             api,
		trackingService: trackingService,
		pendingCarrier:  make(map[int64]string),
	}

	return bot, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	log.Printf("Bot started and listening for messages...")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go b.handleMessage(update.Message)
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	if msg.IsCommand() {
		switch msg.Command() {
		case "start":
			b.handleStart(msg)
		case "help":
			b.handleHelp(msg)
		case "pochta":
			b.handleCarrierCommand(msg, "russianpost")
		default:
			b.sendMessage(msg.Chat.ID, "Unknown Command. Use /help to see what's available.")
		}
		return
	}

	if carrier, ok := b.pendingCarrier[msg.Chat.ID]; ok {
		b.handleTrackingNumber(msg, carrier)
	}
}
