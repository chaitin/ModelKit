# ModelModal 快速开始

5分钟内快速集成 ModelModal 组件到你的项目中！

## ⚡ 超快速集成

### 1. 复制组件（30秒）

```bash
# 复制整个组件目录
cp -r ui/src/components/ModelModal your-project/src/components/
```

### 2. 安装依赖（1分钟）

```bash
cd your-project
npm install @mui/material @emotion/react @emotion/styled react-hook-form
```

### 3. 使用组件（2分钟）

```tsx
import React, { useState } from 'react';
import { ModelModal, createDefaultModelService } from './components/ModelModal/src';

function App() {
  const [open, setOpen] = useState(false);
  
  // 创建服务实例
  const modelService = createDefaultModelService('https://api.example.com');
  
  return (
    <div>
      <button onClick={() => setOpen(true)}>添加AI模型</button>
      
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

### 4. 运行项目（1分钟）

```bash
npm start
```

🎉 **完成！** 你现在有了一个功能完整的AI模型配置界面！

## 🔧 自定义配置（可选）

### 主题定制

```tsx
<ModelModal
  config={{
    theme: {
      primaryColor: '#ff6b6b',
      borderRadius: '16px',
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
      enableModelTesting: false,    // 禁用模型测试
      enableHeaderConfig: true,     // 启用Header配置
    },
  }}
  // ... 其他属性
/>
```

### 本地化

```tsx
<ModelModal
  config={{
    locale: {
      language: 'en-US',            // 英文界面
    },
  }}
  // ... 其他属性
/>
```

## 📚 下一步

- 📖 查看 [完整文档](README.md)
- 🔗 阅读 [集成指南](INTEGRATION_GUIDE.md)
- 💡 探索 [使用示例](examples/)
- 🐛 遇到问题？查看 [常见问题](INTEGRATION_GUIDE.md#常见问题)

## 🚀 高级功能

- 自定义模型提供商
- 自定义验证规则
- 响应式设计
- 事件回调
- 异步验证

---

**需要帮助？** 查看完整文档或联系开发团队！ 