package shortener

import (
	"testing"

	"url-shortener/internal/storage"
)

func TestShorten_ValidURL(t *testing.T) {
	store := storage.New()
	svc := New(store)

	longURL := "https://www.example.com/very/long/path/to/resource"
	shortCode, err := svc.Shorten(longURL)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if shortCode == "" {
		t.Fatal("expected non-empty short code")
	}
	if len(shortCode) != DefaultCodeLength {
		t.Errorf("expected short code length %d, got %d", DefaultCodeLength, len(shortCode))
	}
}

func TestShorten_Idempotency(t *testing.T) {
	store := storage.New()
	svc := New(store)

	longURL := "https://www.example.com/test-idempotency"

	// Shorten the same URL multiple times
	shortCode1, err := svc.Shorten(longURL)
	if err != nil {
		t.Fatalf("first shorten failed: %v", err)
	}

	shortCode2, err := svc.Shorten(longURL)
	if err != nil {
		t.Fatalf("second shorten failed: %v", err)
	}

	shortCode3, err := svc.Shorten(longURL)
	if err != nil {
		t.Fatalf("third shorten failed: %v", err)
	}

	// All short codes should be identical
	if shortCode1 != shortCode2 || shortCode2 != shortCode3 {
		t.Errorf("idempotency failed: got different codes %s, %s, %s", shortCode1, shortCode2, shortCode3)
	}
}

func TestShorten_DifferentURLs(t *testing.T) {
	store := storage.New()
	svc := New(store)

	url1 := "https://www.example.com/path1"
	url2 := "https://www.example.com/path2"

	code1, err := svc.Shorten(url1)
	if err != nil {
		t.Fatalf("shorten url1 failed: %v", err)
	}

	code2, err := svc.Shorten(url2)
	if err != nil {
		t.Fatalf("shorten url2 failed: %v", err)
	}

	if code1 == code2 {
		t.Errorf("different URLs should produce different codes, got %s for both", code1)
	}
}

func TestShorten_EmptyURL(t *testing.T) {
	store := storage.New()
	svc := New(store)

	_, err := svc.Shorten("")
	if err != ErrEmptyURL {
		t.Errorf("expected ErrEmptyURL, got %v", err)
	}
}

func TestShorten_WhitespaceURL(t *testing.T) {
	store := storage.New()
	svc := New(store)

	_, err := svc.Shorten("   ")
	if err != ErrEmptyURL {
		t.Errorf("expected ErrEmptyURL for whitespace-only URL, got %v", err)
	}
}

func TestShorten_InvalidURL(t *testing.T) {
	store := storage.New()
	svc := New(store)

	testCases := []string{
		"not-a-valid-url",
		"example.com", // missing scheme
		"://missing-scheme.com",
		"http://", // missing host
	}

	for _, tc := range testCases {
		_, err := svc.Shorten(tc)
		if err != ErrInvalidURL {
			t.Errorf("expected ErrInvalidURL for %q, got %v", tc, err)
		}
	}
}

func TestExpand_ValidCode(t *testing.T) {
	store := storage.New()
	svc := New(store)

	longURL := "https://www.example.com/expand-test"
	shortCode, err := svc.Shorten(longURL)
	if err != nil {
		t.Fatalf("shorten failed: %v", err)
	}

	expandedURL, err := svc.Expand(shortCode)
	if err != nil {
		t.Fatalf("expand failed: %v", err)
	}

	if expandedURL != longURL {
		t.Errorf("expected %s, got %s", longURL, expandedURL)
	}
}

func TestExpand_NotFound(t *testing.T) {
	store := storage.New()
	svc := New(store)

	_, err := svc.Expand("nonexistent")
	if err != storage.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestExpand_EmptyCode(t *testing.T) {
	store := storage.New()
	svc := New(store)

	_, err := svc.Expand("")
	if err != ErrEmptyURL {
		t.Errorf("expected ErrEmptyURL for empty code, got %v", err)
	}
}
