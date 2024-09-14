package llm_service

// curl -X POST "https://api.groq.com/openai/v1/chat/completions" \
//      -H "Authorization: Bearer $GROQ_API_KEY" \
//      -H "Content-Type: application/json" \
//      -d '{"messages": [{"role": "user", "content": "Explain the importance of fast language models"}], "model": "llama3-8b-8192"}'

type OpenAICompletionsMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Стандарт запроса к OpenAI API
type OpenAICompletionsRequest struct {
	Model    string       							`json:"model"`
	Messages []OpenAICompletionsMessage `json:"messages"`
	Stream   bool         							`json:"stream"`
}

// Стандарт ответа от OpenAI API
type OpenAICompletionsResponse struct {
	Choices []struct {
		Message OpenAICompletionsMessage `json:"message"`
	} `json:"choices"`
}

// Запрос к LLM-сервису
type LLMRequest struct {
	Messages []string `json:"messages"`
	Language string `json:"language"`
	Timezone int `json:"timezone"`
}

// Ответ от LLM-сервиса
type LLMResponse struct {
	Title string `json:"title"`
	DateTime string `json:"datetime"`
}