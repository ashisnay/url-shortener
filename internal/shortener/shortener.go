package shortener

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/url"
	"strings"

	"url-shortener/internal/storage"
)

const (
	// DefaultCodeLength is the default length of generated short codes
	DefaultCodeLength = 8
)

var (
	ErrEmptyURL   = errors.New("URL cannot be empty")
	ErrInvalidURL = errors.New("invalid URL format")
)

// Service handles URL shortening operations
type Service struct {
	storage    *storage.Storage
	codeLength int
}

// New creates a new shortener Service
func New(store *storage.Storage) *Service {
	return &Service{
		storage:    store,
		codeLength: DefaultCodeLength,
	}
}

// Shorten creates a short code for the given long URL
func (s *Service) Shorten(longURL string) (string, error) {
	// Validate input
	if strings.TrimSpace(longURL) == "" {
		return "", ErrEmptyURL
	}

	// Parse and validate URL format
	parsedURL, err := url.Parse(longURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", ErrInvalidURL
	}

	// Check if we already have a short code for this URL (idempotency)
	if existingCode, exists := s.storage.GetShortCode(longURL); exists {
		return existingCode, nil
	}

	// Generate deterministic short code using hash
	shortCode := s.generateCode(longURL)

	// Store the mapping
	s.storage.Store(shortCode, longURL)

	return shortCode, nil
}

// Expand retrieves the original long URL for a given short code
func (s *Service) Expand(shortCode string) (string, error) {
	if strings.TrimSpace(shortCode) == "" {
		return "", ErrEmptyURL
	}

	return s.storage.GetLongURL(shortCode)
}

// generateCode creates a deterministic short code from a URL using SHA-256
// Using a hash ensures the same URL always produces the same code (idempotency)
func (s *Service) generateCode(longURL string) string {
	hash := sha256.Sum256([]byte(longURL))
	// Use URL-safe base64 encoding and take first N characters
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	// Remove any padding characters and take first codeLength characters
	encoded = strings.TrimRight(encoded, "=")
	if len(encoded) > s.codeLength {
		return encoded[:s.codeLength]
	}
	return encoded
}
