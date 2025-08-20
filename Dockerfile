# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and ca-certificates for Go modules
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build Go binary for Linux, statically linked
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main .

ENV GIN_MODE=release
ENV PORT=8080

EXPOSE 8080

CMD ["./main"]
