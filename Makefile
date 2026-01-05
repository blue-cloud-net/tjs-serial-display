# 项目配置
BINARY_NAME=tjc
OUTPUT_DIR=release
GO_FILES=$(shell find . -name '*.go' -type f)
MAIN_FILE=cmd/tjc.go

# Go 编译参数
GOFLAGS=-trimpath
LDFLAGS=-s -w

# 默认目标
.PHONY: all
all: build

# 编译项目
.PHONY: build
build: $(OUTPUT_DIR)/$(BINARY_NAME)

$(OUTPUT_DIR)/$(BINARY_NAME): $(GO_FILES)
	@mkdir -p $(OUTPUT_DIR)
	@echo "Building $(BINARY_NAME)..."
	go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(OUTPUT_DIR)/$(BINARY_NAME)"

# 编译（开发模式，不优化）
.PHONY: build-dev
build-dev:
	@mkdir -p $(OUTPUT_DIR)
	@echo "Building $(BINARY_NAME) (development mode)..."
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# 运行程序
.PHONY: run
run: build
	./$(OUTPUT_DIR)/$(BINARY_NAME)

# 清理构建产物
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(OUTPUT_DIR)
	@go clean
	@echo "Clean complete"

# 运行测试
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# 运行测试并生成覆盖率报告
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Format complete"

# 检查代码
.PHONY: vet
vet:
	@echo "Vetting code..."
	go vet ./...
	@echo "Vet complete"

# 下载依赖
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies updated"

# 安装到系统
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME)..."
	@install -m 755 $(OUTPUT_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "Installed to /usr/local/bin/$(BINARY_NAME)"

# 卸载
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstalled"

# 显示帮助信息
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all (default)    - Build the project"
	@echo "  build            - Build release binary"
	@echo "  build-dev        - Build development binary (no optimization)"
	@echo "  run              - Build and run the program"
	@echo "  clean            - Remove build artifacts"
	@echo "  test             - Run tests"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  fmt              - Format code"
	@echo "  vet              - Vet code"
	@echo "  deps             - Download and tidy dependencies"
	@echo "  install          - Install binary to /usr/local/bin"
	@echo "  uninstall        - Uninstall binary from /usr/local/bin"
	@echo "  help             - Show this help message"
