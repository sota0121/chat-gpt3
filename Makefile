.PHONY: build clean run test

APP_NAME = "chatgpt3"

build:
	@echo "Building..."
	@go build -o bin/$(APP_NAME) cmd/main.go

clean:
	@echo "Cleaning..."
	@rm -rf bin

run:
	@echo "Running..."
	@go run cmd/main.go

test:
	@echo "Testing..."
	@go test -v ./...
