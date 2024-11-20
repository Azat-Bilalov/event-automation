package email

import "net/mail"

func IsEmail(text string) bool {
	_, err := mail.ParseAddress(text)
	return err == nil
}
