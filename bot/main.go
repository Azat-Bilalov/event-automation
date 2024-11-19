package main

import (
	"context"
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

		log.Printf("fff %s  %s", update.Message.From.UserName, update.Message.Text)
		log.Printf("asd %v ", update.Message.ForwardFrom)
		log.Printf("adsd %v ", update.Message.ForwardSenderName)
		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to the bot!")
			bot.Send(msg)
		case "help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "There are instructions!")
			bot.Send(msg)
		case "register":
			handlers.Register(bot, store, update.Message)
		default:
			if update.Message.ForwardFrom != nil || update.Message.ForwardSenderName != "" {
				if !store.IsExist(update.Message.From.ID) {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для доступа к этому функционалу необходимо зарегистрироваться. Используйте команду /register <Ваша@почта>")
					bot.Send(msg)
					continue
				}
				if counter == 0 {
					go func() {
						ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*40))
						defer cancel()

						<-ctx.Done()

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Начинаю обработку сообщений")
						bot.Send(msg)
						handlers.CreateEvent(bot, store, update.Message)
						counter = 0

					}()
				}
				log.Printf("пришел")
				counter++
				handlers.CollectMessage(bot, store, update.Message)
			}
		}
	}
}
