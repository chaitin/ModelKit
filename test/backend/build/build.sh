#!/bin/bash

# ModelKit Backend Docker 构建脚本
# 用于构建 x86 架构的 Docker 镜像

set -e  # 遇到错误时退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认配置
IMAGE_NAME="modelkit-backend"
IMAGE_TAG="latest"
PLATFORM="linux/amd64"
DOCKERFILE_PATH="./Dockerfile"
CONTEXT_PATH="../../.."

# 显示帮助信息
show_help() {
    echo "ModelKit Backend Docker 构建脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -n, --name NAME        设置镜像名称 (默认: $IMAGE_NAME)"
    echo "  -t, --tag TAG          设置镜像标签 (默认: $IMAGE_TAG)"
    echo "  -p, --platform PLATFORM 设置目标平台 (默认: $PLATFORM)"
    echo "  -f, --file FILE        指定 Dockerfile 路径 (默认: $DOCKERFILE_PATH)"
    echo "  -c, --context PATH     设置构建上下文路径 (默认: $CONTEXT_PATH)"
    echo "  --no-cache             不使用缓存构建"
    echo "  --push                 构建后推送到仓库"
    echo "  -h, --help             显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0                                    # 使用默认配置构建"
    echo "  $0 -n my-app -t v1.0.0               # 指定名称和标签"
    echo "  $0 --no-cache                        # 不使用缓存构建"
    echo "  $0 -t v1.0.0 --push                  # 构建并推送"
}

# 解析命令行参数
NO_CACHE=""
PUSH=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--name)
            IMAGE_NAME="$2"
            shift 2
            ;;
        -t|--tag)
            IMAGE_TAG="$2"
            shift 2
            ;;
        -p|--platform)
            PLATFORM="$2"
            shift 2
            ;;
        -f|--file)
            DOCKERFILE_PATH="$2"
            shift 2
            ;;
        -c|--context)
            CONTEXT_PATH="$2"
            shift 2
            ;;
        --no-cache)
            NO_CACHE="--no-cache"
            shift
            ;;
        --push)
            PUSH=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}错误: 未知选项 $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# 构建完整的镜像名称
FULL_IMAGE_NAME="$IMAGE_NAME:$IMAGE_TAG"

echo -e "${BLUE}=== ModelKit Backend Docker 构建 ===${NC}"
echo -e "${YELLOW}镜像名称:${NC} $FULL_IMAGE_NAME"
echo -e "${YELLOW}目标平台:${NC} $PLATFORM"
echo -e "${YELLOW}Dockerfile:${NC} $DOCKERFILE_PATH"
echo -e "${YELLOW}构建上下文:${NC} $CONTEXT_PATH"
echo ""

# 检查 Docker 是否可用
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: Docker 未安装或不在 PATH 中${NC}"
    exit 1
fi

# 检查 Dockerfile 是否存在
if [[ ! -f "$DOCKERFILE_PATH" ]]; then
    echo -e "${RED}错误: Dockerfile 不存在: $DOCKERFILE_PATH${NC}"
    exit 1
fi

# 检查构建上下文是否存在
if [[ ! -d "$CONTEXT_PATH" ]]; then
    echo -e "${RED}错误: 构建上下文目录不存在: $CONTEXT_PATH${NC}"
    exit 1
fi

# 开始构建
echo -e "${GREEN}开始构建 Docker 镜像...${NC}"
echo ""

# 构建命令
BUILD_CMD="docker build \
    --platform $PLATFORM \
    --file $DOCKERFILE_PATH \
    --tag $FULL_IMAGE_NAME \
    $NO_CACHE \
    $CONTEXT_PATH"

echo -e "${YELLOW}执行命令:${NC}"
echo "$BUILD_CMD"
echo ""

# 执行构建
if eval "$BUILD_CMD"; then
    echo ""
    echo -e "${GREEN}✅ 镜像构建成功: $FULL_IMAGE_NAME${NC}"
    
    # 显示镜像信息
    echo ""
    echo -e "${BLUE}镜像信息:${NC}"
    docker images "$IMAGE_NAME" --format "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.CreatedAt}}\t{{.Size}}"
    
    # 推送镜像（如果指定）
    if [[ "$PUSH" == true ]]; then
        echo ""
        echo -e "${GREEN}推送镜像到仓库...${NC}"
        if docker push "$FULL_IMAGE_NAME"; then
            echo -e "${GREEN}✅ 镜像推送成功${NC}"
        else
            echo -e "${RED}❌ 镜像推送失败${NC}"
            exit 1
        fi
    fi
    
else
    echo ""
    echo -e "${RED}❌ 镜像构建失败${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}=== 构建完成 ===${NC}"
echo -e "${YELLOW}运行镜像:${NC} docker run --rm -p 8080:8080 $FULL_IMAGE_NAME"
echo -e "${YELLOW}查看日志:${NC} docker logs <container_id>"
echo -e "${YELLOW}进入容器:${NC} docker exec -it <container_id> /bin/sh"