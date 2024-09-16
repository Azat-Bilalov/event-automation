package llm_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// CallLLM отправляет запрос к LLM сервису (Groq)
func CallLLM(systemMessage string, userMessages []string) (string, error) {
	llmAPIURL := os.Getenv("LLM_API_URL")
	bearerToken := os.Getenv("LLM_API_TOKEN") // Получаем токен из переменных окружения

	llmMessages := make([]OpenAICompletionsMessage, len(userMessages)+1)

	llmMessages[0] = OpenAICompletionsMessage{Role: "system", Content: systemMessage}
	for i, message := range userMessages {
		llmMessages[i+1] = OpenAICompletionsMessage{Role: "user", Content: message}
	}

	llmRequest := OpenAICompletionsRequest{
		Model:    Model,
		Messages: llmMessages,
		Stream:   false,
	}

	// Кодируем запрос в JSON
	requestBody, err := json.Marshal(llmRequest)
	if err != nil {
		return "", err
	}

	// Создаём новый HTTP запрос
	req, err := http.NewRequest("POST", llmAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	// Создаём HTTP клиент и выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Декодируем ответ
	var llmResponse OpenAICompletionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&llmResponse); err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM service returned non-OK status: %d", resp.StatusCode)
	}

	return llmResponse.Choices[0].Message.Content, nil
}
