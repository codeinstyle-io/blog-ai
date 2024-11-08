.PHONY: fmt run docker-build build clean docker-build

fmt:
	@echo "Formatting the project..."
	gofmt -s -w .

run:
	@echo "Running the project..."
	go run main.go

build: dist
	@echo "Building the project..."
	go build -o dist/main main.go

clean:
	@echo "Cleaning up..."
	rm -rf dist/

dist:
	mkdir -p dist

docker-build:
	@echo "Building Docker image..."
	docker build -t codeinstyle:latest .