package llm_service

type OpenAICompletionsMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Стандарт запроса к OpenAI API
type OpenAICompletionsRequest struct {
	Model    string                     `json:"model"`
	Messages []OpenAICompletionsMessage `json:"messages"`
	Stream   bool                       `json:"stream"`
}

// Стандарт ответа от OpenAI API
type OpenAICompletionsResponse struct {
	Choices []struct {
		Message OpenAICompletionsMessage `json:"message"`
	} `json:"choices"`
}

// Запрос на генерацию встречи к LLM-сервису
type LLMGenerationMeetRequest struct {
	Messages []string `json:"messages"`
	Language string   `json:"language"`
	Timezone int      `json:"timezone"`
}

// Запрос на уточнение встречи к LLM-сервису
type LLMClarificationMeetRequest struct {
	Meet     string `json:"meet"`
	Prompt   string `json:"prompt"`
	Language string `json:"language"`
	Timezone int    `json:"timezone"`
}

// Ответ с данными встречи от LLM-сервиса
type LLMMeetResponse struct {
	Title         string `json:"title"`
	StartDateTime string `json:"start_datetime"`
	EndDateTime   string `json:"end_datetime"`
}
