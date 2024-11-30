package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"event-automation/bot/sender"
	"event-automation/bot/storage"
	validate "event-automation/lib/email"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageData struct {
	Messages []string `json:"messages"`
	Language string   `json:"language"`
	Timezone int64    `json:"timezone"`
}

type LlmResponse struct {
	Title         string `json:"title"`
	StartDatetime string `json:"start_datetime"`
	EndDatetime   string `json:"end_datetime"`
}

type EventData struct {
	Title         string   `json:"title"`
	Desc          string   `json:"description"`
	Emails        []string `json:"recipients_emails"`
	StartDatetime string   `json:"start_datetime"`
	EndDatetime   string   `json:"end_datetime"`
	Timezone      int64    `json:"timezone"`
}

var userMessages = make(map[int64][]string)
var emails []string

func Register(sender *sender.Sender, store storage.Storage, message *tgbotapi.Message) bool {
	email := message.Text
	if !validate.IsEmail(email) {
		sender.SendLocalizedMessage(message.From.ID, message.From.LanguageCode, "check email")
		return false
	}
	store.SetEmail(message.From.ID, email)
	sender.SendLocalizedMessage(message.From.ID, message.From.LanguageCode, "successful registration")
	return true
}

func ChangeEmail(sender *sender.Sender, store storage.Storage, message *tgbotapi.Message) (changeSessionState bool) {
	changeSessionState = false
	if message.Command() == "/return" {
		changeSessionState = true
		return changeSessionState
	}
	email := message.Text
	if !validate.IsEmail(email) {
		sender.SendLocalizedMessage(message.From.ID, message.From.LanguageCode, "check email")
		return changeSessionState
	}
	store.SetEmail(message.From.ID, email)
	sender.SendLocalizedMessage(message.From.ID, message.From.LanguageCode, "successful email change")
	changeSessionState = true
	return changeSessionState
}

func addEmailReceiver(_ *sender.Sender, store storage.Storage, message *tgbotapi.Message) {
	if email := store.GetEmail(message.ForwardFrom.ID); email != "" {
		emails = append(emails, store.GetEmail(message.ForwardFrom.ID))
	} else {
		// TODO: вот тут надо прикинуть как прокидывать имя юзера внутрь функции сендер
		// msg := tgbotapi.NewMessage(message.Chat.ID, "Пользователь %s не зарегистрирован в боте, невозможно создать событие для него")
		// bot.Send(msg)
	}
}

func parseMessage(sender *sender.Sender, store storage.Storage, message *tgbotapi.Message) string {
	var name string
	if message.ForwardFrom == nil {
		name = message.ForwardSenderName
	} else {
		name = message.ForwardFrom.FirstName + " " + message.From.LastName
		addEmailReceiver(sender, store, message)
	}

	text := fmt.Sprintf("%s: %s", name, strings.TrimSpace(message.Text))

	return text
}

func CollectMessage(sender *sender.Sender, store storage.Storage, message *tgbotapi.Message) {
	text := parseMessage(sender, store, message)
	userMessages[message.From.ID] = append(userMessages[message.From.ID], text)
}

func CreateEvent(sender *sender.Sender, store storage.Storage, message *tgbotapi.Message) error {
	userID := message.From.ID
	data := MessageData{
		Messages: userMessages[userID],
		Language: message.From.LanguageCode,
		Timezone: 3,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:8080/new_meet", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error from LLM service: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response LlmResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	log.Println("Response Body:", string(body))

	reqBody := EventData{
		Title:         response.Title,
		Desc:          response.Title,
		Emails:        emails,
		StartDatetime: response.StartDatetime,
		EndDatetime:   response.EndDatetime,
		Timezone:      3,
	}

	payload, err = json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err = http.Post("http://localhost:8080/create_event", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	log.Println(resp)

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(body)

	delete(userMessages, userID)

	sender.SendLocalizedMessage(message.From.ID, message.From.LanguageCode, "success")
	//TODO: Рассылка сообщений всем участникам переписки (для скрытых)

	return nil
}
