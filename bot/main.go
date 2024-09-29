package main

import (
	"log"

	"event-automation/bot/clients/telegram"
	event_consumer "event-automation/bot/consumer/event-consumer"
	tgClient "event-automation/bot/events/telegram"
	storage "event-automation/bot/storage/map_storage"
	"event-automation/config"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

func main() {
	eventsProcessor := tgClient.New(telegram.New(tgBotHost, mustToken()), storage.New())
	log.Printf("starting telegram bot")
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatalf("failed to start consumer: %s", err)
	}
}

func mustToken() string {
	config.LoadEnv()
	token := config.GetEnv("TELEGRAM_API_TOKEN", "")
	if token == "" {
		log.Fatal("token is required")
	}

	return token
}
