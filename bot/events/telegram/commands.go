package telegram

import (
	"log"
	"net/mail"
	"strings"
)

const (
	HelpCmd      = "/help"
	StartCmd     = "/start"
	SaveEmailCmd = "/saveEmail"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	user, exists := p.users[username]
	if !exists {
		user = &User{State: StateDefault}
		p.users[username] = user
	}

	log.Printf("got new command %s from %s (state: %d)", text, username, user.State)

	switch user.State {
	case StateDefault:
		return p.handleDefaultState(text, chatID, username, user)
	case StateAwaitingEmail:
		return p.handleEmailState(text, chatID, username, user)
	case StateAwaitingEventDetails:
		return p.handleEventState(text, chatID, username, user)
	default:
		return p.tg.SendMessage(chatID, "Unknown state")
	}

	// if isSaveEmailCmd(text) {
	// 	// TODO: save email
	// }

	// // 1) Save email
	// // 2) Check email + Change email
	// // 3) make event
	// // 4) help
	// // 5) start

	// switch text {
	// case StartCmd:
	// 	return p.sendHello(chatID)
	// case HelpCmd:
	// 	return p.sendHelp(chatID)
	// default:
	// 	return p.tg.SendMessage(chatID, "Unknown command")
	// }
}

func (p *Processor) handleDefaultState(text string, chatID int, username string, user *User) error {
	switch text {
	case StartCmd:
		return p.sendHello(chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	default:
		if strings.Contains(text, SaveEmailCmd) {
			user.State = StateAwaitingEmail
			return p.tg.SendMessage(chatID, "Напишите свой email")
		}
		return p.tg.SendMessage(chatID, "Unknown command")
	}
}

func (p *Processor) handleEmailState(text string, chatID int, username string, user *User) error {
	if isEmail(text) {
		user.Email = text
		user.State = StateDefault
		p.saveEmail(user.Email, []string{username})
		user.State = StateAwaitingEventDetails
		return p.tg.SendMessage(chatID, "Email сохранен")
	} else {
		return p.tg.SendMessage(chatID, "Некорректный email")
	}
}

func (p *Processor) handleEventState(text string, chatID int, username string, user *User) error {
	//TODO: обработка события, нужно изменить структуры, чтоб можно было отлавливать пересланные сообщения
	// user.State = StateDefault
	return p.tg.SendMessage(chatID, "Event created")
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) saveEmail(email string, usernames []string) {
	for _, key := range usernames {
		if !p.storage.IsExist(key) {
			p.storage.Save(key, email)
		}
	}
}

func isSaveEmailCmd(text string) bool {
	return isEmail(text)
}

func isEmail(text string) bool {
	_, err := mail.ParseAddress(text)
	return err == nil
}
