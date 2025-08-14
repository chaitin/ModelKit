# ModelModal

一个用于管理AI模型配置的React组件，基于Material-UI构建。

## 安装

```bash
npm install @your-org/model-modal
```

## 使用方法

### 基本用法

```tsx
import React from 'react';
import { ModelModal, ConstsModelType } from '@your-org/model-modal';
import { ThemeProvider, createTheme } from '@mui/material/styles';

// 创建主题时需要包含 paper2 背景色
const theme = createTheme({
  palette: {
    background: {
      default: '#fff',
      paper: '#F1F2F8',
      paper2: '#F8F9FA', // 重要：需要定义 paper2 背景色
    },
  },
});

function App() {
  const [open, setOpen] = React.useState(false);
  
  const modelService = {
    createModel: async (data) => {
      // 实现创建模型的逻辑
      return { model: {} };
    },
    listModel: async (data) => {
      // 实现获取模型列表的逻辑
      return { models: [] };
    },
    checkModel: async (data) => {
      // 实现检查模型的逻辑
      return { model: {} };
    },
    updateModel: async (data) => {
      // 实现更新模型的逻辑
      return { model: {} };
    },
  };

  return (
    <ThemeProvider theme={theme}>
      <ModelModal
        open={open}
        onClose={() => setOpen(false)}
        refresh={() => console.log('refresh')}
        data={null}
        type={ConstsModelType.ModelTypeLLM}
        modelService={modelService}
      />
    </ThemeProvider>
  );
}
```

### 主题配置

为了确保组件样式正确显示，你需要在主题中定义 `background.paper2` 属性：

```tsx
import { createTheme } from '@mui/material/styles';
import { mergeThemeWithDefaults } from '@your-org/model-modal';

// 方法1：直接在主题中定义
const theme = createTheme({
  palette: {
    background: {
      default: '#fff',
      paper: '#F1F2F8',
      paper2: '#F8F9FA', // 必需
    },
  },
});

// 方法2：使用提供的合并工具
const baseTheme = createTheme({
  palette: {
    primary: {
      main: '#1976d2',
    },
  },
});

const theme = mergeThemeWithDefaults(baseTheme);
```

### TypeScript 支持

如果你使用 TypeScript，需要确保主题类型扩展被正确导入：

```tsx
// 在你的类型声明文件中（如 vite-env.d.ts 或 global.d.ts）
import '@your-org/model-modal/dist/types/theme';
```

或者在使用组件的文件中导入：

```tsx
import '@your-org/model-modal/dist/types/theme';
import { ModelModal } from '@your-org/model-modal';
```

## 常见问题

### 样式显示不正确

如果组件的样式显示不正确，通常是因为主题中缺少 `background.paper2` 的定义。请确保：

1. 在主题配置中添加 `background.paper2` 属性
2. 使用 `ThemeProvider` 包装你的应用
3. 导入了正确的主题类型声明

### TypeScript 类型错误

如果遇到 TypeScript 类型错误，如 "Property 'paper2' does not exist"，请确保：

1. 导入了主题类型扩展：`import '@your-org/model-modal/dist/types/theme'`
2. 重启 TypeScript 服务

## API

### ModelModalProps

| 属性 | 类型 | 必需 | 描述 |
|------|------|------|------|
| open | boolean | ✓ | 控制模态框的显示/隐藏 |
| onClose | () => void | ✓ | 关闭模态框的回调函数 |
| refresh | () => void | ✓ | 刷新数据的回调函数 |
| data | Model \| null | ✓ | 编辑时的模型数据，新建时为 null |
| type | ConstsModelType | ✓ | 模型类型 |
| modelService | ModelService | ✓ | 模型服务接口 |

### ModelService

```tsx
interface ModelService {
  createModel: (data: CreateModelReq) => Promise<{ model: Model }>;
  listModel: (data: ListModelReq) => Promise<{ models: ModelListItem[] }>;
  checkModel: (data: CheckModelReq) => Promise<{ model: Model }>;
  updateModel: (data: UpdateModelReq) => Promise<{ model: Model }>;
}
```

## 许可证

MIT
