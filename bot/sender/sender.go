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

func (s *Sender) SendLocalizedMessage(chatID int64, lang string, key string) {
	text := messages.GetMessage(lang, key)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	s.bot.Send(msg)
}
