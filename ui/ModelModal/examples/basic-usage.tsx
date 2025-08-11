import React, { useState } from 'react';
import { Button, Box, Typography, Container } from '@mui/material';
import { ModelModal, createDefaultModelService } from '../src';

// 模拟的模型数据
const mockModelData = {
  id: '1',
  model: 'gpt-4',
  type: 'chat' as const,
  provider: 'OpenAI',
  base_url: 'https://api.openai.com/v1',
  api_key: 'sk-xxx',
  api_version: '',
  api_header: '',
};

export const BasicUsageExample: React.FC = () => {
  const [open, setOpen] = useState(false);
  const [modelData, setModelData] = useState<typeof mockModelData | null>(null);
  const [isEditMode, setIsEditMode] = useState(false);

  // 创建默认服务实例
  const modelService = createDefaultModelService('https://api.example.com');

  const handleOpenAdd = () => {
    setModelData(null);
    setIsEditMode(false);
    setOpen(true);
  };

  const handleOpenEdit = () => {
    setModelData(mockModelData);
    setIsEditMode(true);
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setModelData(null);
    setIsEditMode(false);
  };

  const handleRefresh = () => {
    console.log('刷新模型列表');
    // 这里可以调用API刷新数据
  };

  const handleBeforeSubmit = async (data: any) => {
    console.log('提交前验证:', data);
    // 可以在这里添加自定义验证逻辑
    if (data.api_key.length < 10) {
      alert('API密钥长度不足');
      return false;
    }
    return true;
  };

  const handleAfterSubmit = (data: any, result: any) => {
    console.log('提交成功:', { data, result });
    alert(isEditMode ? '模型更新成功！' : '模型添加成功！');
  };

  const handleError = (error: string) => {
    console.error('操作失败:', error);
    alert(`操作失败: ${error}`);
  };

  return (
    <Container maxWidth="md">
      <Box sx={{ py: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          ModelModal 基本使用示例
        </Typography>
        
        <Typography variant="body1" paragraph>
          这是一个展示如何使用 ModelModal 组件的示例。组件支持添加和编辑两种模式。
        </Typography>

        <Box sx={{ display: 'flex', gap: 2, mb: 4 }}>
          <Button variant="contained" onClick={handleOpenAdd}>
            添加新模型
          </Button>
          <Button variant="outlined" onClick={handleOpenEdit}>
            编辑现有模型
          </Button>
        </Box>

        <Typography variant="h6" gutterBottom>
          功能说明：
        </Typography>
        <ul>
          <li>支持添加和编辑两种模式</li>
          <li>内置多种AI模型提供商</li>
          <li>支持模型配置验证</li>
          <li>可自定义主题和样式</li>
          <li>支持国际化</li>
        </ul>

        {/* ModelModal 组件 */}
        <ModelModal
          open={open}
          data={modelData}
          type="chat"
          onClose={handleClose}
          refresh={handleRefresh}
          modelService={modelService}
          onBeforeSubmit={handleBeforeSubmit}
          onAfterSubmit={handleAfterSubmit}
          onError={handleError}
          config={{
            theme: {
              primaryColor: '#1976d2',
              borderRadius: '12px',
            },
            locale: {
              language: 'zh-CN',
            },
            features: {
              enableModelTesting: true,
              enableHeaderConfig: true,
              enableApiVersion: true,
              enableProviderSelection: true,
            },
            styles: {
              modalWidth: 900,
              sidebarWidth: 220,
            },
          }}
        />
      </Box>
    </Container>
  );
};

export default BasicUsageExample; 