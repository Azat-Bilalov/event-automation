package calendar_service

import "errors"

// FakeMemcached эмуляция хранилища никнеймов и почт
var FakeMemcached = map[string]string{
	"tolkachev_r": "tolkachev.rodion.03@gmail.com",
	"azat_bil":    "az.bilalov@gmail.com",
	"azat":        "az.bilalov@mail.ru",
	"ivan":        "elnikvolk908@gmail.com",
	"123":         "tolkachev.rodion.03@gmail.com",
	"345":         "az.bilalov@gmail.com",
}

func GetEmailByID(id string) (string, error) {
	email, exists := FakeMemcached[id]
	if !exists {
		return "", errors.New("email not found for id: " + id)
	}
	return email, nil
}
