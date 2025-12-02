package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"url-shortener/internal/handlers"
	"url-shortener/internal/storage"
)

func main() {
	// Get configuration from environment variables with defaults
	port := getEnv("PORT", "8080")
	baseURL := getEnv("BASE_URL", "http://localhost:"+port)

	// Initialize storage and handlers
	store := storage.New()
	handler := handlers.New(baseURL, store)

	// Setup routes
	mux := http.NewServeMux()

	// GET /api/v1/health - Health check endpoint
	mux.HandleFunc("/api/v1/health", handler.Health)

	// POST /api/v1/shortURL - Create a short URL
	mux.HandleFunc("/api/v1/shortURL", handler.Shorten)

	// GET /{shortCode} - Redirect to original URL
	mux.HandleFunc("/", handler.Redirect)

	// Start server
	addr := ":" + port
	fmt.Printf("URL Shortener service starting on %s\n", addr)
	fmt.Printf("Base URL: %s\n", baseURL)
	fmt.Println("Endpoints:")
	fmt.Println("  GET  /api/v1/health - Health check")
	fmt.Println("  POST /api/v1/shortURL - Create a short URL")
	fmt.Println("  GET  /{shortCode} - Redirect to original URL")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
