package handlers

import (
	"event-automation/bot/sender"
	"event-automation/bot/storage"
	"event-automation/bot/utils"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func parseMessage(message *tgbotapi.Message) string {
	var name string
	if name = utils.ExtractName(message.ForwardFrom); name == "Unknown" {
		name = message.ForwardSenderName
	}
	text := fmt.Sprintf("%s: %s", name, strings.TrimSpace(message.Text))

	return text
}

type UserProcessingState struct {
	// мапа с массивом пересланных сообщений от каждого пользователя
	ForwardedMessages map[int64][]string
	// мапа с флагом начала обработки пересланных сообщений для пользователя
	IsProcessing map[int64]bool
	// список имен пользователей, которые не получат событие - закрытый акк
	InaccessibleClosed map[int64][]string
	// список имен пользователей, которые не получат событие - нет в базе
	InaccessibleNotInDB map[int64][]string
}

func NewUserProcessingState() *UserProcessingState {
	return &UserProcessingState{
		ForwardedMessages:   make(map[int64][]string),
		IsProcessing:        make(map[int64]bool),
		InaccessibleClosed:  make(map[int64][]string),
		InaccessibleNotInDB: make(map[int64][]string),
	}
}

func handleClosedAccount(state *UserProcessingState, userID int64, name string) {
	state.InaccessibleClosed[userID] = append(state.InaccessibleClosed[userID], name)
}

func handleNotInDB(state *UserProcessingState, userID int64, name string) {
	state.InaccessibleNotInDB[userID] = append(state.InaccessibleNotInDB[userID], name)
}

func handleForwardedMessage(state *UserProcessingState, store storage.Storage, message *tgbotapi.Message) {
	userID := message.From.ID
	text := parseMessage(message)
	state.ForwardedMessages[userID] = append(state.ForwardedMessages[userID], text)

	if message.ForwardFrom == nil {
		handleClosedAccount(state, userID, message.ForwardSenderName)
	} else if email := store.GetEmail(message.ForwardFrom.ID); email != "" {
		emails = append(emails, email)
	} else {
		handleNotInDB(state, userID, utils.ExtractName(message.ForwardFrom))
	}
}

func processMessagesForUser(
	state *UserProcessingState,
	sender *sender.Sender,
	message *tgbotapi.Message,
	delay time.Duration,
) {
	userID := message.From.ID

	// Если обработка уже запущена, выходим
	if state.IsProcessing[userID] {
		return
	}
	state.IsProcessing[userID] = true

	// Запускаем горутину для обработки
	go func() {
		defer func() { state.IsProcessing[userID] = false }()
		time.Sleep(delay)

		language := message.From.LanguageCode
		sender.SendLocalizedMessage(userID, language, "processing")

		ctx := &UserEventContext{
			UserID:   userID,
			Language: language,
			Messages: state.ForwardedMessages[userID],
			Emails:   emails,
			Timezone: 3,
		}

		calendarResponse, err := createEvent(ctx)
		if err != nil {
			log.Printf("Error processing messages for user %d: %v", userID, err)
		}

		// Отправка подтверждения пользователю
		sender.SendLocalizedMessage(
			ctx.UserID,
			ctx.Language,
			"success",
			calendarResponse.EventLink,
			state.InaccessibleClosed[userID],
			state.InaccessibleNotInDB[userID],
		)

		state.ForwardedMessages[userID] = nil
		state.InaccessibleClosed[userID] = nil
		state.InaccessibleNotInDB[userID] = nil
	}()
}

func CollectMessageAndSendEvent(
	state *UserProcessingState,
	sender *sender.Sender,
	store storage.Storage,
	message *tgbotapi.Message,
) {
	handleForwardedMessage(state, store, message)
	processMessagesForUser(state, sender, message, 2*time.Second)
}
