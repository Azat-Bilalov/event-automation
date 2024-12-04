package models

type CreateEventRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Emails      []string `json:"recipients_emails"`
	StartDate   string   `json:"start_datetime"`
	EndDate     string   `json:"end_datetime"`
	Timezone    int64    `json:"timezone"`
}

type LLMMessageRequest struct {
	Messages []string `json:"messages"`
	Language string   `json:"language"`
	Timezone int64    `json:"timezone"`
}

type LLMResponse struct {
	Title     string `json:"title"`
	StartDate string `json:"start_datetime"`
	EndDate   string `json:"end_datetime"`
}

type GoogleCalendarResponse struct {
	EventLink string `json:"event_link"`
}
