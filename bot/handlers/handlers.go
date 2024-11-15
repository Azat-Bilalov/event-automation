package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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
	Emails        []string `json:"user_ids"`
	StartDatetime string   `json:"start_datetime"`
	EndDatetime   string   `json:"end_datetime"`
	Timezone      int64    `json:"timezone"`
}

var userMessages = make(map[int64][]string)
var emails []string

func getMessageText(message *tgbotapi.Message) string {
	text := strings.TrimSpace(message.Text)
	return text
}

func Register(bot *tgbotapi.BotAPI, store storage.Storage, message *tgbotapi.Message) {
	email := message.CommandArguments()
	if !validate.IsEmail(email) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ghjdth!")
		bot.Send(msg)
		return
	}
	if store.IsExist(message.From.ID) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You already registered! Want to change email?") //TODO: change email logic
		bot.Send(msg)
		return
	}
	store.SetEmail(message.From.ID, email)
}

func CollectMessage(message *tgbotapi.Message, store storage.Storage) {
	if message.ForwardFrom == nil {
		return
	}
	log.Println("AVTOR", message.ForwardFrom.ID, message.ForwardFrom.FirstName)
	text := fmt.Sprintf("%s: %s", message.ForwardFrom.FirstName, strings.TrimSpace(message.Text))
	userMessages[message.From.ID] = append(userMessages[message.From.ID], text)
	emails = append(emails, store.GetEmail(message.ForwardFrom.ID))
}

func CreateEvent(bot *tgbotapi.BotAPI, store storage.Storage, message *tgbotapi.Message) error {
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
	body, err := ioutil.ReadAll(resp.Body)
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
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(body)
	delete(userMessages, userID)

	return nil
}
