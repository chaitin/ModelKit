# ModelKit UI

一个基于 React 和 Material-UI 的模型管理组件库。

## 安装

```bash
npm install @yokowu/modelkit-ui
# 或
yarn add @yokowu/modelkit-ui
# 或
pnpm add @yokowu/modelkit-ui
```

## 使用方法

### 基本使用

```tsx
import React, { useState } from 'react';
import { ModelAdd, ModelProvider } from '@yokowu/modelkit-ui';

function App() {
  const [modalOpen, setModalOpen] = useState(false);

  const handleModalClose = () => {
    setModalOpen(false);
  };

  const handleRefresh = () => {
    // 刷新逻辑
  };

  return (
    <div>
      <ModelAdd
        open={modalOpen}
        data={null}
        type="chat"
        onClose={handleModalClose}
        refresh={handleRefresh}
      />
    </div>
  );
}

export default App;
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
