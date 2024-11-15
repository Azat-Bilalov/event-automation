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

		log.Printf("fff[%s] %s", update.Message.From.UserName, update.Message.Text)
		log.Printf("asd[%s] %s", update.Message)

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
			if counter == 0 {
				go func() {
					ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*40))
					defer cancel()
					select {
					case <-ctx.Done():
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Let me think...")
						bot.Send(msg)
						handlers.CreateEvent(bot, store, update.Message)
						counter = 0
					}
				}()
			}
			counter++
			handlers.CollectMessage(update.Message, store)
		}
	}
}
