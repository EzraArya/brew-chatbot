# ── Stage 1: Builder ──────────────────────────────────────────────────────────
# Use the official Go image — it has the compiler and all Go tools
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first
# Docker caches each step — if these haven't changed, it skips re-downloading deps
COPY go.mod go.sum ./

# Download dependencies (this layer is cached unless go.mod/go.sum change)
RUN go mod download

# Now copy the rest of the source code
COPY . .

# Compile the Go binary
# CGO_ENABLED=0 → pure Go binary, no C dependencies (needed for alpine)
# -o brew-chatbot → output file name
RUN CGO_ENABLED=0 GOOS=linux go build -o brew-chatbot .

# ── Stage 2: Runner ───────────────────────────────────────────────────────────
# Alpine is a minimal Linux distro (~5MB) — no Go compiler, just enough to run our binary
FROM alpine:latest

WORKDIR /app

# Copy ONLY the compiled binary from the builder stage
# Everything else (source code, Go tools) is left behind
COPY --from=builder /app/brew-chatbot .

# Document which port the app listens on (doesn't actually publish it)
EXPOSE 8080

# The command to run when the container starts
CMD ["./brew-chatbot"]
