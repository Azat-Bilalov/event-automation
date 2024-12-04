package handlers

import (
	"event-automation/bot/models"
	"event-automation/bot/processing"
	"event-automation/bot/sender"
	"event-automation/bot/storage"
	"event-automation/bot/utils"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const processDelay = 2 * time.Second

func parseMessage(message *tgbotapi.Message) string {
	var name string
	if name = utils.ExtractName(message.ForwardFrom); name == "Unknown" {
		name = message.ForwardSenderName
	}
	text := fmt.Sprintf("%s: %s", name, strings.TrimSpace(message.Text))

	return text
}

func handleForwardedMessage(state *processing.ProcessingState, store storage.Storage, message *tgbotapi.Message) {
	userID := message.From.ID
	text := parseMessage(message)
	state.ForwardedMessages[userID] = append(state.ForwardedMessages[userID], text)

	if message.ForwardFrom == nil {
		state.AddToClosedAccount(userID, message.ForwardSenderName)
	} else if email := store.GetEmail(message.ForwardFrom.ID); email != "" {
		state.Emails[userID] = append(state.Emails[userID], email)
	} else {
		state.AddToNotInDB(userID, utils.ExtractName(message.ForwardFrom))
	}
}

func processMessagesForUser(
	state *processing.ProcessingState,
	sender *sender.Sender,
	message tgbotapi.Message,
) {
	userID := message.From.ID

	if state.IsProcessing[userID] {
		return
	}
	state.IsProcessing[userID] = true

	go func(msg tgbotapi.Message) {
		userID := msg.From.ID

		defer func() { state.IsProcessing[userID] = false }()
		time.Sleep(processDelay)

		language := msg.From.LanguageCode
		sender.SendLocalizedMessage(userID, language, "processing")

		event := &models.UserEvent{
			UserID:   userID,
			Language: language,
			Messages: state.ForwardedMessages[userID],
			Emails:   state.Emails[userID],
			Timezone: 3,
		}

		calendarResponse, err := handleEventCreation(event)
		if err != nil {
			log.Printf("Error processing messages for user %d: %v", userID, err)
		}

		sender.SendLocalizedMessage(
			event.UserID,
			event.Language,
			"success",
			calendarResponse.EventLink,
			state.InaccessibleClosed[userID],
			state.InaccessibleNotInDB[userID],
		)

		state.ClearUserData(userID)
	}(message)
}

func CollectMessageAndSendEvent(
	state *processing.ProcessingState,
	sender *sender.Sender,
	store storage.Storage,
	message *tgbotapi.Message,
) {
	handleForwardedMessage(state, store, message)
	processMessagesForUser(state, sender, *message)
}
