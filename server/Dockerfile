# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git for private modules
RUN apk add --no-cache git

# Copy go.mod files first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# Create writable directory in build stage (works because Alpine has a shell)
RUN mkdir -p /data

# ---------- Runtime stage ----------
FROM gcr.io/distroless/base-debian12

# Set working directory for binary
WORKDIR /app

# Copy server binary from builder
COPY --from=builder /app/server /app/server


# Expose port Railway will use
EXPOSE 8080

# Run as non-root user
USER nonroot:nonroot

# Start the Go server
ENTRYPOINT ["/app/server"]

