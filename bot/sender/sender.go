package sender

import (
	"event-automation/bot/messages"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendLocalizedMessage(bot *tgbotapi.BotAPI, chatID int64, lang string, key string) {
	text := messages.GetMessage(lang, key)
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}
