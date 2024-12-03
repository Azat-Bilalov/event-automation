package utils

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ExtractName(user *tgbotapi.User) string {
	if user == nil {
		return "Unknown"
	}

	if user.UserName != "" {
		return "@" + user.UserName
	}

	if user.FirstName != "" && user.LastName != "" {
		return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}

	if user.FirstName != "" {
		return user.FirstName
	}

	return "Unknown"
}
