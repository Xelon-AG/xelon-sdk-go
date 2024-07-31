# Project variables
PROJECT_NAME := xelon-sdk-go

# Build variables
.DEFAULT_GOAL = test
BUILD_DIR := build
TOOLS_DIR := $(shell pwd)/tools
TOOLS_BIN_DIR := ${TOOLS_DIR}/bin


## tools: Install required tooling.
.PHONY: tools
tools:
	@echo "==> Installing required tooling..."
	@cd ${TOOLS_DIR} && GOBIN=${TOOLS_BIN_DIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint


## clean: Delete the build directory.
.PHONY: clean
clean:
	@echo "==> Removing '$(BUILD_DIR)' directory..."
	@rm -rf $(BUILD_DIR)


## lint: Lint code with golangci-lint.
.PHONY: lint
lint:
	@echo "==> Linting code with 'golangci-lint'..."
	@${TOOLS_BIN_DIR}/golangci-lint run


## test: Run all tests.
.PHONY: test
test:
	@echo "==> Running tests..."
	@mkdir -p $(BUILD_DIR)
	@go test -count=1 -v -cover -coverprofile=$(BUILD_DIR)/coverage.out ./...


help: Makefile
	@echo "Usage: make <command>"
	@echo ""
	@echo "Commands:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
