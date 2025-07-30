FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY internal/app/googleDrive/service-account.json ./internal/app/googleDrive/
RUN CGO_ENABLED=0 go build -o /app/bin/main ./cmd/main/main.go
FROM alpine AS runner
COPY --from=builder /app/bin/main /app/main
COPY --from=builder /app/internal/app/googleDrive/service-account.json ./internal/app/googleDrive/
COPY --from=builder /app/config/config.yaml ./config/
CMD ["/app/main"]