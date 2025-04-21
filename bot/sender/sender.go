package sender

import (
	"event-automation/bot/messages"
	"log"

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

func getKeyboardForState(state string) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton

	switch state {
	case "start":
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("ℹ️ Обо мне"),
			tgbotapi.NewKeyboardButton("📝 Зарегистрироваться"),
		))

	case "awaiting_messages":
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📧 Сменить Email"),
			tgbotapi.NewKeyboardButton("ℹ️ Обо мне"),
		))

	case "change_email", "awaiting_new_email":
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔙 Главное меню"),
		))

	default:
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("ℹ️ Обо мне"),
			tgbotapi.NewKeyboardButton("🔙 Главное меню"),
		))
	}

	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard:        rows,
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}

func (s *Sender) SendMenu(chatID int64, state string) {
	msg := tgbotapi.NewMessage(chatID, "Выбери действие: ") // Используем неразрывный пробел (U+200E)
	msg.ReplyMarkup = getKeyboardForState(state)

	// Отправляем сообщение
	if _, err := s.bot.Send(msg); err != nil {
		log.Printf("Error sending menu: %v", err)
	}
}
