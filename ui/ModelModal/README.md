# @modelkit/model-modal

一个功能完整的AI模型配置模态框组件，支持添加、编辑和测试各种AI模型。该组件可以在其他React项目中复用。

> **注意**: 这是一个私有包，仅供内部使用，不对外发布。

## ✨ 特性

- 🚀 **开箱即用** - 提供完整的模型配置功能
- 🎨 **高度可定制** - 支持主题、样式、功能开关等配置
- 🌍 **国际化支持** - 内置中英文支持，可扩展其他语言
- 🔌 **依赖注入** - 通过接口注入模型服务，易于集成
- 📱 **响应式设计** - 基于Material-UI，支持各种屏幕尺寸
- 🧪 **类型安全** - 完整的TypeScript类型定义
- 🎯 **功能开关** - 可配置启用/禁用特定功能

## 📦 安装

由于这是一个私有包，你需要通过以下方式使用：

### 方式1: 直接复制源码

将整个 `ModelModal` 目录复制到你的项目中：

```bash
cp -r ui/src/components/ModelModal your-project/src/components/
```

### 方式2: 使用相对路径导入

在你的项目中直接导入：

```tsx
import { ModelModal } from './components/ModelModal/src';
```

### 方式3: 构建后使用

先构建包，然后使用构建后的文件：

```bash
cd ui/src/components/ModelModal
npm install
npm run build
```

构建完成后，将 `dist` 目录复制到你的项目中使用。

## 🔧 依赖要求

```json
{
  "dependencies": {
    "react": ">=18.0.0",
    "react-dom": ">=18.0.0",
    "@mui/material": ">=5.0.0",
    "@emotion/react": ">=11.0.0",
    "@emotion/styled": ">=11.0.0",
    "react-hook-form": ">=7.0.0"
  }
}
```

## 🚀 快速开始

### 基本使用

```tsx
import React, { useState } from 'react';
import { ModelModal, createDefaultModelService } from './components/ModelModal/src';

function App() {
  const [open, setOpen] = useState(false);
  const [modelData, setModelData] = useState(null);
  
  // 创建默认服务实例
  const modelService = createDefaultModelService('https://api.example.com');
  
  const handleClose = () => {
    setOpen(false);
    setModelData(null);
  };
  
  const handleRefresh = () => {
    console.log('刷新模型列表');
  };
  
  return (
    <div>
      <button onClick={() => setOpen(true)}>添加模型</button>
      
      <ModelModal
        open={open}
        data={modelData}
        type="chat"
        onClose={handleClose}
        refresh={handleRefresh}
        modelService={modelService}
      />
    </div>
  );
}
```

### 使用预配置组件

```tsx
import { createModelModal } from './components/ModelModal/src';

// 创建预配置的组件
const CustomModelModal = createModelModal({
  theme: {
    primaryColor: '#ff6b6b',
    borderRadius: '16px',
  },
  locale: {
    language: 'en-US',
  },
  features: {
    enableModelTesting: false,
  },
});

// 使用预配置组件
<CustomModelModal
  open={open}
  type="chat"
  onClose={handleClose}
  refresh={handleRefresh}
  modelService={modelService}
/>
```

## ⚙️ 配置选项

### 主题配置

```tsx
const config = {
  theme: {
    primaryColor: '#1976d2',      // 主色调
    secondaryColor: '#dc004e',    // 次色调
    borderRadius: '10px',         // 圆角大小
    spacing: 2,                   // 间距倍数
  },
};
```

### 本地化配置

```tsx
const config = {
  locale: {
    language: 'zh-CN',            // 语言: 'zh-CN' | 'en-US'
    messages: {                   // 自定义消息
      'custom.message': '自定义消息',
    },
  },
};
```

### 验证配置

```tsx
const config = {
  validation: {
    requiredFields: ['provider', 'model', 'base_url'], // 必填字段
    customValidators: {           // 自定义验证器
      'api_key': (value) => {
        if (value.length < 10) return 'API密钥长度不足';
        return undefined;
      },
    },
  },
};
```

### 功能开关

```tsx
const config = {
  features: {
    enableModelTesting: true,     // 启用模型测试
    enableHeaderConfig: true,     // 启用Header配置
    enableApiVersion: true,       // 启用API版本
    enableProviderSelection: true, // 启用提供商选择
  },
};
```

### 样式配置

```tsx
const config = {
  styles: {
    modalWidth: 800,              // 模态框宽度
    sidebarWidth: 200,            // 侧边栏宽度
    customCSS: '.custom-class { color: red; }', // 自定义CSS
  },
};
```

## 🔌 服务接口

组件通过 `ModelService` 接口与后端交互，你需要实现以下方法：

```tsx
interface ModelService {
  createModel: (data: CreateModelData) => Promise<{ ModelName: string }>;
  getModelNameList: (data: GetModelNameData) => Promise<{ models: { model: string }[] }>;
  testModel: (data: CheckModelData) => Promise<{ error: string }>;
  updateModel: (data: UpdateModelData) => Promise<void>;
}
```

### 使用默认服务实现

```tsx
import { createDefaultModelService } from './components/ModelModal/src';

const modelService = createDefaultModelService('https://api.example.com', {
  'Authorization': 'Bearer your-token',
});

// 设置认证头
modelService.setAuthToken('your-token');

// 设置基础URL
modelService.setBaseURL('https://new-api.example.com');
```

### 自定义服务实现

