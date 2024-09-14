package calendar_service

import (
	"encoding/json"
	"event-automation/config"
	"log"
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

	// Получаем email'ы по никнеймам
	var attendees []string
	for _, nickname := range req.Nicknames {
		email, err := GetEmailByNickname(nickname)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, "Invalid nickname: "+nickname, http.StatusBadRequest)
			return
		}
		attendees = append(attendees, email)
	}

	// Создаем клиента для Google Calendar
	credentialsPath := config.GetEnv("GOOGLE_CALENDAR_CREDENTIALS", "credentials.json")
	tokenPath := config.GetEnv("GOOGLE_CALENDAR_TOKEN", "token.json")

	client, err := NewCalendarClient(credentialsPath, tokenPath)
	if err != nil {
		http.Error(w, "Error creating Google Calendar client", http.StatusInternalServerError)
		return
	}

	// Парсим дату в формате "2022-01-01T15:04:05Z"
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}
	
	// Создаем событие
	event, err := CreateEvent(client, req.Title, req.Description, attendees, date)
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
