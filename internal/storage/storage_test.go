package storage

import (
	"sync"
	"testing"
)

func TestStorage_StoreAndRetrieve(t *testing.T) {
	store := New()

	shortCode := "abc123"
	longURL := "https://www.example.com/test"

	store.Store(shortCode, longURL)

	// Test GetLongURL
	gotURL, err := store.GetLongURL(shortCode)
	if err != nil {
		t.Fatalf("GetLongURL failed: %v", err)
	}
	if gotURL != longURL {
		t.Errorf("expected %s, got %s", longURL, gotURL)
	}

	// Test GetShortCode
	gotCode, exists := store.GetShortCode(longURL)
	if !exists {
		t.Fatal("GetShortCode returned false for existing URL")
	}
	if gotCode != shortCode {
		t.Errorf("expected %s, got %s", shortCode, gotCode)
	}
}

func TestStorage_GetLongURL_NotFound(t *testing.T) {
	store := New()

	_, err := store.GetLongURL("nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestStorage_GetShortCode_NotFound(t *testing.T) {
	store := New()

	_, exists := store.GetShortCode("https://nonexistent.com")
	if exists {
		t.Error("expected exists to be false for non-stored URL")
	}
}

func TestStorage_ConcurrentAccess(t *testing.T) {
	store := New()

	var wg sync.WaitGroup
	numGoroutines := 100

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			shortCode := string(rune('a' + n%26))
			longURL := "https://example.com/" + string(rune('a'+n%26))
			store.Store(shortCode, longURL)
		}(i)
	}
	wg.Wait()

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			shortCode := string(rune('a' + n%26))
			store.GetLongURL(shortCode)
		}(i)
	}
	wg.Wait()
}
