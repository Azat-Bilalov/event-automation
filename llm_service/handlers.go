package llm_service

import (
	"encoding/json"
	"net/http"
)

func GetMeetDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем данные из запроса
	var llmReq LLMRequest
	if err := json.NewDecoder(r.Body).Decode(&llmReq); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Вызываем LLM сервис
	llmResponse, err := CallLLM(llmReq.Messages, llmReq.Language, llmReq.Timezone)
	if err != nil {
		http.Error(w, "Error calling LLM service", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(llmResponse); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}