.PHONY: build test test-coverage clean run docker-build docker-run help

# 默认目标
.DEFAULT_GOAL := help

# 应用名称
APP_NAME := go-ocr
VERSION := 1.0
DOCKER_IMAGE := xsdhy/go-ocr:$(VERSION)

# Go相关变量
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOCLEAN := $(GOCMD) clean
GOMOD := $(GOCMD) mod

# 构建目录
BUILD_DIR := build
BINARY_NAME := $(APP_NAME)

## build: 构建应用程序
build:
	@echo "构建应用程序..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v .
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

## test: 运行单元测试
test:
	@echo "运行单元测试..."
	$(GOTEST) ./... -tags=test -v

## test-coverage: 运行测试并生成覆盖率报告
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	$(GOTEST) ./... -tags=test -coverprofile=coverage.out
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

## test-short: 运行快速测试
test-short:
	@echo "运行快速测试..."
	$(GOTEST) ./... -tags=test -short

## clean: 清理构建文件
clean:
	@echo "清理构建文件..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -rf tmp/
	@echo "清理完成"

## run: 运行应用程序
run:
	@echo "运行应用程序..."
	$(GOCMD) run . 

## run-dev: 以开发模式运行应用程序
run-dev:
	@echo "以开发模式运行应用程序..."
	GIN_MODE=debug $(GOCMD) run .

## deps: 下载依赖
deps:
	@echo "下载依赖..."
	$(GOMOD) download
	$(GOMOD) tidy

## fmt: 格式化代码
fmt:
	@echo "格式化代码..."
	$(GOCMD) fmt ./...

## vet: 运行go vet
vet:
	@echo "运行go vet..."
	$(GOCMD) vet ./...

## lint: 运行golangci-lint (需要先安装)
lint:
	@echo "运行golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint未安装，跳过lint检查"; \
	fi

## docker-build: 构建Docker镜像
docker-build:
	@echo "构建Docker镜像..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "Docker镜像构建完成: $(DOCKER_IMAGE)"

## docker-run: 运行Docker容器
docker-run:
	@echo "运行Docker容器..."
	docker run --name $(APP_NAME) --rm -d -p 8080:8080 $(DOCKER_IMAGE)
	@echo "Docker容器已启动，访问 http://localhost:8080"

## docker-stop: 停止Docker容器
docker-stop:
	@echo "停止Docker容器..."
	docker stop $(APP_NAME) || true

## check: 运行所有检查（格式化、vet、测试）
check: fmt vet test
	@echo "所有检查完成"

## install-tools: 安装开发工具
install-tools:
	@echo "安装开发工具..."
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

## help: 显示帮助信息
help:
	@echo "可用的命令:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## all: 运行完整的构建流程
all: clean deps fmt vet test build
	@echo "完整构建流程完成" 