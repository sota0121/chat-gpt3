.PHONY: setup build clean run test xbuild

APP_NAME = "gochat"
VERSION = "0.0.1"
GOBIN ?= $(shell go env GOPATH)/bin
ASSET_DIR = "dist"

setup: $(GOBIN)/goxz $(GOBIN)/ghr $(GOBIN)/gobump
	@echo "Setting up..."

build:
	@echo "Building..."
	@go build -o bin/$(APP_NAME) cmd/*.go

xbuild: $(GOBIN)/goxz
	@echo "Cross building..."
	@goxz -pv=v$(VERSION) -d=$(ASSET_DIR) -os="linux darwin windows" -arch="amd64 arm64" -build-ldflags="-s -w" -build-tags="release" -n=$(APP_NAME) cmd/*.go

upload-release-assets: $(GOBIN)/ghr
	@echo "Uploading assets..."
	@ghr "v$(VERSION)" $(ASSET_DIR)

clean:
	@echo "Cleaning..."
	@rm -rf bin

run:
	@echo "Running..."
	@go run cmd/main.go

test:
	@echo "Testing..."
	@go test -v ./...

$(GOBIN)/goxz:
	@go install github.com/Songmu/goxz/cmd/goxz@latest

$(GOBIN)/ghr:
	@go install github.com/tcnksm/ghr@latest

$(GOBIN)/gobump:
	@go install github.com/x-motemen/gobump/cmd/gobump@master