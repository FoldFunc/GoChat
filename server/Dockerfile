FROM golang:1.25 AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=intel64 go build \
  -o server ./cmd/server

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/server /app/server

EXPOSE 42069

USER nonroot:nonroot

ENTRYPOINT ["/app/server"]
