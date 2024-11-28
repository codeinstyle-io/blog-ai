.PHONY: build run clean \
        docker-build docker-run \
        dev lint fmt \
        test test-coverage \
        create-user update-password quality

# Variables
BINARY_NAME=dist/captain
MAIN_FILE=main.go
DOCKER_IMAGE=captain
DOCKER_TAG=latest

# Build commands
build:
	mkdir -p dist
	go build -o $(BINARY_NAME) $(MAIN_FILE)

clean:
	go clean
	rm -f $(BINARY_NAME)

# Run commands
run: build
	./$(BINARY_NAME) run

dev:
	go run $(MAIN_FILE) run

# Docker commands
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Test commands
test:
	go test ./... -v

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Code quality commands
lint:
	golangci-lint run

fmt:
	go fmt ./...

quality: fmt lint test

# User management commands
create-user: build
	$(BINARY_NAME) user create

update-password: build
	$(BINARY_NAME) user update-password

.DEFAULT_GOAL := build