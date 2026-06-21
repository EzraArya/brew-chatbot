# 🍺 BrewBot — iOS

A SwiftUI iOS app that provides a conversational interface to BrewBot — an AI-powered brewing assistant powered by the Gemini API via a Go backend.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Swift 6 |
| UI Framework | SwiftUI (iOS 26+) |
| Architecture | Clean Architecture + MVVM |
| State Management | `@Observable` (iOS 17+) |
| Networking | `URLSession` + `async/await` |
| Concurrency | Swift 6 strict concurrency (`@MainActor`) |

---

## Project Structure

```
BrewBot/
├── App/
│   ├── BrewBotApp.swift            # App entry point
│   └── ContentView.swift          # Root view — renders ChatView
│
├── Data/
│   ├── Models/
│   │   └── DTO.swift              # Codable structs mapping to/from Go API JSON
│   └── Network/
│       └── ChatService.swift      # URLSession calls to Go backend (@MainActor)
│
├── Domain/
│   └── Models/
│       └── Message.swift          # Core domain model — one chat message
│
└── Presentation/
    ├── Chat/
    │   ├── ChatView.swift         # Main chat screen
    │   ├── ChatViewModel.swift    # @Observable ViewModel — owns state + calls service
    │   └── MessageBubble.swift    # Individual message bubble with Markdown rendering
    └── Components/
        └── InputBar.swift         # Text field + send button component
```

---

## Architecture Overview

```
ChatView (SwiftUI)
  │  user taps Send
  ▼
ChatViewModel (@Observable, @MainActor)
  │  calls service with history + new message
  ▼
ChatService (@MainActor)
  │  builds ChatRequestDTO, POST /chat
  ▼
Go Backend → Gemini API
  │
  ▼
ChatService
  │  decodes ChatResponseDTO → domain Message
  ▼
ChatViewModel
  │  appends reply to messages array
  ▼
ChatView re-renders with new message bubble
```

### Layer Responsibilities

| Layer | Responsibility | Knows About |
|---|---|---|
| `Presentation` | UI + state | Domain models only |
| `Domain` | Pure business models | Nothing external |
| `Data` | Network calls + JSON mapping | Domain models + DTOs |

> Dependencies only flow inward: `Presentation → Domain ← Data`

---

## Key Components

### `Message.swift` — Domain Model
The core concept of the app. A single chat bubble with a role (`user` or `model`), content, and timestamp. Conforms to `Identifiable`, `Equatable`, and `Sendable` for Swift 6 concurrency safety.

### `DTO.swift` — Data Transfer Objects
Raw `Codable` structs that map 1:1 to the Go API's JSON contract:
- `MessageDTO` — `{ role, content }`
- `ChatRequestDTO` — `{ history: [MessageDTO], userMessage }`
- `ChatResponseDTO` — `{ reply }`

Includes mapping extensions to convert between DTOs and domain `Message`.

### `ChatService.swift` — Network Layer
`@MainActor final class` responsible for all HTTP communication. Builds and fires the `POST /chat` request and maps the response back to a domain `Message`.

### `ChatViewModel.swift` — Presentation Logic
`@Observable @MainActor final class` that owns the full conversation state:
- `messages: [Message]` — the full history displayed in the UI
- `isLoading: Bool` — controls typing indicator
- `errorMessage: String?` — surfaces network errors

The iOS app maintains conversation history locally and sends it with every request — the backend is stateless.

### `MessageBubble.swift` & Widgets — UI Components
Renders a single chat bubble with:
- User messages right-aligned in accent color
- Bot messages left-aligned in secondary background
- Full Markdown rendering via `AttributedString` (line-by-line to preserve newlines)

If the backend sends a Tool Call (e.g. `[generate_brew_recipe]`), the `WidgetFactory` decodes it and the `WidgetRenderer` dynamically inserts a beautiful, interactive SwiftUI card (like `ManualBrewRecipeCardView`) directly into the chat stream, completely bypassing the text bubble.

### `TypewriterStreamer.swift` — Reusable Utility
A decoupled concurrency utility that acts as an async iterator. It drains Server-Sent Event streams instantly to avoid network timeouts, but yields characters back to the UI on a paced 15ms delay to create a smooth typewriter effect. It also intercepts Tool Calls to bypass the animation entirely.

### `InputBar.swift` — UI Component
Text input + send button with:
- Multi-line `TextField` (grows up to 5 lines)
- Animated send/loading state transition
- Disabled state when empty or loading

---

## Running Locally

### Prerequisites
- Xcode 26+
- iOS 26 device or simulator
- Go backend running locally (see [backend README](../backend/README.md))

### Setup

**1. Open the project:**
```
ios/BrewBot/BrewBot.xcodeproj
```

**2. Make sure the Go backend is running:**
```bash
cd backend
go run main.go
# Server starts at http://localhost:8080
```

**3. Select a simulator and run** (`⌘R`)

---

## Running on a Physical Device

The app uses `http://localhost:8080` by default, which only works in the simulator.

For a physical device:
1. Find your Mac's local IP: `ipconfig getifaddr en0`
2. Update `baseURL` in `ChatService.swift` to your Mac's IP (e.g. `http://192.168.1.42:8080`)
3. Both device and Mac must be on the same WiFi network
4. Add `NSAllowsLocalNetworking` to `Info.plist` to allow plain HTTP

---

## Key Design Decisions

### `@Observable` over `ObservableObject`
iOS 17+ `@Observable` macro replaces `ObservableObject` + `@Published`. Fewer annotations, better performance, and more granular view updates.

### `@MainActor` on `ChatService`
Xcode 26 project templates apply `@MainActor` to the whole module by default, causing synthesized `Codable` conformances to be `@MainActor` isolated. Using `@MainActor final class` on `ChatService` keeps it in the same isolation domain — no cross-actor conflicts.

### Stateless History Management
The iOS app owns the full conversation history in memory. Every request sends the complete history to the Go backend. This keeps the server stateless and the iOS code simple. Persistence (SwiftData, iCloud sync) can be added later.

### Line-by-line Markdown Rendering
SwiftUI's `Text` doesn't support block-level markdown (newlines, lists, code blocks). We split responses on `\n`, parse each line's inline markdown with `AttributedString`, then rejoin with actual newline characters.

### Explicit Dependency Injection
`ChatViewModel` always receives a `ChatService` via `init` — no hidden defaults. This makes dependencies visible and easy to mock in tests.

### Server-Driven Generative UI
Instead of hardcoding every possible UI state, we use Gemini Function Calling. The Go backend streams down structured JSON inside a `[tool_name]` wrapper. The iOS app intercepts this in the networking layer, decodes it into a type-safe Swift struct using `WidgetFactory`, and `WidgetRenderer` draws a beautiful native SwiftUI card. This allows the AI to dynamically spawn interactive widgets directly in the chat!
