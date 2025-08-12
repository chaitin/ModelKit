#!/bin/bash

echo "🚀 开始构建 ModelModal 私有包..."

# 检查依赖
if ! command -v npm &> /dev/null; then
    echo "❌ 错误: 未找到 npm"
    exit 1
fi

# 清理之前的构建
echo "🧹 清理之前的构建..."
rm -rf dist

# 安装依赖
echo "📦 安装依赖..."
npm install

# 类型检查
echo "🔍 运行类型检查..."
npm run type-check

if [ $? -ne 0 ]; then
    echo "❌ 类型检查失败"
    exit 1
fi

# 构建
echo "🔨 构建包..."
npm run build

if [ $? -ne 0 ]; then
    echo "❌ 构建失败"
    exit 1
fi

# 检查构建结果
if [ ! -d "dist" ]; then
    echo "❌ 构建目录不存在"
    exit 1
fi

echo "✅ 构建成功！"
echo "📁 构建结果:"
ls -la dist/

echo ""
echo "🎉 ModelModal 私有包构建完成！"
echo "📦 包位置: $(pwd)/dist/"
echo "📚 现在可以将 dist/ 目录复制到其他项目中使用"
echo ""
echo "💡 使用方法:"
echo "1. 将 dist/ 目录复制到目标项目"
echo "2. 在目标项目中导入: import { ModelModal } from './dist'"
echo "3. 确保目标项目安装了必要的依赖" 