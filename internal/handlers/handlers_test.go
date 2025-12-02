package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"url-shortener/internal/storage"
)

func TestShorten_Success(t *testing.T) {
	store := storage.New()
	handler := New("http://localhost:8080", store)

	reqBody := `{"url": "https://www.example.com/test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shortURL", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Shorten(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var resp ShortenResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ShortCode == "" {
		t.Error("expected non-empty short code")
	}
	if resp.LongURL != "https://www.example.com/test" {
		t.Errorf("expected long URL https://www.example.com/test, got %s", resp.LongURL)
	}
	if !strings.HasPrefix(resp.ShortURL, "http://localhost:8080/") {
		t.Errorf("short URL should start with base URL, got %s", resp.ShortURL)
	}
}

func TestShorten_Idempotency(t *testing.T) {
	store := storage.New()
	handler := New("http://localhost:8080", store)

	reqBody := `{"url": "https://www.example.com/idempotent"}`

	// First request
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/shortURL", bytes.NewBufferString(reqBody))
	rec1 := httptest.NewRecorder()
	handler.Shorten(rec1, req1)

	var resp1 ShortenResponse
	json.NewDecoder(rec1.Body).Decode(&resp1)

	// Second request with same URL
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/shortURL", bytes.NewBufferString(reqBody))
	rec2 := httptest.NewRecorder()
	handler.Shorten(rec2, req2)

	var resp2 ShortenResponse
	json.NewDecoder(rec2.Body).Decode(&resp2)

	if resp1.ShortCode != resp2.ShortCode {
		t.Errorf("idempotency failed: got %s and %s", resp1.ShortCode, resp2.ShortCode)
	}
}

func TestShorten_InvalidJSON(t *testing.T) {
	store := storage.New()
	handler := New("http://localhost:8080", store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/shortURL", bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()

	handler.Shorten(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestShorten_EmptyURL(t *testing.T) {
	store := storage.New()
	handler := New("http://localhost:8080", store)

	reqBody := `{"url": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shortURL", bytes.NewBufferString(reqBody))
	rec := httptest.NewRecorder()

	handler.Shorten(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestShorten_InvalidURL(t *testing.T) {
	store := storage.New()
	handler := New("http://localhost:8080", store)

	reqBody := `{"url": "not-a-valid-url"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shortURL", bytes.NewBufferString(reqBody))
	rec := httptest.NewRecorder()

	handler.Shorten(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestShorten_WrongMethod(t *testing.T) {
	store := storage.New()
	handler := New("http://localhost:8080", store)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/shortURL", nil)
	rec := httptest.NewRecorder()

	handler.Shorten(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestRedirect_Success(t *testing.T) {
	store := storage.New()
	handler := New("http://localhost:8080", store)

	// First, create a short URL
	reqBody := `{"url": "https://www.example.com/redirect-test"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/shortURL", bytes.NewBufferString(reqBody))
	createRec := httptest.NewRecorder()
	handler.Shorten(createRec, createReq)

	var createResp ShortenResponse
	json.NewDecoder(createRec.Body).Decode(&createResp)

	// Now test the redirect
	redirectReq := httptest.NewRequest(http.MethodGet, "/"+createResp.ShortCode, nil)
	redirectRec := httptest.NewRecorder()
	handler.Redirect(redirectRec, redirectReq)

	if redirectRec.Code != http.StatusFound {
		t.Errorf("expected status %d, got %d", http.StatusFound, redirectRec.Code)
	}

	location := redirectRec.Header().Get("Location")
	if location != "https://www.example.com/redirect-test" {
		t.Errorf("expected redirect to https://www.example.com/redirect-test, got %s", location)
	}
}

func TestRedirect_NotFound(t *testing.T) {
	store := storage.New()
	handler := New("http://localhost:8080", store)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rec := httptest.NewRecorder()

	handler.Redirect(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestRedirect_EmptyCode(t *testing.T) {
	store := storage.New()
	handler := New("http://localhost:8080", store)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.Redirect(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestRedirect_WrongMethod(t *testing.T) {
	store := storage.New()
	handler := New("http://localhost:8080", store)

	req := httptest.NewRequest(http.MethodPost, "/somecode", nil)
	rec := httptest.NewRecorder()

	handler.Redirect(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}
