# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Download dependencies first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o captain -ldflags='-w -s' .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /build/captain .

# Copy themes
COPY --from=builder /build/themes /app/themes

# Create necessary directories
RUN mkdir -p /app/media

EXPOSE 8080

ENTRYPOINT ["./captain"]
CMD ["run", "-b", "0.0.0.0", "-p", "8080"]