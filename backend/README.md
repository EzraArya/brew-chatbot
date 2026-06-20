# 🍺 Brew Chatbot — Backend

A Go backend that powers the Brew Chatbot — an AI-powered brewing assistant built on the Gemini API. It exposes a simple REST API consumed by the iOS app.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.26+ |
| AI | Google Gemini API (`gemini-2.5-flash`) |
| HTTP | Go standard library (`net/http`) |
| Containerization | Docker (multi-stage build) |

---

## Project Structure

```
backend/
├── main.go                     # Entry point — starts the HTTP server
├── Dockerfile                  # Multi-stage Docker build
├── .env                        # Local secrets (never committed)
├── config/
│   └── config.go               # Loads environment variables into a typed Config struct
├── gemini/
│   └── client.go               # Gemini API client wrapper + system prompt
├── handler/
│   └── chat.go                 # HTTP handler for POST /chat
└── internal/
    └── httputil/
        └── response.go         # Shared JSON response helpers
```

---

## Architecture Overview

```
iOS App
  │
  │  POST /chat  { history: [...], userMessage: "..." }
  ▼
handler/chat.go         — validates request, calls Gemini client
  │
  ▼
gemini/client.go        — builds chat session with system prompt + history, calls Gemini API
  │
  ▼
Gemini API (gemini-2.5-flash)
  │
  ▼
handler/chat.go         — returns { reply: "..." }
  │
  ▼
iOS App
```

---

## API Endpoints

### `GET /health`
Health check. Returns `200 OK` if the server is running.

```bash
curl http://localhost:8080/health
```

---

### `POST /chat`
Send a message to BrewBot. The iOS client maintains the full conversation history and sends it with every request, keeping the server stateless.

**Request body:**
```json
{
  "history": [
    { "role": "user",  "content": "What hops work for an IPA?" },
    { "role": "model", "content": "Cascade and Centennial are great choices..." }
  ],
  "userMessage": "What about bitterness levels?"
}
```

**Response:**
```json
{
  "reply": "For bitterness in an IPA, you're looking at IBUs..."
}
```

**Error response:**
```json
{
  "error": "userMessage cannot be empty"
}
```

---

## Running Locally

### Prerequisites
- Go 1.26+
- A Gemini API key from [Google AI Studio](https://aistudio.google.com)

### Setup

**1. Clone and navigate:**
```bash
cd backend
```

**2. Create your `.env` file:**
```bash
cp .env.example .env
# then add your key inside .env
```

**3. Install dependencies:**
```bash
go mod download
```

**4. Run:**
```bash
go run main.go
```

Server starts at `http://localhost:8080`

---

## Running with Docker

**Build the image:**
```bash
docker build -t brew-chatbot .
```

**Run the container:**
```bash
docker run -p 8080:8080 \
  -e GEMINI_API_KEY=your_key_here \
  brew-chatbot
```

> The `.env` file is never baked into the Docker image. Secrets are always injected at runtime via the `-e` flag.

---

## Key Design Decisions

### Stateless Chat
The server does not store conversation history. The iOS client sends the full history with every request. This keeps the backend simple, scalable, and easy to deploy.

### Standard Library HTTP
No web framework is used — Go's `net/http` is production-grade for this use case. A lightweight router (e.g. `chi`) can be added later if path parameters or middleware are needed.

### Multi-stage Docker Build
The builder stage compiles the Go binary; the runner stage contains only the binary on a minimal Alpine Linux base. This results in a ~15MB final image with no source code or Go toolchain exposed.

### Internal Package
Shared HTTP helpers live in `internal/httputil/` — Go's `internal/` convention prevents these from being imported by code outside this module.

---

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `GEMINI_API_KEY` | ✅ Yes | Your Google Gemini API key |

---

## What's Next

- iOS frontend (SwiftUI)
- CI/CD pipeline
- Deploy to cloud (Fly.io / GCP / AWS)
