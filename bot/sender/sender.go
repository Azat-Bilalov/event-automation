package sender

import (
	"event-automation/bot/messages"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender struct {
	bot *tgbotapi.BotAPI
}

func NewSender(bot *tgbotapi.BotAPI) *Sender {
	return &Sender{
		bot: bot,
	}
}

// TODO: придумать способ типизировать аргументы
func (s *Sender) SendLocalizedMessage(chatID int64, lang string, key string, args ...interface{}) {
	text := messages.GetMessage(lang, key, args...)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	s.bot.Send(msg)
}
