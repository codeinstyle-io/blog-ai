# Build Stage
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go binary with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -o codeinstyle main.go

# Final Stage
FROM scratch

COPY --from=builder /app/codeinstyle /codeinstyle

EXPOSE 8080

ENTRYPOINT ["/codeinstyle"]
