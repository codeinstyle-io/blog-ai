# Stage 1: Build the binary
FROM golang:1.23.2-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Stage 2: Copy the binary to a minimal image
FROM scratch

# Copy the binary from the builder stage
COPY --from=builder /app/main /main

# Copy the templates and static files
COPY --from=builder /app/templates /templates
COPY --from=builder /app/blog/static /blog/static

# Command to run the executable
ENTRYPOINT ["/main"]