package main

import (
	"context"
	"event-automation/bot/fsm"
	"event-automation/bot/handlers"
	"event-automation/bot/sender"
	"event-automation/bot/storage"
	"event-automation/config"
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
	counter := 0 // временное решение для начала отсчета пересланных сообщений

	for update := range updates {
		userID := update.Message.From.ID
		language := update.Message.From.LanguageCode
		chatID := update.Message.Chat.ID
		userState := session.GetState(userID, language)

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
			if update.Message.Command() == "register" {
				sender.SendLocalizedMessage(chatID, userState.Language, "already registered")
				session.SetState(userID, "change_email")
				log.Printf("Смена состояния на изменения ящика")
				continue
			}
			if update.Message.ForwardFrom != nil || update.Message.ForwardSenderName != "" {
				handlers.CollectMessage(sender, store, update.Message)
				counter++
				if counter == 1 {
					go func(message *tgbotapi.Message) {
						ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*40))
						defer cancel()

						<-ctx.Done()

						sender.SendLocalizedMessage(chatID, userState.Language, "processing")
						err := handlers.CreateEvent(sender, store, message)
						if err != nil {
							log.Printf("Error: %v", err)
						}
						counter = 0
					}(update.Message)
				}
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
