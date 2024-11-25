.PHONY: build run test clean docker-build docker-run

# Go parameters
BINARY_NAME=dist/captain
MAIN_FILE=main.go

# Docker parameters
DOCKER_IMAGE=captain
DOCKER_TAG=latest

build:
	mkdir -p dist
	go build -o $(BINARY_NAME) $(MAIN_FILE)

run: build
	./$(BINARY_NAME) run

test:
	go test ./... -v

clean:
	go clean
	rm -f $(BINARY_NAME)

docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Development helpers
dev:
	go run $(MAIN_FILE) run

lint:
	golangci-lint run

fmt:
	go fmt ./...

# User management helpers
create-user: build
	$(BINARY_NAME) user create

update-password: build
	$(BINARY_NAME) user update-password

.DEFAULT_GOAL := build