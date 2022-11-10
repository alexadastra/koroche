package service

import (
	"context"

	"github.com/alexadastra/koroche/internal/app/models"
	"github.com/pkg/errors"
)

// Storage defines entity to store stuff
type Storage interface {
	AddURL(ctx context.Context, userURL models.UserURL, shortURL models.ShortURL) error
	GetUserURL(ctx context.Context, url models.ShortURL) (models.UserURL, error)
}

// Service handles all the domain
type Service struct {
	storage Storage
}

// NewService constructs new service
func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

// AddURL creates short url from user url, saves the pair and returns the short one
func (s *Service) AddURL(ctx context.Context, url models.UserURL) (models.ShortURL, error) {
	// TODO: hash URL here
	shortURL := models.NewShortURL("aaaaaaa.com")
	if err := s.storage.AddURL(ctx, url, shortURL); err != nil {
		return models.ShortURL(""), errors.Wrap(err, "failed to save URLs into db")
	}

	return shortURL, nil
}

// GetURL gets user url by short url
func (s *Service) GetURL(ctx context.Context, url models.ShortURL) (models.UserURL, error) {
	userURL, err := s.storage.GetUserURL(ctx, url)
	if err != nil {
		return models.UserURL(""), errors.Wrap(err, "failed to get user URL from db")
	}
	return userURL, nil
}
