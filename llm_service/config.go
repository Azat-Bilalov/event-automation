package llm_service

import (
	"event-automation/utils"
	"log"
	"time"
)

const Model = "llama-3.1-8b-instant"

func GetSystemMessage(prefix string, language string, timezone int) string {
	today := time.Now().UTC()
	dateWithUserTimezone := utils.GetDateWithTimezoneFromUTC(today, timezone)
	formattedDate := dateWithUserTimezone.Format(time.RFC3339)

	log.Printf("Today: %s", formattedDate)

	return prefix + " Any time before eight o'clock is taken as pm. By default, the difference between the start Datetime and the end Datetime of the event is half an hour. Today: " + formattedDate + ". Language: " + language + ". Answer in only JSON with keys: title, start_datetime, end_datetime."
}
