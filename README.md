# ModelKit

[![Go Version](https://img.shields.io/badge/Go-1.24.0+-blue.svg)](https://golang.org)
[![React Version](https://img.shields.io/badge/React-19.0.0+-blue.svg)](https://reactjs.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

ModelKit 是一个强大的AI模型管理平台，支持多种AI服务提供商，提供统一的模型管理、配置验证服务。

## 🚀 功能特性

- **多模型提供商支持**: 支持 OpenAI、Ollama、DeepSeek、SiliconFlow、Moonshot、Azure OpenAI、百智云、腾讯混元、百炼、火山引擎、Gemini、智谱等主流AI服务商
- **模型类型管理**: 支持聊天模型、嵌入模型、重排序模型、视觉模型、代码模型、函数调用等多种模型类型
- **配置验证**: 提供模型配置的实时验证功能，确保API配置正确性
- **现代化Web界面**: 基于React 19和Material-UI构建的响应式用户界面
- **国际化支持**: 内置中英文多语言支持
- **可复用组件**: 提供开箱即用的ModelModal组件，支持在其他项目中快速集成

## 在项目中集成ModelKit

### 1. 安装依赖

#### 后端依赖
```bash
go get github.com/chaitin/ModelKit
go get github.com/labstack/echo/v4
go get github.com/go-playground/validator/v10
```

#### 前端依赖
```bash
npm install @yokowu/modelkit-ui
# 或
yarn add @yokowu/modelkit-ui
```

### 2. 实现接口

需要实现以下4个接口，其中 `listModel` 和 `checkModel` 已提供业务逻辑，在handler中调用即可：

#### ListModel 接口
- **请求参数** (`domain.ModelListReq`):
  ```go
  type ModelListReq struct {
      Provider  string `json:"provider" validate:"required,oneof=SiliconFlow OpenAI Ollama DeepSeek Moonshot AzureOpenAI BaiZhiCloud Hunyuan BaiLian Volcengine Gemini ZhiPu"`
      BaseURL   string `json:"base_url" validate:"required"`
      APIKey    string `json:"api_key"`
      APIHeader string `json:"api_header"`
      Type      string `json:"type" validate:"required,oneof=chat embedding rerank"`
  }
  ```
- **响应参数** (`domain.ModelListResp`):
  ```go
  type ModelListResp struct {
      Models []ModelListItem `json:"models"`
  }
  type ModelListItem struct {
      Model string `json:"model"`
  }
  ```

#### CheckModel 接口
- **请求参数** (`domain.CheckModelReq`):
  ```go
  type CheckModelReq struct {
      Provider   string `json:"provider" validate:"required,oneof=OpenAI Ollama DeepSeek SiliconFlow Moonshot Other AzureOpenAI BaiZhiCloud Hunyuan BaiLian Volcengine Gemini ZhiPu"`
      Model      string `json:"model" validate:"required"`
      BaseURL    string `json:"base_url" validate:"required"`
      APIKey     string `json:"api_key"`
      APIHeader  string `json:"api_header"`
      APIVersion string `json:"api_version"` // for azure openai
      Type       string `json:"type" validate:"required,oneof=chat embedding rerank"`
  }
  ```
- **响应参数** (`domain.CheckModelResp`):
  ```go
  type CheckModelResp struct {
      Error   string `json:"error"`
      Content string `json:"content"`
  }
  ```

#### CreateModel 接口
- **请求参数** (`CreateModelReq`):
  ```go
  type CreateModelReq struct {
      APIBase    string        `json:"api_base" validate:"required"`
      APIHeader  string        `json:"api_header"`
      APIKey     string        `json:"api_key"`
      APIVersion string        `json:"api_version"`
      ModelName  string        `json:"model_name" validate:"required"`
      ModelType  string        `json:"model_type" validate:"oneof=llm coder embedding audio reranker"`
      Param      *ModelParam   `json:"param"`
      Provider   string        `json:"provider" validate:"required,oneof=SiliconFlow OpenAI Ollama DeepSeek Moonshot AzureOpenAI BaiZhiCloud Hunyuan BaiLian Volcengine Other"`
      ShowName   string        `json:"show_name"`
  }
  
  type ModelParam struct {
      ContextWindow      int  `json:"context_window"`
      MaxTokens          int  `json:"max_tokens"`
      R1Enabled          bool `json:"r1_enabled"`
      SupportComputerUse bool `json:"support_computer_use"`
      SupportImages      bool `json:"support_images"`
      SupportPromptCache bool `json:"support_prompt_cache"`
  }
  ```
- **响应参数** (`CreateModelResp`):
  ```go
  type CreateModelResp struct {
      Model Model `json:"model"`
  }
  ```

#### UpdateModel 接口
- **请求参数** (`UpdateModelReq`):
  ```go
  type UpdateModelReq struct {
      ID         string        `json:"id" validate:"required"`
      APIBase    string        `json:"api_base"`
      APIHeader  string        `json:"api_header"`
      APIKey     string        `json:"api_key"`
      APIVersion string        `json:"api_version"`
      ModelName  string        `json:"model_name"`
      Param      *ModelParam   `json:"param"`
      Provider   string        `json:"provider" validate:"oneof=SiliconFlow OpenAI Ollama DeepSeek Moonshot AzureOpenAI BaiZhiCloud Hunyuan BaiLian Volcengine Other"`
      ShowName   string        `json:"show_name"`
      Status     string        `json:"status" validate:"oneof=active inactive"`
  }
  ```
- **响应参数** (`UpdateModelResp`):
  ```go
  type UpdateModelResp struct {
      Model Model `json:"model"`
  }
  ```

#### 通用Model结构
```go
type Model struct {
    ID          string        `json:"id"`
    APIBase     string        `json:"api_base"`
    APIHeader   string        `json:"api_header"`
    APIKey      string        `json:"api_key"`
    APIVersion  string        `json:"api_version"`
    CreatedAt   int64         `json:"created_at"`
    Input       int           `json:"input"`
    IsActive    bool          `json:"is_active"`
    IsInternal  bool          `json:"is_internal"`
    ModelName   string        `json:"model_name"`
    ModelType   string        `json:"model_type"`
    Output      int           `json:"output"`
    Param       *ModelParam   `json:"param"`
    Provider    string        `json:"provider"`
    ShowName    string        `json:"show_name"`
    Status      string        `json:"status"`
    UpdatedAt   int64         `json:"updated_at"`
}
```

### 3. 后端使用方式

在handler中调用 `listModel` 与 `checkModel` 业务逻辑：

```go
package v1

import (
    "net/http"
    "github.com/chaitin/ModelKit/domain"
    "github.com/chaitin/ModelKit/usecase"
    "github.com/labstack/echo/v4"
)

type ModelKit struct{}

// GetModelList 获取模型列表
func (m *ModelKit) GetModelList(c echo.Context) error {
    var req domain.ModelListReq
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, domain.Response{
            Success: false,
            Message: "参数绑定失败: " + err.Error(),
        })
    }

    // 调用ModelKit提供的业务逻辑
    resp, err := usecase.ModelList(c.Request().Context(), &req)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, domain.Response{
            Success: false,
            Message: err.Error(),
        })
    }

    return c.JSON(http.StatusOK, domain.Response{
        Success: true,
        Message: "获取模型列表成功",
        Data:    resp,
    })
}

// CheckModel 检查模型
func (m *ModelKit) CheckModel(c echo.Context) error {
    var req domain.CheckModelReq
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, domain.Response{
            Success: false,
            Message: "参数绑定失败: " + err.Error(),
        })
    }

    // 调用ModelKit提供的业务逻辑
    resp, err := usecase.CheckModel(c.Request().Context(), &req)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, domain.Response{
            Success: false,
            Message: err.Error(),
        })
    }

    if resp.Error != "" {
        return c.JSON(http.StatusBadRequest, domain.Response{
            Success: false,
            Message: "模型检查失败",
            Data:    resp,
        })
    }

    return c.JSON(http.StatusOK, domain.Response{
        Success: true,
        Message: "模型检查成功",
        Data:    resp,
    })
}

// CreateModel 创建模型 (需要自行实现业务逻辑)
func (m *ModelKit) CreateModel(c echo.Context) error {
    // 实现创建模型的业务逻辑
    // ...
}

// UpdateModel 更新模型 (需要自行实现业务逻辑)
func (m *ModelKit) UpdateModel(c echo.Context) error {
    // 实现更新模型的业务逻辑
    // ...
}
```

### 4. 前端使用方式

```typescript
// 1. 引入ModelKit组件和类型
import { 
    ModelModal, 
    Model, 
    ModelService, 
    ConstsModelType as ModelKitType, 
    ModelListItem 
} from '@yokowu/modelkit-ui';

// 2. 创建符合ModelService接口的服务实现
const modelService: ModelService = {
    createModel: async (params) => {
        const response = await postCreateModel(params as unknown as DomainCreateModelReq);
        return { model: response as unknown as Model };
    },
    listModel: async (params) => {
        const response = await getGetProviderModelList(params as unknown as GetGetProviderModelListParams);
        return { models: response?.models || [] };
    },
    checkModel: async (params) => {
        const response = await postCheckModel(params as unknown as DomainCheckModelReq);
        return { model: response as unknown as Model };
    },
    updateModel: async (params) => {
        const response = await putUpdateModel(params as unknown as DomainUpdateModelReq);
        return { model: response as unknown as Model };
    }
};

// 3. 使用ModelModal组件
function App() {
    const [open, setOpen] = useState(false);
    const [editData, setEditData] = useState<Model | null>(null);
    const [modelType, setModelType] = useState<ModelKitType>(ModelKitType.CHAT);

    const refreshModel = () => {
        // 刷新模型列表的逻辑
    };

    return (
        <ModelModal
            open={open}
            onClose={() => {
                setOpen(false);
                setEditData(null);
            }}
            refresh={refreshModel}
            data={editData}
            type={modelType}
            modelService={modelService}
            language="zh-CN"
        />
    );
}
```

## 🤝 贡献指南

我们欢迎所有形式的贡献！无论是bug修复、新功能开发、文档改进还是问题反馈，都对项目的发展非常有价值。

### 如何贡献

1. **Fork 项目**
   ```bash
   git clone https://github.com/chaitin/ModelKit.git
   cd ModelKit
   ```

2. **创建功能分支**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **进行开发**
   - 遵循现有的代码风格和约定

4. **提交更改**
   ```bash
   git add .
   git commit -m "feat: 添加新功能描述"
   ```

5. **推送到远程仓库**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **创建 Pull Request**
   - 提供清晰的PR描述
   - 说明更改的目的和影响
   - 关联相关的Issue（如果有）

### 代码规范

- **Go代码**: 遵循 `gofmt` 和 `golint` 标准
- **TypeScript/React代码**: 遵循 ESLint 和 Prettier 配置
- **提交信息**: 使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式
- **测试**: 新功能必须包含相应的单元测试

### 开发环境设置

1. **后端开发**
   ```bash
   go mod tidy
   go run main.go
   ```

2. **前端开发**
   ```bash
   cd ui/ModelModal
   npm install
   npm run dev
   ```

## 📄 许可证

本项目采用 [MIT License](LICENSE) 开源协议。

## 🙏 致谢

感谢所有为ModelKit项目做出贡献的开发者和用户！

---

**ModelKit** - 让AI模型管理更简单 🚀