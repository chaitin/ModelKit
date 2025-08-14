# ModelModal

一个用于管理AI模型配置的React组件，基于Material-UI构建。

## 安装

```bash
npm install @yokowu/modelkit-ui
# 或
yarn add @yokowu/modelkit-ui
# 或
pnpm add @yokowu/modelkit-ui
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

## 组件

### ModelAdd

模型添加/编辑弹窗组件。

#### Props

- `open: boolean` - 是否显示弹窗
- `data: ModelListItem | null` - 编辑时的模型数据
- `type: 'chat' | 'embedding' | 'rerank'` - 模型类型
- `onClose: () => void` - 关闭弹窗回调
- `refresh: () => void` - 刷新数据回调

## 依赖要求

确保你的项目已安装以下依赖：

```json
{
  "peerDependencies": {
    "react": ">=18.0.0",
    "react-dom": ">=18.0.0",
    "@mui/material": ">=5.0.0",
    "@mui/icons-material": ">=5.0.0",
    "@emotion/react": ">=11.0.0",
    "@emotion/styled": ">=11.0.0",
    "react-hook-form": ">=7.0.0"
  }
}
```

## 开发

本项目使用 pnpm 作为包管理器：

```bash
# 安装依赖
pnpm install

# 开发模式
pnpm dev

# 构建
pnpm build
```

## 许可证

MIT
