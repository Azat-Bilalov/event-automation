package utils

import "time"

// сдвиг на часовой пояс от UTC
func GetDateWithTimezoneFromUTC(date time.Time, timezone int) time.Time {
	return date.Add(time.Duration(timezone) * time.Hour)
}

// дата по UTC без часового пояса
func GetDateWithoutTimezone(date time.Time, timezone int) time.Time {
	return date.Add(time.Duration(-timezone) * time.Hour)
}
