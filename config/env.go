package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv загружает переменные из .env файла
func LoadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// GetEnv возвращает значение переменной окружения или значение по умолчанию, если переменная не найдена
func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
