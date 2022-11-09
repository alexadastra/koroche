package service

import (
	"context"

	"github.com/alexadastra/koroche/internal/app/models"
	"github.com/pkg/errors"
)

// Service handles all the domain
type Service struct{}

// AddURL creates short url from user url, saves the pair and returns the short one
func (s *Service) AddURL(ctx context.Context, url models.UserURL) (models.ShortURL, error) {
	return models.ShortURL(""), errors.New("unimplemented")
}

// GetURL gets user url by short url
func (s *Service) GetURL(ctx context.Context, url models.ShortURL) (models.UserURL, error) {
	return models.UserURL(""), errors.New("unimplemented")
}
