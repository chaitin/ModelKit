# ModelModal 集成指南

本指南将帮助你如何将这个私有包集成到其他React项目中。

## 🎯 集成方式

### 方式1: 源码集成（推荐）

直接将源码复制到你的项目中，这样可以完全控制代码，便于调试和定制。

```bash
# 复制整个目录
cp -r ui/src/components/ModelModal your-project/src/components/

# 或者只复制需要的文件
mkdir -p your-project/src/components/ModelModal
cp -r ui/src/components/ModelModal/src your-project/src/components/ModelModal/
```

### 方式2: 构建后集成

先构建包，然后使用构建后的文件。

```bash
# 在ModelModal目录下构建
cd ui/src/components/ModelModal
npm install
npm run build

# 将构建结果复制到目标项目
cp -r dist your-project/src/components/ModelModal/
```

## 📦 依赖安装

确保你的项目安装了必要的依赖：

```bash
npm install @mui/material @emotion/react @emotion/styled react-hook-form
# 或者使用yarn
yarn add @mui/material @emotion/react @emotion/styled react-hook-form
# 或者使用pnpm
pnpm add @mui/material @emotion/react @emotion/styled react-hook-form
```

## 🔧 基本集成步骤

### 步骤1: 复制组件

选择上述任一方式将组件复制到你的项目中。

### 步骤2: 导入组件

```tsx
// 如果使用源码集成
import { ModelModal, createDefaultModelService } from './components/ModelModal/src';

// 如果使用构建后集成
import { ModelModal, createDefaultModelService } from './components/ModelModal/dist';
```

### 步骤3: 实现服务接口

```tsx
import { ModelService } from './components/ModelModal/src';

// 方式1: 使用默认实现
const modelService = createDefaultModelService('https://api.example.com');

// 方式2: 自定义实现
class CustomModelService implements ModelService {
  async createModel(data) {
    const response = await fetch('/api/models', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    return response.json();
  }
  
  async getModelNameList(data) {
    const response = await fetch('/api/models/list', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    return response.json();
  }
  
  async testModel(data) {
    const response = await fetch('/api/models/test', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    return response.json();
  }
  
  async updateModel(data) {
    await fetch(`/api/models/${data.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
  }
}

const modelService = new CustomModelService();
```

### 步骤4: 使用组件

```tsx
import React, { useState } from 'react';
import { ModelModal } from './components/ModelModal/src';

function App() {
  const [open, setOpen] = useState(false);
  
  return (
    <div>
      <button onClick={() => setOpen(true)}>添加模型</button>
      
      <ModelModal
        open={open}
        type="chat"
        onClose={() => setOpen(false)}
        refresh={() => console.log('刷新')}
        modelService={modelService}
      />
    </div>
  );
}
```

## 🎨 自定义配置

### 主题配置

```tsx
<ModelModal
  config={{
    theme: {
      primaryColor: '#ff6b6b',
      borderRadius: '16px',
      spacing: 3,
    },
  }}
  // ... 其他属性
/>
```

### 本地化配置

```tsx
<ModelModal
  config={{
    locale: {
      language: 'en-US',
      messages: {
        'custom.message': 'Custom Message',
      },
    },
  }}
  // ... 其他属性
/>
```

### 功能开关

```tsx
<ModelModal
  config={{
    features: {
      enableModelTesting: false,
      enableHeaderConfig: true,
      enableApiVersion: false,
      enableProviderSelection: true,
    },
  }}
  // ... 其他属性
/>
```

## 🔌 自定义提供商

```tsx
const customProviders = {
  MyAI: {
    label: 'MyAI',
    cn: '我的AI',
    icon: 'icon-myai',
    urlWrite: true,
    secretRequired: true,
    customHeader: false,
    modelDocumentUrl: 'https://docs.myai.com',
    defaultBaseUrl: 'https://api.myai.com/v1',
  },
};

<ModelModal
  customProviders={customProviders}
  // ... 其他属性
/>
```

## 📱 响应式配置

```tsx
<ModelModal
  config={{
    styles: {
      modalWidth: { xs: '95%', sm: 600, md: 800, lg: 1000 },
      sidebarWidth: { xs: 150, sm: 180, md: 200 },
    },
  }}
  // ... 其他属性
