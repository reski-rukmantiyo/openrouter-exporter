BINARY_NAME := openrouter-exporter
BUILD_DIR := ./bin
GO := go
LDFLAGS := -s -w

.PHONY: all build clean $(BUILD_DIR)

all: build-all

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

build: $(BUILD_DIR)
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .

build-all: build-linux build-darwin build-windows

build-linux: build-linux-amd64 build-linux-arm64

build-linux-amd64: $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

build-linux-arm64: $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .

build-darwin: build-darwin-amd64 build-darwin-arm64

build-darwin-amd64: $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .

build-darwin-arm64: $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .

build-windows: build-windows-amd64 build-windows-arm64

build-windows-amd64: $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

build-windows-arm64: $(BUILD_DIR)
	GOOS=windows GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe .

clean:
	rm -rf $(BUILD_DIR)
