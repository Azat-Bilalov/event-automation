package main

import (
	"event-automation/bot/fsm"
	"event-automation/bot/handlers"
	"event-automation/bot/sender"
	"event-automation/bot/storage"
	"event-automation/config"
	"fmt"
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

	// TODO: держать вместе userMessages и lastUpdateID
	lastUpdateID := make(map[int64]int)

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		userID := update.Message.From.ID
		language := update.Message.From.LanguageCode
		chatID := update.Message.Chat.ID
		userState := session.GetState(userID, language)

		if len(handlers.UserMessages[userID]) > 0 && update.UpdateID != lastUpdateID[userID] {
			fmt.Printf("Processing messages for user %d: %v\n", userID, handlers.UserMessages[userID])
			sender.SendLocalizedMessage(userID, userState.Language, "processing")
			// TODO: не передавать экземпляр сообщения, а только нужные свойства
			err := handlers.CreateEvent(sender, store, update.Message)
			if err != nil {
				log.Printf("Error: %v", err)
			}
		}

		switch userState.State {
		case "start":
			if store.IsExist(chatID) {
				// TODO: В будущем в этом блоке будем добавлять обработку старта с параметрами
				sender.SendLocalizedMessage(chatID, userState.Language, "waiting")
				session.SetState(userID, "awaiting_messages")
			} else {
				sender.SendLocalizedMessage(chatID, userState.Language, "welcome")
				sender.SendLocalizedMessage(chatID, userState.Language, "register required")
				session.SetState(userID, "initial")
			}
		case "initial":
			registered := handlers.Register(sender, store, update.Message)
			if !registered {
				continue
			}
			session.SetState(userID, "awaiting_messages")

		case "awaiting_messages":
			if update.Message.ForwardFrom != nil || update.Message.ForwardSenderName != "" {
				handlers.CollectMessage(sender, store, update.Message)
			} else {
				sender.SendLocalizedMessage(chatID, userState.Language, "waiting")
			}
		case "change_email":
			if update.Message.Command() == "yes" {
				sender.SendLocalizedMessage(chatID, userState.Language, "waiting email")
				session.SetState(userID, "awaiting_new_email")
				continue
			}
			sender.SendLocalizedMessage(chatID, userState.Language, "cancel email change")
			session.SetState(userID, "awaiting_messages")
		case "awaiting_new_email":
			changeSessionState := handlers.ChangeEmail(sender, store, update.Message)
			if changeSessionState {
				session.SetState(userID, "awaiting_messages")
				continue
			}
		default:
			sender.SendLocalizedMessage(chatID, userState.Language, "error")
		}
	}
}
