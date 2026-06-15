SHELL=/bin/bash

all: build test

bin:
	@mkdir -p bin

build: bin
	@go build -o bin/api cmd/api/main.go

run:
	@go run cmd/api/main.go

test:
	@echo "Testing..."
	@go test ./... -v

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -rf bin/

# Live Reload
watch:
	@bash scripts/watch.sh

migrate:
	@bash scripts/migrate.sh $(action)

.PHONY: all bin build run test clean watch docker-run docker-down itest