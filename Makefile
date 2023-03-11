.PHONY: build clean run test

APP_NAME = "gochat"

build:
	@echo "Building..."
	@go build -o bin/$(APP_NAME) cmd/*.go

clean:
	@echo "Cleaning..."
	@rm -rf bin

run:
	@echo "Running..."
	@go run cmd/main.go

test:
	@echo "Testing..."
	@go test -v ./...
