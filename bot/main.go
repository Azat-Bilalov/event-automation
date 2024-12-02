package main

import (
	"event-automation/bot/fsm"
	"event-automation/bot/handlers"
	"event-automation/bot/sender"
	"event-automation/bot/storage"
	"event-automation/config"
	"fmt"
	"log"
	"time"

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

	// Буфер для сообщений и флаги обработки
	isProcessing := make(map[int64]bool) // Флаг обработки сообщений для пользователя

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
				session.SetState(userID, "awaiting_messages")
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
				handlers.CollectMessage(sender, store, update.Message)
				if !isProcessing[userID] {
					isProcessing[userID] = true

					go func(message *tgbotapi.Message) {
						time.Sleep(2 * time.Second)

						// Обрабатываем накопленные сообщения
						if len(handlers.UserMessages[userID]) > 0 {
							fmt.Printf("Processing messages for user %d: %v\n", userID, handlers.UserMessages[userID])
							sender.SendLocalizedMessage(userID, userState.Language, "processing")
							handlers.CreateEvent(sender, store, message)
							if err != nil {
								log.Printf("Error processing messages for user %d: %v", userID, err)
							}
							handlers.UserMessages[userID] = nil
						}
						isProcessing[userID] = false
					}(update.Message)
				}
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
