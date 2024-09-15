package calendar_service

// CreateEventRequest структура запроса для создания события
type CreateEventRequest struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	UserIDs       []string `json:"user_ids"`
	StartDatetime string   `json:"start_datetime"`
	EndDatetime   string   `json:"end_datetime"`
	Timezone      int      `json:"timezone"`
}
