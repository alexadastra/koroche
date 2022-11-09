package models

import (
	"net/url"

	"github.com/pkg/errors"
)

// UserURL is URL passed by user
type UserURL string

// NewUserURL constructs new UserURL
func NewUserURL(rawURL string) (UserURL, error) {
	if _, err := url.Parse(rawURL); err != nil {
		return UserURL(""), errors.Wrap(err, "failed to validate url")
	}
	return UserURL(rawURL), nil
}

// ToString exports UserURL
func (uu UserURL) ToString() string {
	return string(uu)
}

// ShortURL is URL generated by system as a short representation of UserURL
type ShortURL string

// NewShortURL constructs new ShortURL
func NewShortURL(rawURL string) ShortURL {
	return ShortURL(rawURL)
}

// ToString exports ShortURL
func (su ShortURL) ToString() string {
	return string(su)
}