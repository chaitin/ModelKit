# ModelKit

[![Go Version](https://img.shields.io/badge/Go-1.24.0+-blue.svg)](https://golang.org)
[![React Version](https://img.shields.io/badge/React-19.0.0+-blue.svg)](https://reactjs.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

ModelKit 是一个强大的AI模型管理平台，支持多种AI服务提供商，提供统一的模型管理、配置验证和API接口服务。

## 🚀 功能特性

- **多模型提供商支持**: 支持 OpenAI、Ollama、DeepSeek、SiliconFlow、Moonshot、Azure OpenAI、百智云、腾讯混元、百炼、火山引擎、Gemini、智谱等主流AI服务商
- **模型类型管理**: 支持聊天模型、嵌入模型、重排序模型、视觉模型、代码模型、函数调用等多种模型类型
- **配置验证**: 提供模型配置的实时验证功能，确保API配置正确性
- **统一API接口**: 提供标准化的RESTful API，简化AI模型集成
- **现代化Web界面**: 基于React 19和Material-UI构建的响应式用户界面
- **国际化支持**: 内置中英文多语言支持

## 🏗️ 技术架构

### 后端 (Go)
- **框架**: Echo v4 (HTTP框架)
- **语言**: Go 1.24.0+
- **架构**: 分层架构 (Handler -> UseCase -> Domain)
- **依赖管理**: Go Modules

### 前端 (React)
- **框架**: React 19 + TypeScript
- **UI库**: Material-UI v6 + CT-MUI
- **构建工具**: Vite 6
- **状态管理**: Redux Toolkit
- **路由**: React Router v7
- **编辑器**: TipTap (富文本编辑器)

## 📁 项目结构

```
ModelKit/
├── consts/           # 常量定义
├── domain/           # 领域模型和接口
├── errcode/          # 错误码和国际化
├── handler/          # HTTP处理器
├── pkg/              # 公共包
├── usecase/          # 业务用例
├── ui/               # 前端应用
│   ├── src/          # 源代码
│   ├── public/       # 静态资源
│   └── package.json  # 前端依赖
├── utils/            # 工具函数
├── go.mod            # Go模块文件
└── Makefile          # 构建脚本
```

## 🛠️ 安装部署

### 环境要求

- Go 1.24.0+
- Node.js 18+
- pnpm 10.12.1+

### 前端部署(开发中)

1. 进入前端目录
```bash
cd ui
```

2. 安装依赖
```bash
pnpm install
```

3. 开发模式运行
```bash
pnpm dev
```

4. 构建生产版本
```bash
pnpm build
```

## 📚 API文档

### 获取模型列表

```http
GET /api/v1/modelkit/models
```

**查询参数:**
- `provider`: 模型提供商 (必需)
- `base_url`: 基础URL (必需)
- `api_key`: API密钥 (可选)
- `api_header`: API头部 (可选)
- `type`: 模型类型 (必需)

**响应示例:**
```json
{
  "success": true,
  "message": "获取模型列表成功",
  "data": {
    "models": [
      {
        "id": "gpt-4",
        "object": "model",
        "owned_by": "OpenAI",
        "model_type": "chat"
      }
    ]
  }
}
```

### 检查模型配置

```http
POST /api/v1/modelkit/check
```

**请求体:**
```json
{
  "provider": "OpenAI",
  "model": "gpt-4",
  "base_url": "https://api.openai.com",
  "api_key": "sk-xxx",
  "type": "chat"
}
```

**响应示例:**
```json
{
  "success": true,
  "message": "模型配置验证成功",
  "data": {
    "error": "",
    "content": "配置验证通过"
  }
}
```

## 🔧 配置说明

### 支持的模型提供商

- **OpenAI**: GPT系列模型
- **Ollama**: 本地部署模型
- **DeepSeek**: DeepSeek系列模型
- **SiliconFlow**: SiliconFlow模型
- **Moonshot**: Moonshot系列模型
- **Azure OpenAI**: Azure OpenAI服务
- **百智云**: 通义千问等模型
- **腾讯混元**: 混元系列模型
- **百炼**: 百炼系列模型
- **火山引擎**: 火山引擎模型
- **Gemini**: Google Gemini模型
- **智谱**: 智谱AI模型

### 支持的模型类型

- **chat**: 聊天模型
- **embedding**: 嵌入模型
- **rerank**: 重排序模型
- **vision**: 视觉模型
- **coder**: 代码模型
- **functioncall**: 函数调用模型


## 📝 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📞 联系我们

- 项目主页: [https://github.com/chaitin/ModelKit](https://github.com/chaitin/ModelKit)
- 问题反馈: [Issues](https://github.com/chaitin/ModelKit/issues)

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者和用户！
