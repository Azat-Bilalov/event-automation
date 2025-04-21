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

		text := update.Message.Text

		switch text {
		case "ℹ️ О боте":
			sender.SendLocalizedMessage(chatID, userState.Language, "welcome")
			sender.SendMenu(chatID, userState.State)
			continue

		case "📝 Зарегистрироваться":
			session.SetState(userID, "initial")
			sender.SendLocalizedMessage(chatID, userState.Language, "register required")
			sender.SendMenu(chatID, "initial")
			continue

		case "📧 Сменить Email":
			session.SetState(userID, "change_email")
			sender.SendMenu(chatID, "change_email")
			continue

		case "🔙 Главное меню":
			session.SetState(userID, "awaiting_messages")
			sender.SendLocalizedMessage(chatID, userState.Language, "back to menu")
			sender.SendMenu(chatID, "awaiting_messages")
			continue
		}

		switch userState.State {
		case "start":
			sender.SendMenu(chatID, userState.State)
			if store.IsExist(chatID) {
				sender.SendLocalizedMessage(chatID, userState.Language, "waiting")
				session.SetState(userID, "awaiting messages")
			} else {
				sender.SendLocalizedMessage(chatID, userState.Language, "welcome")
				sender.SendLocalizedMessage(chatID, userState.Language, "register required")
			}

		case "initial":
			registered := handlers.Register(sender, store, update.Message)
			if registered {
				session.SetState(userID, "awaiting_messages")
				sender.SendMenu(chatID, "awaiting_messages")
			}

		case "awaiting_messages":
			if update.Message.ForwardFrom != nil || update.Message.ForwardSenderName != "" {
				handlers.CollectMessageAndSendEvent(state, sender, store, update.Message)
			} else {
				sender.SendLocalizedMessage(chatID, userState.Language, "waiting")
			}

		case "change_email":
			sender.SendLocalizedMessage(chatID, userState.Language, "waiting email")
			session.SetState(userID, "awaiting_new_email")

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
