package messages

var Messages = map[string]map[string]string{
	"en": {
		"welcome": "Welcome to the bot!",
		"help":    "Here are the available commands...",
		"error":   "Something went wrong!",
	},
	"ru": {
		"welcome": "Ку. Я помогу тебе с твоими событиями в google-календаре",
		"register required": "Тебе нужно зарегистрировать свой gmail, " +
			"отправьте его в формате `example@gmail.com`",
		"register": "Сначала зарегистрируйтесь с помощью команды /register",
		"already registered": "Вы уже зарегистрированы." +
			"Хотите сменить email?	Введите /yes, если да",
		"check email": "Проверьте введенную почту, неверный формат " +
			"или домен (в текущей версии обязателен gmail)",
		"successful registration": "Вы зарегистрировались!",
		"successful email change": "Почта изменена!",
		"unknown user":            "",
		"processing":              "Начинаю обработку сообщений",
		"waiting":                 "Перешлите сообщение для обработки",
		"waiting email":           "Введите новый почтовый ящик",
		"cancel email change":     "Отмена операции смены email",
		"success":                 "Встреча успешно создана",
		"help":                    "Вот доступные команды...",
		"error":                   "Что-то пошло не так!",
	},
}

func GetMessage(lang string, key string) string {
	if langMessages, ok := Messages[lang]; ok {
		if message, ok := langMessages[key]; ok {
			return message
		}
	}
	message := Messages[lang][key]
	return message
}
