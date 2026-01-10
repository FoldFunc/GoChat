# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (sometimes needed for private modules)
RUN apk add --no-cache git

# Copy go mod files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the binary
# 👇 change ./cmd/app if your folder name is different
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o server ./cmd/server


# ---------- Runtime stage ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy only the binary
COPY --from=builder /app/server /app/server

# If you NEED a writable directory (sqlite, logs, etc)
# create it at runtime instead of baking files in
WORKDIR /data

EXPOSE 443

USER nonroot:nonroot

ENTRYPOINT ["/app/server"]

