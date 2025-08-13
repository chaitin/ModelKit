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
- **可复用组件**: 提供开箱即用的ModelModal组件，支持在其他项目中快速集成

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
│   │   ├── components/     # 可复用组件
│   │   │   └── Card/       # 卡片组件
│   │   ├── constant/       # 常量定义
│   │   ├── api/            # API接口
│   │   └── services/       # 服务层
│   ├── public/       # 静态资源
│   └── package.json  # 前端依赖
├── utils/            # 工具函数
├── go.mod            # Go模块文件
└── Makefile          # 构建脚本
```

## 🛠️ 安装部署

### 环境要求

- Go 1.24.0+

## 📝 许可证

本项目采用 [MIT 许可证](LICENSE)。
