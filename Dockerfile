# Stage 1: Build the binary
FROM golang:1.23-bookworm AS builder
ENV CGO_ENABLED=1

RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      build-essential \
      libsqlite3-dev

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

FROM debian:bookworm
RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      libsqlite3-0

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main /main

# Copy the templates and static files
COPY --from=builder /app/templates /templates
COPY --from=builder /app/static /static
COPY --from=builder /app/data /data

EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/main"]