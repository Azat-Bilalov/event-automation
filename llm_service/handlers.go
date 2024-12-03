package llm_service

import (
	"encoding/json"
	"net/http"
)

func GetNewMeetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// TODO: выпилить!!! для тестов:
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&LLMMeetResponse{
		Title:         "Встреча",
		StartDatetime: "2023-12-03T15:04:05",
		EndDatetime:   "2023-12-03T15:34:05",
	}); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	// 	// Получаем данные из запроса
	// 	var llmReq LLMGenerationMeetRequest
	// 	if err := json.NewDecoder(r.Body).Decode(&llmReq); err != nil {
	// 		http.Error(w, "Bad request", http.StatusBadRequest)
	// 		return
	// 	}

	// 	// Генерируем системное сообщение
	// 	systemMessage := GetSystemMessage("You are the generator of the title and date time of the new meeting.", llmReq.Language, llmReq.Timezone)

	// 	// Вызываем LLM сервис
	// 	llmResponse, err := CallLLM(systemMessage, llmReq.Messages)
	// 	if err != nil {
	// 		http.Error(w, "Error calling LLM service", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// Парсим ответ от LLM сервиса
	// 	// TODO: В дальнейшем, если llm отдаёт неверный формат, нужно перезапрашивать
	// 	var response LLMMeetResponse
	// 	if err := json.Unmarshal([]byte(llmResponse), &response); err != nil {
	// 		http.Error(w, "Error parsing LLM response", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// Отправляем ответ
	// 	w.Header().Set("Content-Type", "application/json")
	// 	if err := json.NewEncoder(w).Encode(response); err != nil {
	// 		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	// 		return
	// 	}
}

func GetDetailedMeetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Получаем данные из запроса
	var llmReq LLMClarificationMeetRequest
	if err := json.NewDecoder(r.Body).Decode(&llmReq); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Генерируем системное сообщение
	systemMessage := GetSystemMessage("You clarify the details for the meeting title and Datetime.", llmReq.Language, llmReq.Timezone)

	// Вызываем LLM сервис
	llmResponse, err := CallLLM(systemMessage, []string{llmReq.Meet, llmReq.Prompt})
	if err != nil {
		http.Error(w, "Error calling LLM service", http.StatusInternalServerError)
		return
	}

	// Парсим ответ от LLM сервиса
	// TODO: В дальнейшем, если llm отдаёт неверный формат, нужно перезапрашивать
	var response LLMMeetResponse
	if err := json.Unmarshal([]byte(llmResponse), &response); err != nil {
		http.Error(w, "Error parsing LLM response", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
