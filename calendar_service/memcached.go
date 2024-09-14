package calendar_service

import "errors"

// FakeMemcached эмуляция хранилища никнеймов и почт
var FakeMemcached = map[string]string{
	"tolkachev_r": "tolkachev.rodion.03@gmail.com",
	"azat_bil": "az.bilalov@gmail.com",
	"azat": "az.bilalov@mail.ru",
}

func GetEmailByNickname(nickname string) (string, error) {
	email, exists := FakeMemcached[nickname]
	if !exists {
		return "", errors.New("email not found for nickname: " + nickname)
	}
	return email, nil
}
