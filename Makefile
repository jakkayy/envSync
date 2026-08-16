.PHONY: all build test lint clean run-cli run-server

BINARY_CLI=bin/envsync
BINARY_SERVER=bin/server

all: lint test build

build:
	@echo "==> Building binaries..."
	@mkdir -p bin
	go build -o $(BINARY_CLI) ./cmd/envsync
	go build -o $(BINARY_SERVER) ./cmd/server

test:
	@echo "==> Running unit tests..."
	go test -v -race -cover ./...

lint:
	@echo "==> Running golangci-lint..."
	golangci-lint run ./...

clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf bin/ dist/ coverage.out

run-cli: build
	./$(BINARY_CLI)

run-server: build
	./$(BINARY_SERVER)
