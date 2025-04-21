package messages

import (
	"fmt"
	"strings"
)

type Message struct {
	Static  string
	Dynamic func(args ...interface{}) string
}

type Localization struct {
	Welcome                Message
	Help                   Message
	Error                  Message
	Success                Message
	Waiting                Message
	RegisterRequired       Message
	Register               Message
	AlreadyRegistered      Message
	CheckEmail             Message
	SuccessfulRegistration Message
	SuccessfulEmailChange  Message
	UnknownUser            Message
	Processing             Message
	WaitingEmail           Message
	CancelEmailChange      Message
}

var Localizations = map[string]Localization{
	"en": {
		Welcome: Message{Static: "Welcome to the bot!"},
		Help:    Message{Static: "Here are the available commands..."},
		Error:   Message{Static: "Something went wrong!"},
		Success: Message{
			Dynamic: func(args ...interface{}) string {
				eventLink := args[0].(string)
				inaccessibleClosed := args[1].([]string)
				inaccessibleNotInDB := args[2].([]string)

				closedList := strings.Join(inaccessibleClosed, ", ")
				notInDBList := strings.Join(inaccessibleNotInDB, ", ")

				return fmt.Sprintf(
					"Event created successfully! Link: %s.\nNot delivered to:\n- Closed accounts: %s\n- Not in DB: %s",
					eventLink, closedList, notInDBList,
				)
			},
		},
		RegisterRequired:       Message{Static: "You need to register your Gmail. Please send it in the format `example@gmail.com`."},
		Register:               Message{Static: "First, please register using the /register command."},
		AlreadyRegistered:      Message{Static: "You are already registered. Do you want to change your email? Enter /yes if yes."},
		CheckEmail:             Message{Static: "Please check the email you entered. The format or domain is incorrect (Gmail is required in this version)."},
		SuccessfulRegistration: Message{Static: "You have successfully registered!"},
		SuccessfulEmailChange:  Message{Static: "Your email has been changed!"},
		UnknownUser:            Message{Static: "Unknown user."},
		Processing:             Message{Static: "Starting to process the messages."},
		Waiting:                Message{Static: "Please forward a message for processing."},
		WaitingEmail:           Message{Static: "Please enter a new email address."},
		CancelEmailChange:      Message{Static: "Email change operation canceled."},
	},
	"ru": {
		Welcome: Message{Static: "Ку. Я помогу тебе с твоими событиями в google-календаре"},
		Help:    Message{Static: "Вот доступные команды..."},
		Error:   Message{Static: "Что-то пошло не так!"},
		Success: Message{
			Dynamic: func(args ...interface{}) string {
				eventLink := args[0].(string)
				inaccessibleClosed := args[1].([]string)
				inaccessibleNotInDB := args[2].([]string)

				closedList := strings.Join(inaccessibleClosed, ", ")
				notInDBList := strings.Join(inaccessibleNotInDB, ", ")

				return fmt.Sprintf(
					"Встреча успешно создана! Ссылка на событие: %s.\nНе доставлено:\n- Закрытые аккаунты: %s\n- Нет в базе: %s",
					eventLink, closedList, notInDBList,
				)
			},
		},
		RegisterRequired:       Message{Static: "Тебе нужно зарегистрировать свой Gmail, отправьте его в формате `example@gmail.com`."},
		Register:               Message{Static: "Сначала зарегистрируйтесь с помощью команды /register."},
		AlreadyRegistered:      Message{Static: "Вы уже зарегистрированы. Хотите сменить email? Введите /yes, если да."},
		CheckEmail:             Message{Static: "Проверьте введенную почту, неверный формат или домен (в текущей версии обязателен gmail)."},
		SuccessfulRegistration: Message{Static: "Вы успешно зарегистрировались!"},
		SuccessfulEmailChange:  Message{Static: "Почта успешно изменена!"},
		UnknownUser:            Message{Static: "Неизвестный пользователь."},
		Processing:             Message{Static: "Начинаю обработку сообщений."},
		Waiting:                Message{Static: "Перешлите сообщение для обработки."},
		WaitingEmail:           Message{Static: "Введите новый почтовый ящик."},
		CancelEmailChange:      Message{Static: "Операция смены email отменена."},
	},
}

// Получаем сообщение по ключу и языку, а также передаем аргументы для динамических сообщений
func GetMessage(lang, key string, args ...interface{}) string {
	// Получаем локализацию для заданного языка, если нет, используем английскую локализацию по умолчанию
	localization, ok := Localizations[lang]
	if !ok {
		localization = Localizations["en"] // fallback на английский
	}

	// В зависимости от ключа выбираем нужное сообщение
	var message Message
	fmt.Printf("you here %v \n", key)
	switch key {
	case "welcome":
		message = localization.Welcome
	case "help":
		message = localization.Help
	case "error":
		message = localization.Error
	case "success":
		message = localization.Success
	case "waiting":
		message = localization.Waiting
	case "register required":
		message = localization.RegisterRequired
	case "register":
		message = localization.Register
	case "already registered":
		message = localization.AlreadyRegistered
	case "check email":
		message = localization.CheckEmail
	case "successful registration":
		message = localization.SuccessfulRegistration
	case "successful email change":
		message = localization.SuccessfulEmailChange
	case "unknown user":
		message = localization.UnknownUser
	case "processing":
		message = localization.Processing
	case "waiting email":
		message = localization.WaitingEmail
	case "cancel email change":
		message = localization.CancelEmailChange
	default:
		return "Message not found ff"
	}

	// Если сообщение динамическое, вызываем функцию и передаем аргументы
	if message.Dynamic != nil {
		return message.Dynamic(args...)
	}

	// Если сообщение статическое, возвращаем его
	return message.Static
}
