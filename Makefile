# 只需要改这2个参数！
APP_NAME := cfst

# 输出目录（编译后的文件直接放这，按平台命名）
OUTPUT_DIR:=./dest
WINDOWS_OUTPUT_DIR := $(OUTPUT_DIR)/windows
LINUX_OUTPUT_DIR := $(OUTPUT_DIR)/linux
DARWIN_OUTPUT_DIR := $(OUTPUT_DIR)/darwin

# 默认目标：执行make时，构建所有3个平台
all: build-windows build-linux build-darwin
	@echo "✅ 所有版本构建完成！文件在 $(OUTPUT_DIR) 目录"

# 构建Windows 64位版本
build-windows:
	rm -rf $(WINDOWS_OUTPUT_DIR)
	@mkdir -p $(WINDOWS_OUTPUT_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(WINDOWS_OUTPUT_DIR)/$(APP_NAME).exe
	cp config.yaml ip.txt ipv6.txt $(WINDOWS_OUTPUT_DIR)
	@echo "✅ Windows版本：$(WINDOWS_OUTPUT_DIR)/$(APP_NAME).exe"

# 构建Linux 64位版本
build-linux:
	rm -rf $(LINUX_OUTPUT_DIR)
	@mkdir -p $(LINUX_OUTPUT_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(LINUX_OUTPUT_DIR)/$(APP_NAME)
	cp config.yaml ip.txt ipv6.txt $(LINUX_OUTPUT_DIR)
	@echo "✅ Linux版本：$(LINUX_OUTPUT_DIR)/$(APP_NAME)"

# 构建MacOS 64位版本（支持Intel芯片，M1/M2也兼容）
build-darwin:
	rm -rf $(DARWIN_OUTPUT_DIR)
	@mkdir -p $(DARWIN_OUTPUT_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(DARWIN_OUTPUT_DIR)/$(APP_NAME)
	cp config.yaml ip.txt ipv6.txt $(DARWIN_OUTPUT_DIR)
	@echo "✅ MacOS版本：$(DARWIN_OUTPUT_DIR)/$(APP_NAME)"

# 清理构建产物
clean:
	rm -rf $(OUTPUT_DIR)
	@echo "✅ 清理完成！"

# 帮助信息
help:
	@echo "使用方法：make [目标]"
	@echo "目标："
	@echo "  all           构建Windows/Linux/MacOS所有版本（默认）"
	@echo "  build-windows 仅构建Windows版本"
	@echo "  build-linux   仅构建Linux版本"
	@echo "  build-darwin  仅构建MacOS版本"
	@echo "  clean         清理构建文件"
	@echo "  help          查看帮助"