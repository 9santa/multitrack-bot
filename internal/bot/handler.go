package bot

import (
	"context"
	"fmt"
	"multitrack-bot/internal/domain"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleStart(msg *tgbotapi.Message) {
	text := `👋 Привет! Я бот для отслеживания посылок.
/pochta - Отслеживание посылок Почты России.`

	b.sendMessage(msg.Chat.ID, text)
}

func (b *Bot) handleHelp(msg *tgbotapi.Message) {
	text := `📖 Справка:

• Отправь трек-номер для отслеживания посылки
• Поддерживается Почта России
• Формат трек-номера: 14 цифр (начинается с 0)

Команды:
/start - начать работу
/help - показать эту справку`

	b.sendMessage(msg.Chat.ID, text)
}

func (b *Bot) handleCarrierCommand(msg *tgbotapi.Message, carrier string) {
	var formatCarrier string
	switch carrier {
	case "russianpost":
		formatCarrier = "Почта России"
	}

	text := fmt.Sprintf("📬 Вы выбрали *%s*.\nПожалуйста, введите трек-номер.", formatCarrier)
	b.sendMessage(msg.Chat.ID, text)

	b.pendingCarrier[msg.Chat.ID] = carrier
}

func (b *Bot) handleTrackingNumber(msg *tgbotapi.Message, carrier string) {
	trackingNumber := strings.TrimSpace(msg.Text)

	// number validation
	if len(trackingNumber) < 8 {
		b.sendMessage(msg.Chat.ID, "❌ Неверный формат трек-номера. Должен содержать не менее 8 символов.")
		return
	}

	delete(b.pendingCarrier, msg.Chat.ID)

	b.trackPackage(msg.Chat.ID, trackingNumber, carrier)

}

func (b *Bot) trackPackage(chatID int64, trackingNumber string, carrier string) {
	processingMsg := b.sendMessage(chatID, "🔄 Отслеживаю посылку...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := b.trackingService.Track(ctx, trackingNumber, carrier)
	if err != nil {
		b.editMessage(processingMsg.Chat.ID, processingMsg.MessageID,
			"❌ Не удалось отследить посылку. Проверьте номер и попробуйте позже.")
		return
	}

	response := b.formatTrackingResponse(result)
	b.sendMessage(processingMsg.Chat.ID, response)
	// b.editMessage(processingMsg.Chat.ID, processingMsg.MessageID, response)
}

func (b *Bot) formatTrackingResponse(result *domain.TrackingResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("📦 *Посылка %s*\n", result.Number))
	sb.WriteString(fmt.Sprintf("🚚 *Курьер:* %s\n", result.Courier))
	sb.WriteString(fmt.Sprintf("📊 *Статус:* %s\n", result.Status))
	sb.WriteString(fmt.Sprintf("📝 *Описание:* %s\n\n", result.Description))

	if len(result.Checkpoints) > 0 {
		sb.WriteString("*Последние события:*\n")
		for i, checkpoint := range result.Checkpoints {
			if i >= 3 { // display only the latest 3 tracking updates
				break
			}
			sb.WriteString(fmt.Sprintf("• %s - %s\n",
				checkpoint.Date.Format("02.01.2006 15:04"),
				checkpoint.Description))
		}
	}
	return sb.String()
}

func (b *Bot) sendMessage(chatID int64, text string) tgbotapi.Message {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"

	message, _ := b.api.Send(msg)
	return message
}

func (b *Bot) editMessage(chatID int64, messageID int, text string) {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ParseMode = "Markdown"

	b.api.Send(msg)
}
