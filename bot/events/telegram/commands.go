package telegram

import (
	"log"
	"net/mail"
	"strings"
)

const (
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command %s from %s", text, username)

	if isSaveEmailCmd(text) {
		// TODO: save email
	}

	// 1) Save email
	// 2) Check email + Change email
	// 3) make event
	// 4) help
	// 5) start

	switch text {

	}
}

func isSaveEmailCmd(text string) bool {
	return isEmail(text)
}

func isEmail(text string) bool {
	_, err := mail.ParseAddress(text)
	return err == nil
}
