package telegram

import "event-automation/bot/clients/telegram"

type Processor struct {
	tg     *telegram
	offset int
}

func New(client *telegram.Client) *Processor {
	return &Processor{
		tg: client,
	}
}