/>
```

## 🧪 测试集成

### 单元测试

```tsx
import { render, screen } from '@testing-library/react';
import { ModelModal } from './components/ModelModal/src';

// Mock服务
const mockModelService = {
  createModel: jest.fn(),
  getModelNameList: jest.fn(),
  testModel: jest.fn(),
  updateModel: jest.fn(),
};

test('renders ModelModal', () => {
  render(
    <ModelModal
      open={true}
      type="chat"
      onClose={jest.fn()}
      refresh={jest.fn()}
      modelService={mockModelService}
    />
  );
  
  expect(screen.getByText(/模型供应商/)).toBeInTheDocument();
});
```

### 集成测试

```tsx
test('can submit form', async () => {
  const mockService = {
    createModel: jest.fn().mockResolvedValue({ ModelName: '1' }),
    // ... 其他方法
  };
  
  render(
    <ModelModal
      open={true}
      type="chat"
      onClose={jest.fn()}
      refresh={jest.fn()}
      modelService={mockService}
    />
  );
  
  // 填写表单并提交
  // ... 测试逻辑
});
```

## 🚨 常见问题

### 问题1: 依赖冲突

**症状**: 运行时出现依赖版本冲突错误

**解决方案**: 
- 检查依赖版本兼容性
- 使用 `npm ls` 查看依赖树
- 考虑使用 `resolutions` 或 `overrides` 强制版本

### 问题2: 样式不生效

**症状**: 组件样式显示异常

**解决方案**:
- 确保 Material-UI 主题配置正确
- 检查 CSS 导入顺序
- 验证 Emotion 配置

### 问题3: TypeScript 类型错误

**症状**: 编译时出现类型错误

**解决方案**:
- 确保 TypeScript 版本兼容
- 检查类型定义文件是否正确导入
- 验证接口实现是否完整

### 问题4: 构建失败

**症状**: 构建过程中出现错误

**解决方案**:
- 检查 Node.js 版本兼容性
- 清理 node_modules 重新安装
- 验证构建配置是否正确

## 📚 最佳实践

### 1. 服务层抽象

将模型服务逻辑抽象到独立的服务层：

```tsx
// services/modelService.ts
export class ModelService {
  private baseURL: string;
  private headers: Record<string, string>;
  
  constructor(baseURL: string, headers: Record<string, string> = {}) {
    this.baseURL = baseURL;
    this.headers = headers;
  }
  
  // 实现接口方法...
}

// 在组件中使用
const modelService = new ModelService('/api');
```

### 2. 配置管理

将配置集中管理：

```tsx
// config/modelModal.ts
export const modelModalConfig = {
  theme: {
    primaryColor: process.env.REACT_APP_PRIMARY_COLOR || '#1976d2',
  },
  features: {
    enableModelTesting: process.env.NODE_ENV === 'development',
  },
};

// 在组件中使用
<ModelModal config={modelModalConfig} />
```

### 3. 错误处理

实现统一的错误处理：

```tsx
const handleError = (error: string) => {
  // 记录错误日志
  console.error('ModelModal Error:', error);
  
  // 显示用户友好的错误消息
  showNotification('error', '操作失败，请稍后重试');
  
  // 上报错误
  reportError('ModelModal', error);
};
```

### 4. 性能优化

```tsx
// 使用 useMemo 缓存配置
const memoizedConfig = useMemo(() => ({
  theme: { primaryColor: theme.palette.primary.main },
  features: { enableModelTesting: true },
}), [theme.palette.primary.main]);

// 使用 useCallback 缓存回调函数
const handleClose = useCallback(() => {
  setOpen(false);
}, []);
```

## 🔗 相关资源

- [Material-UI 文档](https://mui.com/)
- [React Hook Form 文档](https://react-hook-form.com/)
- [TypeScript 文档](https://www.typescriptlang.org/)
- [ModelKit 项目](https://github.com/chaitin/ModelKit)

## 🤝 获取帮助

如果遇到问题，可以：

1. 查看本文档的常见问题部分
2. 检查示例代码
3. 提交 Issue 到 ModelKit 项目
4. 联系开发团队 