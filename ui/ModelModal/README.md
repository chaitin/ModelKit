# ModelKit UI

一个基于 React 和 Material-UI 的模型管理组件库。

## 安装

### 从 GitHub Packages 安装

由于此包发布在 GitHub Packages，需要先配置 npm 注册表。你可以通过以下几种方式进行配置：

#### 方式一：全局配置（推荐）

```bash
npm config set @yokowu:registry https://npm.pkg.github.com
```

#### 方式二：使用 .npmrc 文件

在项目根目录创建 `.npmrc` 文件：

### 登录 GitHub Packages, 账号密码找于康

```bash
npm login --registry=https://npm.pkg.github.com --scope=@yokowu
```

然后安装包：

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
    "react": ">=19.0.0",
    "react-dom": ">=19.0.0",
    "@mui/material": ">=6.0.0",
    "@mui/icons-material": ">=6.0.0",
    "@emotion/react": ">=11.0.0",
    "@emotion/styled": ">=11.0.0"
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
