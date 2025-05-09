package calendar_service

import (
	"encoding/json"
	"event-automation/config"
	"net/http"
	"time"
)

// CreateEventHandler обрабатывает запросы на создание событий
func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var attendees []string
	attendees = append(attendees, req.RecipientsEmails...)

	// Создаем клиента для Google Calendar
	credentialsPath := config.GetEnv("GOOGLE_CALENDAR_CREDENTIALS", "credentials.json")
	tokenPath := config.GetEnv("GOOGLE_CALENDAR_TOKEN", "token.json")

	client, err := NewCalendarClient(credentialsPath, tokenPath)
	if err != nil {
		http.Error(w, "Error creating Google Calendar client", http.StatusInternalServerError)
		return
	}

	startDatetime, err := time.Parse(config.LayoutDatetime, req.StartDatetime)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}
	endDatetime, err := time.Parse(config.LayoutDatetime, req.EndDatetime)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	eventRequest := &EventRequest{
		Title:         req.Title,
		Description:   req.Description,
		Attendees:     attendees,
		StartDatetime: startDatetime,
		EndDatetime:   endDatetime,
		Timezone:      req.Timezone,
	}

	// Создаем событие
	event, err := CreateEvent(client, eventRequest)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	// Возвращаем ссылку на событие в ответе
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"event_link": event.HtmlLink,
	})
}
