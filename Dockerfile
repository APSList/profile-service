# --- BUILD STAGE ---
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o profile-service .

# --- FINAL STAGE ---
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/profile-service .
ENTRYPOINT ["./profile-service"]