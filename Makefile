# Makefile for cross-compiling agent and updater

BINARY_NAME_AGENT=demo
# Output to bin/ directory in root (absolute path to avoid issues with cd)
ROOT_DIR:=$(shell pwd)
BUILD_DIR:=$(ROOT_DIR)/bin

.PHONY: all clean windows linux mac

all: windows linux mac

windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME_AGENT)-windows.exe .
	@if [ -f agent/config.json ]; then cp agent/config.json $(BUILD_DIR)/config.json; fi

linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME_AGENT)-linux .
	@if [ -f agent/config.json ]; then cp agent/config.json $(BUILD_DIR)/config.json; fi

mac:
	@echo "Building for macOS (Darwin)..."
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME_AGENT)-darwin .
	@if [ -f agent/config.json ]; then cp agent/config.json $(BUILD_DIR)/config.json; fi

clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)
