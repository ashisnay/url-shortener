package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"url-shortener/internal/shortener"
	"url-shortener/internal/storage"
)

// Handler holds HTTP handlers for the URL shortener service
type Handler struct {
	shortener *shortener.Service
	baseURL   string
}

// ShortenRequest represents the JSON request body for creating a short URL
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse represents the JSON response for a successful shortening
type ShortenResponse struct {
	ShortURL  string `json:"short_url"`
	ShortCode string `json:"short_code"`
	LongURL   string `json:"long_url"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// New creates a new Handler with the given base URL
func New(baseURL string, store *storage.Storage) *Handler {
	return &Handler{
		shortener: shortener.New(store),
		baseURL:   strings.TrimRight(baseURL, "/"),
	}
}

// Shorten handles POST requests to create shortened URLs
func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	shortCode, err := h.shortener.Shorten(req.URL)
	if err != nil {
		status := http.StatusBadRequest
		h.writeError(w, err.Error(), status)
		return
	}

	resp := ShortenResponse{
		ShortURL:  h.baseURL + "/" + shortCode,
		ShortCode: shortCode,
		LongURL:   req.URL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Redirect handles GET requests to redirect short URLs to their original destinations
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract short code from path (remove leading slash)
	shortCode := strings.TrimPrefix(r.URL.Path, "/")
	if shortCode == "" {
		h.writeError(w, "short code is required", http.StatusBadRequest)
		return
	}

	longURL, err := h.shortener.Expand(shortCode)
	if err != nil {
		if err == storage.ErrNotFound {
			h.writeError(w, "short URL not found", http.StatusNotFound)
			return
		}
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Perform HTTP redirect (302 Found)
	http.Redirect(w, r, longURL, http.StatusFound)
}

// writeError writes a JSON error response
func (h *Handler) writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
