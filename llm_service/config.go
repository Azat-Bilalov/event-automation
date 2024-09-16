package llm_service

import (
	"log"
	"time"
)

const Model = "llama-3.1-8b-instant"

// функция с получением системного сообщения для LLM (промпта)
func GetSystemMessage(today time.Time, language string, timezone int) string {
	dateWithouTimezone := today.UTC().Add(time.Duration(timezone) * time.Hour)
	formattedDate := dateWithouTimezone.Format(time.RFC3339)
	log.Printf("You are generator of meet title and date. Today: %v Language: %v. Answer in only JSON: {title:..., datetime:...}", formattedDate, language)
	return "You are generator of meet title and date. Today: " + formattedDate + ". Language: " + language + ". Answer in only JSON: {title:..., datetime:...}"
}
