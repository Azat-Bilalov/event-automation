package handlers

import (
	"event-automation/bot/sender"
	"event-automation/bot/storage"
	validate "event-automation/lib/email"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Register(sender *sender.Sender, store storage.Storage, message *tgbotapi.Message) bool {
	email := message.Text
	if !validate.IsEmail(email) {
		sender.SendLocalizedMessage(message.From.ID, message.From.LanguageCode, "check email")
		return false
	}
	store.SetEmail(message.From.ID, email)
	sender.SendLocalizedMessage(message.From.ID, message.From.LanguageCode, "successful registration")
	return true
}

func ChangeEmail(sender *sender.Sender, store storage.Storage, message *tgbotapi.Message) (changeSessionState bool) {
	changeSessionState = false
	if message.Command() == "/return" {
		changeSessionState = true
		return changeSessionState
	}
	email := message.Text
	if !validate.IsEmail(email) {
		sender.SendLocalizedMessage(message.From.ID, message.From.LanguageCode, "check email")
		return changeSessionState
	}
	store.SetEmail(message.From.ID, email)
	sender.SendLocalizedMessage(message.From.ID, message.From.LanguageCode, "successful email change")
	changeSessionState = true
	return changeSessionState
}
