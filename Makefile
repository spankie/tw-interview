include .env

# Build the application
all: build

build:
	@echo "Building..."
	@go build -o main cmd/web/*.go

# Run the application
run:
	@go run cmd/web/*.go

lint:
	@echo "Linting..."
	@golangci-lint run

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v


# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test -tags="integration" -v ./...


# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

.PHONY: all build run lint test clean 

