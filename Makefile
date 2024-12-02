.PHONY: build run clean \
        docker-build docker-run \
        dev lint fmt \
        test test-coverage \
        create-user update-password quality \
        release

# Variables
BINARY_NAME=dist/captain
MAIN_FILE=main.go
DOCKER_IMAGE=captain
DOCKER_TAG=latest
GIN_MODE=release

# Build commands
build:
	mkdir -p dist
	go build -o $(BINARY_NAME) $(MAIN_FILE)

clean:
	go clean
	rm -rf dist

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

# Release commands
release: test
	@echo "Building release binaries..."
	mkdir -p dist/zip
	GIN_MODE=release GOOS=linux GOARCH=amd64 go build -v -o "dist/captain-linux-amd64/captain" .
	GIN_MODE=release GOOS=linux GOARCH=arm64 go build -v -o "dist/captain-linux-arm64/captain" .
	GIN_MODE=release GOOS=darwin GOARCH=amd64 go build -v -o "dist/captain-darwin-amd64/captain" .
	GIN_MODE=release GOOS=darwin GOARCH=arm64 go build -v -o "dist/captain-darwin-arm64/captain" .
	GIN_MODE=release GOOS=windows GOARCH=amd64 go build -v -o "dist/captain-windows-amd64/captain.exe" .
	GIN_MODE=release GOOS=windows GOARCH=arm64 go build -v -o "dist/captain-windows-arm64/captain.exe" .
	cd dist && \
	for dir in */; do \
		if [ -d "$$dir" ] && [ "$$dir" != "zip/" ]; then \
			platform=$${dir%/}; \
			zip -r "zip/captain-$${platform##captain-}.zip" "$$platform"; \
		fi \
	done
	@echo "Release build complete. Binaries are in dist/ and zip archives are in dist/zip/"

# User management commands
create-user: build
	$(BINARY_NAME) user create

update-password: build
	$(BINARY_NAME) user update-password

.DEFAULT_GOAL := build

release: test
	@echo "Building release binaries..."
	mkdir -p dist/zip
	GIN_MODE=release GOOS=linux GOARCH=amd64 go build -v -o "dist/captain-linux-amd64/captain" .
	GIN_MODE=release GOOS=linux GOARCH=arm64 go build -v -o "dist/captain-linux-arm64/captain" .
	GIN_MODE=release GOOS=darwin GOARCH=amd64 go build -v -o "dist/captain-darwin-amd64/captain" .
	GIN_MODE=release GOOS=darwin GOARCH=arm64 go build -v -o "dist/captain-darwin-arm64/captain" .
	GIN_MODE=release GOOS=windows GOARCH=amd64 go build -v -o "dist/captain-windows-amd64/captain.exe" .
	GIN_MODE=release GOOS=windows GOARCH=arm64 go build -v -o "dist/captain-windows-arm64/captain.exe" .
	cd dist && \
	for dir in */; do \
		if [ -d "$$dir" ] && [ "$$dir" != "zip/" ]; then \
			platform=$${dir%/}; \
			zip -r "zip/captain-$${platform##captain-}.zip" "$$platform"; \
		fi \
	done
	@echo "Release build complete. Binaries are in dist/ and zip archives are in dist/zip/"
