package main

import (
	"event-automation/bot/fsm"
	"event-automation/bot/handlers"
	"event-automation/bot/processing"
	"event-automation/bot/sender"
	"event-automation/bot/storage"
	"event-automation/config"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	config.LoadEnv()

	bot, err := tgbotapi.NewBotAPI(config.GetEnv("TELEGRAM_API_TOKEN", ""))
	if err != nil {
		log.Fatalf("Error while creating bot: %v", err)
	}

	store := storage.NewStore()
	session := fsm.NewSession()
	sender := sender.NewSender(bot)

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	state := processing.NewProcessingState()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		language := update.Message.From.LanguageCode
		chatID := update.Message.Chat.ID
		userState := session.GetState(userID, language)

		switch userState.State {
		case "start":
			if store.IsExist(chatID) {
				sender.SendLocalizedMessage(chatID, userState.Language, "waiting")
				session.SetState(userID, "awaiting messages")
			} else {
				sender.SendLocalizedMessage(chatID, userState.Language, "welcome")
				sender.SendLocalizedMessage(chatID, userState.Language, "register required")
				session.SetState(userID, "initial")
			}

		case "initial":
			registered := handlers.Register(sender, store, update.Message)
			if registered {
				session.SetState(userID, "awaiting_messages")
			}

		case "awaiting_messages":
			if update.Message.ForwardFrom != nil || update.Message.ForwardSenderName != "" {
				handlers.CollectMessageAndSendEvent(state, sender, store, update.Message)
			} else {
				sender.SendLocalizedMessage(chatID, userState.Language, "waiting")
			}

		case "change_email":
			if update.Message.Command() == "yes" {
				sender.SendLocalizedMessage(chatID, userState.Language, "waiting email")
				session.SetState(userID, "awaiting_new_email")
			} else {
				sender.SendLocalizedMessage(chatID, userState.Language, "cancel email change")
				session.SetState(userID, "awaiting_messages")
			}

		case "awaiting_new_email":
			changeSessionState := handlers.ChangeEmail(sender, store, update.Message)
			if changeSessionState {
				session.SetState(userID, "awaiting_messages")
			}

		default:
			sender.SendLocalizedMessage(chatID, userState.Language, "error")
		}
	}
}