```tsx
import { ModelService } from './components/ModelModal/src';

class CustomModelService implements ModelService {
  async createModel(data) {
    // 自定义实现
    const response = await fetch('/api/models', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return response.json();
  }
  
  // 实现其他方法...
}

const modelService = new CustomModelService();
```

## 🎨 自定义提供商

你可以自定义模型提供商配置：

```tsx
const customProviders = {
  CustomProvider: {
    label: 'CustomProvider',
    cn: '自定义提供商',
    icon: 'icon-custom',
    urlWrite: true,
    secretRequired: true,
    customHeader: false,
    modelDocumentUrl: 'https://docs.example.com',
    defaultBaseUrl: 'https://api.example.com/v1',
  },
};

<ModelModal
  customProviders={customProviders}
  // ... 其他属性
/>
```

## 🔄 事件回调

组件提供多个回调函数：

```tsx
<ModelModal
  onBeforeSubmit={async (data) => {
    // 提交前的验证
    if (data.api_key.length < 10) {
      alert('API密钥长度不足');
      return false;
    }
    return true;
  }}
  onAfterSubmit={(data, result) => {
    // 提交后的处理
    console.log('模型创建成功:', result);
  }}
  onError={(error) => {
    // 错误处理
    console.error('操作失败:', error);
  }}
  // ... 其他属性
/>
```

## 📱 响应式设计

组件基于Material-UI构建，支持响应式设计：

```tsx
const config = {
  styles: {
    modalWidth: { xs: '95%', sm: 600, md: 800, lg: 1000 },
    sidebarWidth: { xs: 150, sm: 180, md: 200 },
  },
};
```

## 🧪 开发和构建

### 开发模式

```bash
# 进入组件目录
cd ui/src/components/ModelModal

# 安装依赖
npm install

# 开发模式（监听文件变化）
npm run dev

# 类型检查
npm run type-check

# 代码检查
npm run lint
```

### 构建

```bash
# 构建生产版本
npm run build

# 清理构建文件
npm run clean
```

构建完成后，`dist` 目录包含以下文件：
- `index.js` - CommonJS 格式
- `index.es.js` - ES Module 格式
- `index.d.ts` - TypeScript 类型定义

## 📚 API 参考

### ModelModal Props

| 属性 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `open` | `boolean` | ✅ | - | 控制模态框显示/隐藏 |
| `data` | `ModelListItem \| null` | ❌ | `null` | 编辑时的模型数据 |
| `type` | `'chat' \| 'embedding' \| 'rerank'` | ✅ | - | 模型类型 |
| `onClose` | `() => void` | ✅ | - | 关闭回调 |
| `refresh` | `() => void` | ✅ | - | 刷新回调 |
| `modelService` | `ModelService` | ✅ | - | 模型服务 |
| `config` | `ModelModalConfig` | ❌ | `DEFAULT_CONFIG` | 配置选项 |
| `customProviders` | `ModelProviderMap` | ❌ | `DEFAULT_MODEL_PROVIDERS` | 自定义提供商 |
| `onBeforeSubmit` | `(data: AddModelForm) => boolean \| Promise<boolean>` | ❌ | - | 提交前回调 |
| `onAfterSubmit` | `(data: AddModelForm, result: any) => void` | ❌ | - | 提交后回调 |
| `onError` | `(error: string) => void` | ❌ | - | 错误回调 |

### 类型定义

```tsx
// 模型类型
type ModelType = 'chat' | 'embedding' | 'rerank';

// 模型提供商配置
interface ModelProviderConfig {
  label: string;
  cn?: string;
  icon: string;
  urlWrite: boolean;
  secretRequired: boolean;
  customHeader: boolean;
  modelDocumentUrl?: string;
  defaultBaseUrl: string;
}

// 表单数据
interface AddModelForm {
  provider: string;
  model: string;
  base_url: string;
  api_version: string;
  api_key: string;
  api_header_key: string;
  api_header_value: string;
  type: ModelType;
}
```

## 📁 项目结构

```
ModelModal/
├── src/                    # 源代码
│   ├── components/         # 组件
│   ├── constants/          # 常量定义
│   ├── services/           # 服务实现
│   ├── types/              # 类型定义
│   ├── utils/              # 工具函数
│   └── index.ts            # 主入口
├── examples/               # 使用示例
│   ├── basic-usage.tsx     # 基本使用示例
│   ├── advanced-config.tsx # 高级配置示例
│   └── index.tsx           # 示例入口
├── dist/                   # 构建输出（构建后生成）
├── package.json            # 包配置
├── tsconfig.json           # TypeScript配置
├── rollup.config.js        # 构建配置
├── README.md               # 说明文档
├── build.sh                # 构建脚本
└── publish.sh              # 发布脚本（私有包不需要）
```

## 🔗 集成到其他项目

### 步骤1: 复制组件

将整个 `ModelModal` 目录复制到你的项目中。

### 步骤2: 安装依赖

确保你的项目安装了必要的依赖：

```bash
npm install @mui/material @emotion/react @emotion/styled react-hook-form
```

### 步骤3: 导入使用

```tsx
import { ModelModal } from './components/ModelModal/src';
```

### 步骤4: 配置服务

实现 `ModelService` 接口或使用默认实现。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 🔗 相关链接

- [ModelKit 项目](https://github.com/chaitin/ModelKit)
- [Material-UI](https://mui.com/)
- [React Hook Form](https://react-hook-form.com/) 