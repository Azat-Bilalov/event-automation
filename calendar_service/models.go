package calendar_service

// CreateEventRequest структура запроса для создания события
type CreateEventRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Nicknames   []string `json:"nicknames"`
	Date 				string   `json:"date"`
}
