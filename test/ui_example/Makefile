# 定义变量
IMAGE_NAME := modelmodal-example
IMAGE_TAG := latest
IMAGE_FULL := $(IMAGE_NAME):$(IMAGE_TAG)
TAR_FILE := $(IMAGE_NAME)-$(IMAGE_TAG).tar

# 默认目标
.PHONY: all
all: help

# 构建前端项目
.PHONY: build
build:
	pnpm build

# 构建Docker镜像并保存到当前文件夹
.PHONY: image
image: build
	@echo "构建Docker镜像: $(IMAGE_FULL)..."
	docker build -f ./build/Dockerfile -t $(IMAGE_FULL) .
	@echo "保存Docker镜像到文件: $(TAR_FILE)..."
	docker save -o $(TAR_FILE) $(IMAGE_FULL)
	@echo "完成! 镜像已保存到: $(TAR_FILE)"

# 清理生成的文件
.PHONY: clean
clean:
	rm -f $(TAR_FILE)

# 帮助信息
.PHONY: help
help:
	@echo "可用的命令:"
	@echo "  make build  - 构建前端项目"
	@echo "  make image  - 构建Docker镜像并保存到当前文件夹"
	@echo "  make clean  - 清理生成的文件"
	@echo "  make help   - 显示帮助信息"