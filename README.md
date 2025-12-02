# URL Shortener Service

A simple, idempotent URL shortening service written in Go.

## Features

- **URL Shortening**: Convert long URLs into short, easy-to-share codes
- **Idempotent**: The same long URL always produces the same short URL
- **HTTP Redirect**: Short URLs automatically redirect to their original destinations
- **Containerized**: Ready to deploy with Docker Compose

## API Endpoints

### Create Short URL

**POST** `/shorten`

Create a shortened URL from a long URL.

**Request:**
```json
{
  "url": "https://www.example.com/very/long/path/to/resource"
}
```

**Response (201 Created):**
```json
{
  "short_url": "http://localhost:8080/abc12345",
  "short_code": "abc12345",
  "long_url": "https://www.example.com/very/long/path/to/resource"
}
```

### Redirect to Original URL

**GET** `/{short_code}`

Redirects (HTTP 302) to the original long URL.

**Example:**
```
GET /abc12345 → Redirects to https://www.example.com/very/long/path/to/resource
```

## Quick Start

### Using Docker Compose (Recommended)

Deploy the service with a single command:

```bash
docker compose up -d
```

This will:
- Build the Docker image
- Start the service on port 8080

To stop the service:
```bash
docker-compose down
```

### Using Docker directly

```bash
# Build the image
docker build -t url-shortener .

# Run the container
docker run -d -p 8080:8080 --name url-shortener url-shortener
```

### Running Locally (without Docker)

Requires Go 1.21 or later.

```bash
# Run the service
go run main.go

# Or build and run
go build -o url-shortener .
./url-shortener
```

## Configuration

The service can be configured using environment variables:

| Variable   | Description                  | Default                  |
|------------|------------------------------|--------------------------|
| `PORT`     | Port to listen on            | `8080`                   |
| `BASE_URL` | Base URL for generated links | `http://localhost:8080`  |

## Running Tests

Run all tests:
```bash
go test ./...
```

Run tests with verbose output:
```bash
go test -v ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Usage Examples

### Create a short URL using curl

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com/search?q=golang"}'
```

**Response:**
```json
{
  "short_url": "http://localhost:8080/RqJhXsGq",
  "short_code": "RqJhXsGq",
  "long_url": "https://www.google.com/search?q=golang"
}
```

### Access the short URL

```bash
# Follow redirect
curl -L http://localhost:8080/RqJhXsGq

# See redirect headers
curl -I http://localhost:8080/RqJhXsGq
```

## Project Structure

```
url-shortener/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── Dockerfile              # Docker build instructions
├── docker-compose.yml      # Docker Compose configuration
├── README.md               # This file
└── internal/
    ├── handlers/           # HTTP handlers
    │   ├── handlers.go
    │   └── handlers_test.go
    ├── shortener/          # URL shortening logic
    │   ├── shortener.go
    │   └── shortener_test.go
    └── storage/            # In-memory storage
        ├── storage.go
        └── storage_test.go
```

## Design Decisions

1. **Idempotency**: Uses SHA-256 hash of the URL to generate deterministic short codes. This ensures the same URL always produces the same short code.

2. **In-Memory Storage**: Uses thread-safe in-memory maps for simplicity. For production use, consider adding persistent storage (Redis, PostgreSQL, etc.).

3. **URL-safe Encoding**: Uses Base64 URL-safe encoding for short codes to ensure they work in URLs without escaping.

4. **HTTP Redirect**: Uses HTTP 302 (Found) for redirects, allowing browsers to cache the redirect while still allowing URL updates.

## License

MIT License
