VERSION ?= 0.3.0
BINARY_NAME = chronoflow
BUILD_DIR = dist

.PHONY: all
all: build

.PHONY: build
build:
	go build -o $(BINARY_NAME) ./cmd/chronoflow

.PHONY: install
install:
	go install ./cmd/chronoflow

.PHONY: run
run: build
	./$(BINARY_NAME)

# Cross-compile for releases
.PHONY: release
release: clean
	mkdir -p $(BUILD_DIR)
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/chronoflow
	# macOS AMD64 (Intel)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/chronoflow
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/chronoflow
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/chronoflow
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/chronoflow
	# Windows ARM64
	GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe ./cmd/chronoflow
	# Create tarballs
	cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)-$(VERSION)-windows-amd64.tar.gz $(BINARY_NAME)-windows-amd64.exe
	cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)-$(VERSION)-windows-arm64.tar.gz $(BINARY_NAME)-windows-arm64.exe
	# Generate checksums
	cd $(BUILD_DIR) && shasum -a 256 *.tar.gz > checksums.txt
	@echo "Release artifacts created in $(BUILD_DIR)/"

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
