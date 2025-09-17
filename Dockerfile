# Stage 1: build
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy go.mod và go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ source
COPY . .

# Build binary từ cmd/main.go
RUN go build -o hestia ./cmd/main.go

# Stage 2: run
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/hestia .

EXPOSE 8080
CMD ["./hestia"]
