# 🍺 BrewBot

An AI-powered brewing assistant chatbot. Ask it anything about beer, coffee, tea, or kombucha brewing — recipes, techniques, ingredients, troubleshooting, and equipment.

Built with a Go backend, Google Gemini API, and a native iOS SwiftUI app. Features **Server-Driven Generative UI**: instead of just returning text, the backend streams strict JSON tool calls that the iOS app decodes into interactive, native SwiftUI widgets (like step-by-step recipe cards)!

---

## Project Structure

```
brew-chatbot/
├── backend/    # Go REST API — Gemini integration + chat endpoint
└── ios/        # iOS SwiftUI app — chat interface
```

---

## Segments

### [Backend](./backend/README.md)
The Go server that acts as the bridge between the iOS app and the Gemini API. Exposes a single `POST /chat` endpoint. Stateless — the iOS client manages conversation history.

| | |
|---|---|
| Language | Go 1.26+ |
| AI | Google Gemini API (`gemini-2.5-flash`) |
| HTTP | Go standard library (`net/http`) |
| Deploy | Docker (multi-stage build) |

### [iOS](./ios/README.md)
The native iOS SwiftUI app. Built with Clean Architecture + MVVM using Swift 6 strict concurrency.

| | |
|---|---|
| Language | Swift 6 |
| UI | SwiftUI (iOS 26+) |
| Architecture | Clean Architecture + MVVM |
| State | `@Observable` |

---

## How It Works

```
User types a brewing question
        │
        ▼
iOS App (SwiftUI)
  — manages conversation history locally
  — sends full history + new message to backend
        │
        │  POST /chat
        ▼
Go Backend
  — validates request
  — passes history + system prompt to Gemini
        │
        ▼
Gemini API (gemini-2.5-flash)
  — responds as BrewBot, a brewing expert
        │
        ▼
Go Backend → iOS App
  — returns reply as JSON
        │
        ▼
iOS App renders response with Markdown formatting
```

---

## Getting Started

See each segment's README for detailed setup instructions:

- **Backend:** [backend/README.md](./backend/README.md)
- **iOS:** [ios/README.md](./ios/README.md)

### Quick Start

**1. Start the backend:**
```bash
cd backend
cp .env.example .env   # add your Gemini API key
go run main.go
```

**2. Run the iOS app:**
```
Open ios/BrewBot/BrewBot.xcodeproj in Xcode
Select a simulator → ⌘R
```

---

## Prerequisites

| Tool | Purpose |
|---|---|
| Go 1.26+ | Backend runtime |
| Docker | Containerized backend deployment |
| Xcode 26+ | iOS development |
| Gemini API Key | AI responses ([Get one here](https://aistudio.google.com)) |
