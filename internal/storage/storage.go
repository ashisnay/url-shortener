package storage

import (
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("short URL not found")
)

// Storage provides thread-safe storage for URL mappings
type Storage struct {
	mu          sync.RWMutex
	longToShort map[string]string
	shortToLong map[string]string
}

// New creates a new Storage instance
func New() *Storage {
	return &Storage{
		longToShort: make(map[string]string),
		shortToLong: make(map[string]string),
	}
}

// Store saves the URL mapping bidirectionally
func (s *Storage) Store(shortCode, longURL string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.longToShort[longURL] = shortCode
	s.shortToLong[shortCode] = longURL
}

// GetShortCode retrieves the short code for a long URL
func (s *Storage) GetShortCode(longURL string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	code, exists := s.longToShort[longURL]
	return code, exists
}

// GetLongURL retrieves the long URL for a short code
func (s *Storage) GetLongURL(shortCode string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	longURL, exists := s.shortToLong[shortCode]
	if !exists {
		return "", ErrNotFound
	}
	return longURL, nil
}
