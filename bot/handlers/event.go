package handlers

import (
	"bytes"
	"encoding/json"
	"event-automation/bot/models"
	"fmt"
	"io"
	"log"
	"net/http"
)

func callLLMService(event *models.UserEvent) (*models.LLMResponse, error) {
	data := models.LLMMessageRequest{
		Messages: event.Messages,
		Language: event.Language,
		Timezone: event.Timezone,
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
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read LLM service error response: %w", err)
		}
		return nil, fmt.Errorf("LLM service error: %s (status %d)", string(body), resp.StatusCode)
	}

	var llmResponse models.LLMResponse
	err = json.NewDecoder(resp.Body).Decode(&llmResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode LLM response: %w", err)
	}

	log.Println("LLM response received:", llmResponse)
	return &llmResponse, nil
}

func createCalendarEvent(event *models.UserEvent, llmResponse *models.LLMResponse) (*models.GoogleCalendarResponse, error) {
	reqBody := models.CreateEventRequest{
		Title:       llmResponse.Title,
		Description: llmResponse.Title,
		Emails:      event.Emails,
		StartDate:   llmResponse.StartDate,
		EndDate:     llmResponse.EndDate,
		Timezone:    event.Timezone,
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

	var calendarResponse models.GoogleCalendarResponse
	err = json.NewDecoder(resp.Body).Decode(&calendarResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Calendar response: %w", err)
	}

	log.Println("Calendar response received: ", calendarResponse)
	return &calendarResponse, nil
}

func handleEventCreation(event *models.UserEvent) (*models.GoogleCalendarResponse, error) {
	llmResponse, err := callLLMService(event)
	if err != nil {
		return nil, err
	}

	calendarResponse, err := createCalendarEvent(event, llmResponse)
	if err != nil {
		return nil, err
	}

	log.Printf("Event created successfully for user %d", event.UserID)
	return calendarResponse, nil
}
