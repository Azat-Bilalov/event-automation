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
	http.HandleFunc("/new_meet", llm_service.GetNewMeetHandler)
	http.HandleFunc("/detailed_meet", llm_service.GetDetailedMeetHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
