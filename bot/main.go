package main

import (
	"context"
	"event-automation/bot/fsm"
	"event-automation/bot/handlers"
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

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	counter := 0 // временное решение для начала отсчета пересланных сообщений

	for update := range updates {
		if update.Message == nil {
			log.Printf("тута")
			continue
		}

		userID := update.Message.From.ID
		userState := session.GetState(userID)

		switch userState.State {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ку. Зарегистрируйтесь с помощью команды /register. Обязательна почта gmail")
			bot.Send(msg)
			session.SetState(userID, "initial")
		case "initial":
			if update.Message.Command() == "register" {
				registered := handlers.Register(bot, store, update.Message)
				if !registered {
					continue
				}
				session.SetState(userID, "awaiting_messages")
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сначала зарегистрируйтесь с помощью команды /register")
				bot.Send(msg)
			}

		case "awaiting_messages":
			if update.Message.Command() == "register" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы уже зарегистрированы."+
					"Хотите сменить email?	Введите /yes, если да ")
				bot.Send(msg)
				session.SetState(userID, "change_email")
				log.Printf("Смена состояния на изменения ящика")
				continue
			}
			if update.Message.ForwardFrom != nil || update.Message.ForwardSenderName != "" {
				handlers.CollectMessage(bot, store, update.Message)
				counter++
				if counter == 1 {
					go func(message *tgbotapi.Message) {
						ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*40))
						defer cancel()

						<-ctx.Done()

						msg := tgbotapi.NewMessage(message.Chat.ID, "Начинаю обработку сообщений")
						bot.Send(msg)
						err := handlers.CreateEvent(bot, store, message)
						if err != nil {
							log.Printf("Error: %v", err)
						}
						counter = 0
					}(update.Message)
				}
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Перешлите сообщение для обработки.")
				bot.Send(msg)
			}
		case "change_email":
			if update.Message.Command() == "yes" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите новый почтовый ящик.")
				bot.Send(msg)
				session.SetState(userID, "awaiting_new_email")
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отмена операции смены email")
			bot.Send(msg)
			session.SetState(userID, "awaiting_messages")
		case "awaiting_new_email":
			changeSessionState := handlers.ChangeEmail(bot, store, update.Message)
			if changeSessionState {
				log.Printf("залез не туда")
				session.SetState(userID, "awaiting_messages")
				continue
			}
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестное состояние.")
			bot.Send(msg)
		}
	}
}
