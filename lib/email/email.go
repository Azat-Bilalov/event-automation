package email

import (
	"errors"
	"net/mail"
	"strings"
)

func IsGmail(addr *mail.Address) error {
	if !strings.HasSuffix(addr.Address, "@gmail.com") {
		return errors.New("email must has gmail domen")
	}
	return nil
}

func IsEmail(text string) bool {
	addr, err := mail.ParseAddress(text)
	if err != nil {
		return err == nil
	}
	err = IsGmail(addr)
	return err == nil
}
