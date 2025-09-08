# ModelKit Backend Makefile

# 默认配置
IMAGE_NAME ?= modelkit-backend
IMAGE_TAG ?= latest
FULL_IMAGE_NAME = $(IMAGE_NAME):$(IMAGE_TAG)
TAR_FILE = $(IMAGE_NAME)-$(IMAGE_TAG).tar

# 构建Docker镜像并保存到当前文件夹
.PHONY: image
image:
	@echo "构建Docker镜像: $(FULL_IMAGE_NAME)"
	@cd build && ./build.sh -n $(IMAGE_NAME) -t $(IMAGE_TAG)
	@echo "保存镜像到文件: $(TAR_FILE)"
	@docker save $(FULL_IMAGE_NAME) -o $(TAR_FILE)
	@echo "✅ 镜像已保存到: $(PWD)/$(TAR_FILE)"
	@ls -lh $(TAR_FILE)

# 清理生成的镜像文件
.PHONY: clean
clean:
	@echo "清理镜像文件..."
	@rm -f *.tar
	@echo "✅ 清理完成"

# 加载镜像文件
.PHONY: load
load:
	@if [ -f "$(TAR_FILE)" ]; then \
		echo "加载镜像文件: $(TAR_FILE)"; \
		docker load -i $(TAR_FILE); \
		echo "✅ 镜像加载完成"; \
	else \
		echo "❌ 镜像文件不存在: $(TAR_FILE)"; \
		exit 1; \
	fi

# 显示帮助信息
.PHONY: help
help:
	@echo "ModelKit Backend Makefile"
	@echo ""
	@echo "可用目标:"
	@echo "  image    - 构建Docker镜像并保存到当前文件夹"
	@echo "  clean    - 清理生成的镜像文件"
	@echo "  load     - 加载镜像文件到Docker"
	@echo "  help     - 显示此帮助信息"
	@echo ""
	@echo "环境变量:"
	@echo "  IMAGE_NAME - 镜像名称 (默认: $(IMAGE_NAME))"
	@echo "  IMAGE_TAG  - 镜像标签 (默认: $(IMAGE_TAG))"
	@echo ""
	@echo "示例:"
	@echo "  make image                    # 使用默认配置"
	@echo "  make image IMAGE_TAG=v1.0.0   # 指定标签"
	@echo "  make clean                    # 清理文件"

# 默认目标
.DEFAULT_GOAL := help