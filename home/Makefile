
.PHONY: fmt run docker-build build

fmt:
	gofmt -s -w .

run:
	go run main.go

clean:
	rm -rf dist

dist:
	mkdir -p dist

build: dist
	go build -o dist/codeinstyle main.go

docker-build:
	docker build -t codeinstyle:latest .
