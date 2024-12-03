package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type UserEventContext struct {
	UserID   int64
	Language string
	Messages []string
	Emails   []string
	Timezone int64
}

func callLLMService(ctx *UserEventContext) (*LlmResponse, error) {
	data := MessageData{
		Messages: ctx.Messages,
		Language: ctx.Language,
		Timezone: ctx.Timezone,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal LLM request: %w", err)
	}

	resp, err := http.Post("http://localhost:8080/new_meet", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to call LLM service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("LLM service error: %s (status %d)", string(body), resp.StatusCode)
	}

	var llmResponse LlmResponse
	err = json.NewDecoder(resp.Body).Decode(&llmResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode LLM response: %w", err)
	}

	log.Println("LLM response received:", llmResponse)
	return &llmResponse, nil
}

func createCalendarEvent(ctx *UserEventContext, llmResponse *LlmResponse) (*CalendarResponse, error) {
	reqBody := EventData{
		Title:         llmResponse.Title,
		Desc:          llmResponse.Title,
		Emails:        ctx.Emails,
		StartDatetime: llmResponse.StartDatetime,
		EndDatetime:   llmResponse.EndDatetime,
		Timezone:      ctx.Timezone,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal calendar event request: %w", err)
	}

	resp, err := http.Post("http://localhost:8080/create_event", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to call calendar service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("calendar service error: %s (status %d)", string(body), resp.StatusCode)
	}

	var calendarResponse CalendarResponse
	err = json.NewDecoder(resp.Body).Decode(&calendarResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Calendar response: %w", err)
	}

	log.Println("Calendar response received: ", calendarResponse)
	return &calendarResponse, nil
}

func createEvent(ctx *UserEventContext) (*CalendarResponse, error) {
	// Создаём данные для LLM сервиса
	llmResponse, err := callLLMService(ctx)
	if err != nil {
		return nil, err
	}

	// Готовим запрос к календарю
	calendarResponse, err := createCalendarEvent(ctx, llmResponse)
	if err != nil {
		return nil, err
	}

	log.Printf("Event created successfully for user %d", ctx.UserID)
	return calendarResponse, nil
}
