BINARY=dapodik-bridge
BUILD_DIR=build
VERSION=1.1.0
LDFLAGS=-ldflags="-s -w -X 'github.com/ardianryan/dapodik-bridge/internal/config.AppVersion=$(VERSION)'"
WINDOWS_GUI_LDFLAGS=-ldflags="-H=windowsgui -s -w -X 'github.com/ardianryan/dapodik-bridge/internal/config.AppVersion=$(VERSION)'"

.PHONY: all build clean test run release build-all build-windows build-linux

all: test build

build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/bridge

run:
	go run ./cmd/bridge -port=4712

test:
	go test -v ./...

clean:
	rm -rf $(BUILD_DIR)

build-windows:
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build $(WINDOWS_GUI_LDFLAGS) -o $(BUILD_DIR)/DapodikBridge-gui-amd64.exe ./cmd/bridge
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe ./cmd/bridge
	GOOS=windows GOARCH=386 go build $(WINDOWS_GUI_LDFLAGS) -o $(BUILD_DIR)/DapodikBridge-gui-386.exe ./cmd/bridge
	GOOS=windows GOARCH=386 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-windows-386.exe ./cmd/bridge

build-linux:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-amd64 ./cmd/bridge
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-arm64 ./cmd/bridge

build-all: clean build-windows build-linux
	@echo "All Windows and Linux binaries successfully built in $(BUILD_DIR)/"
