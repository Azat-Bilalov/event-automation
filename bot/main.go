package main

import (
	"log"

	"event-automation/bot/clients/telegram"
	"event-automation/config"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {

	tgClient := telegram.New(tgBotHost, mustToken())

}

func mustToken() string {
	token := config.GetEnv("TELEGRAM_API_TOKEN", "")
	if token == "" {
		log.Fatal("token is required")
	}

	return token
}
