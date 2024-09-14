package llm_service

import (
	"time"
)

const Model = "llama-3.1-8b-instant"

func GetSystemMessage(prefix string, language string, timezone int) string {
	today := time.Now()

	// Сегодня - это сегодняшняя дата для конкретного часового пояса юзера
	dateWithouTimezone := today.UTC().Add(time.Duration(timezone) * time.Hour)
	formattedDate := dateWithouTimezone.Format(time.RFC3339)

	return prefix + " Any time before eight o'clock is taken as pm. By default, the difference between the start Datetime and the end Datetime of the event is half an hour. Today: " + formattedDate + ". Language: " + language + ". Answer in only JSON with keys: title, start_datetime, end_datetime."
}
