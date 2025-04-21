package calendar_service

import (
	"context"
	"encoding/json"
	"event-automation/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Извлекает токен, сохраняет токен, затем возвращает сгенерированного клиента
func getClient(config *oauth2.Config, tokenPath string) *http.Client {
	// Файл token.json хранит токены доступа и обновления пользователя и создается автоматически
	// при первом завершении процесса авторизации
	tok, err := tokenFromFile(tokenPath)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenPath, tok)
	}
	return config.Client(context.Background(), tok)
}

// Запрашивает токен из веба, затем возвращает полученный токен
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Извлекает токен из файла
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Сохраняет токен в файл
func saveToken(file string, token *oauth2.Token) {
	log.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// NewCalendarClient создает новый клиент для Google Calendar API
func NewCalendarClient(credentialsPath, tokenPath string) (*calendar.Service, error) {
	ctx := context.Background()
	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}
	client := getClient(config, tokenPath)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
		return nil, err
	}

	return srv, nil
}

type EventRequest struct {
	Title         string
	Description   string
	Attendees     []string
	StartDatetime time.Time
	EndDatetime   time.Time
	Timezone      int
}

// CreateEvent создает новое событие в календаре
func CreateEvent(srv *calendar.Service, req *EventRequest) (*calendar.Event, error) {
	startDate := utils.GetDateWithoutTimezone(req.StartDatetime, req.Timezone)
	endDate := utils.GetDateWithoutTimezone(req.EndDatetime, req.Timezone)

	event := &calendar.Event{
		Summary:     req.Title,
		Description: req.Description,
		Start: &calendar.EventDateTime{
			DateTime: startDate.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: endDate.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		Attendees: make([]*calendar.EventAttendee, len(req.Attendees)),
	}

	// Добавляем участников
	for i, email := range req.Attendees {
		event.Attendees[i] = &calendar.EventAttendee{Email: email}
	}

	calendarID := "primary"
	event, err := srv.Events.Insert(calendarID, event).Do()
	if err != nil {
		log.Fatalf("Unable to create event: %v", err)
		return nil, err
	}

	log.Printf("Event created: %s\n", event.HtmlLink)
	return event, nil
}
