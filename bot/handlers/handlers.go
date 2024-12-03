package handlers

import (
	"event-automation/bot/sender"
	"event-automation/bot/storage"
	validate "event-automation/lib/email"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageData struct {
	Messages []string `json:"messages"`
	Language string   `json:"language"`
	Timezone int64    `json:"timezone"`
}

type LlmResponse struct {
	Title         string `json:"title"`
	StartDatetime string `json:"start_datetime"`
	EndDatetime   string `json:"end_datetime"`
}

type CalendarResponse struct {
	EventLink string `json:"event_link"`
}

type EventData struct {
	Title         string   `json:"title"`
	Desc          string   `json:"description"`
	Emails        []string `json:"recipients_emails"`
	StartDatetime string   `json:"start_datetime"`
	EndDatetime   string   `json:"end_datetime"`
	Timezone      int64    `json:"timezone"`
}

var emails []string

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
