package main

import (
	"event-automation/config"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func mainold() {
	config.LoadEnv()

	bot, err := tgbotapi.NewBotAPI(config.GetEnv("TELEGRAM_API_TOKEN", ""))
	if err != nil {
		log.Fatalf("Error while creating bot: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Обработка команд
		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to the bot!")
			bot.Send(msg)
		case "create_event":
			// Вызов Google Calendar микросервиса
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Creating event...")
			bot.Send(msg)
		case "ask_ai":
			// Вызов LLM микросервиса для обработки
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Let me think...")
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
			bot.Send(msg)
		}
	}
}
