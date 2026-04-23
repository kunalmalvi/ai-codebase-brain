# Makefile for AI Coder Context Server

.PHONY: all build clean install test windows macos linux

# Version
VERSION ?= 1.0.0
all: windows macos linux

# Build for current platform
build:
	go build -ldflags="-s -w" -o ai-coder-context-server ./cmd/server

# Windows
windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/ai-coder-context-server-$(VERSION)-windows-amd64.exe ./cmd/server

# macOS
macos:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/ai-coder-context-server-$(VERSION)-darwin-amd64 ./cmd/server
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/ai-coder-context-server-$(VERSION)-darwin-arm64 ./cmd/server

# Linux
linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/ai-coder-context-server-$(VERSION)-linux-amd64 ./cmd/server
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/ai-coder-context-server-$(VERSION)-linux-arm64 ./cmd/server

# Clean build artifacts
clean:
	rm -rf dist/
	rm -f ai-coder-context-server*

# Install to GOPATH/bin
install: build
	go install ./cmd/server

# Run tests
test:
	go test ./...

# Create dist directory
setup:
	mkdir dist