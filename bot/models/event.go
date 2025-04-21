package models

type UserEvent struct {
	UserID   int64
	Language string
	Messages []string
	Emails   []string
	Timezone int
}
