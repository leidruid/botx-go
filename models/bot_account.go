package models

import (
	"errors"
	"net/url"

	"github.com/google/uuid"
)

// BotAccount represents bot credentials without secret.
type BotAccount struct {
	ID     uuid.UUID
	CTSURL string
}

// BotAccountWithSecret contains bot id, CTS url and secret key.
type BotAccountWithSecret struct {
	ID        uuid.UUID
	CTSURL    string
	SecretKey string
}

// Host returns hostname parsed from CTSURL.
func (b BotAccountWithSecret) Host() (string, error) {
	u, err := url.Parse(b.CTSURL)
	if err != nil {
		return "", err
	}
	if u.Hostname() == "" {
		return "", errors.New("could not parse host from CTSURL")
	}
	return u.Hostname(), nil
}
