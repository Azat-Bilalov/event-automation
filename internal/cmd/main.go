package main

import (
	"event-automation/calendar_service"
	"event-automation/config"
	"event-automation/llm_service"
	"log"
	"net/http"
)

func main() {
	config.LoadEnv()

	http.HandleFunc("/create_event", calendar_service.CreateEventHandler)
	http.HandleFunc("/meet_details", llm_service.GetMeetDetailsHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
