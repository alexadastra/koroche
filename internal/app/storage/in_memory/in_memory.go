package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/alexadastra/koroche/internal/app/models"
)

// Storage stores URL pairs in memory
type Storage struct {
	mapping map[models.ShortURL]models.UserURL
	m       *sync.RWMutex
}

// NewStorage constructs new Storage
func NewStorage() *Storage {
	return &Storage{
		mapping: make(map[models.ShortURL]models.UserURL, 0),
		m:       &sync.RWMutex{},
	}
}

// AddURL adds URL pair to map
func (s *Storage) AddURL(ctx context.Context, userURL models.UserURL, shortURL models.ShortURL) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.mapping[shortURL] = userURL
	return nil
}

// GetUserURL fetches user URL by short URL
func (s *Storage) GetUserURL(ctx context.Context, url models.ShortURL) (models.UserURL, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	if userURL, ok := s.mapping[url]; ok {
		return userURL, nil
	}

	return models.UserURL(""), errors.New("url not found")
}
