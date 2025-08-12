import React, { useState } from 'react';
import { Button, Box, Typography, Container, Switch, FormControlLabel } from '@mui/material';
import { ModelModal, createDefaultModelService, createModelModal } from '../src';

// 自定义模型提供商
const customProviders = {
  CustomAI: {
    label: 'CustomAI',
    cn: '自定义AI',
    icon: 'icon-custom',
    urlWrite: true,
    secretRequired: true,
    customHeader: true,
    modelDocumentUrl: 'https://docs.customai.com',
    defaultBaseUrl: 'https://api.customai.com/v1',
  },
  EnterpriseAI: {
    label: 'EnterpriseAI',
    cn: '企业级AI',
    icon: 'icon-enterprise',
    urlWrite: false,
    secretRequired: true,
    customHeader: false,
    modelDocumentUrl: 'https://docs.enterpriseai.com',
    defaultBaseUrl: 'https://api.enterpriseai.com/v1',
  },
};

// 自定义验证器
const customValidators = {
  'api_key': (value: string) => {
    if (value.length < 20) return 'API密钥长度必须至少20个字符';
    if (!/^sk-[a-zA-Z0-9]+$/.test(value)) return 'API密钥格式不正确，应以sk-开头';
    return undefined;
  },
  'base_url': (value: string) => {
    if (!value.startsWith('https://')) return 'API地址必须使用HTTPS协议';
    return undefined;
  },
};

export const AdvancedConfigExample: React.FC = () => {
  const [open, setOpen] = useState(false);
  const [config, setConfig] = useState({
    enableModelTesting: true,
    enableHeaderConfig: true,
    enableApiVersion: true,
    enableProviderSelection: true,
    language: 'zh-CN' as const,
    theme: 'default' as const,
  });

  // 创建默认服务实例
  const modelService = createDefaultModelService('https://api.example.com');

  // 创建预配置的组件
  const CustomModelModal = createModelModal({
    theme: {
      primaryColor: config.theme === 'default' ? '#1976d2' : '#ff6b6b',
      secondaryColor: config.theme === 'default' ? '#dc004e' : '#4ecdc4',
      borderRadius: '16px',
      spacing: 3,
    },
    locale: {
      language: config.language,
      messages: {
        'custom.welcome': '欢迎使用自定义AI模型配置',
        'custom.success': '配置成功！',
      },
    },
    validation: {
      requiredFields: ['provider', 'model', 'base_url', 'api_key'],
      customValidators,
    },
    features: {
      enableModelTesting: config.enableModelTesting,
      enableHeaderConfig: config.enableHeaderConfig,
      enableApiVersion: config.enableApiVersion,
      enableProviderSelection: config.enableProviderSelection,
    },
    styles: {
      modalWidth: 1000,
      sidebarWidth: 250,
      customCSS: `
        .custom-modal .MuiBox-root {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }
        .custom-modal .MuiTypography-root {
          color: white;
        }
      `,
    },
  });

  const handleOpen = () => setOpen(true);
  const handleClose = () => setOpen(false);

  const handleRefresh = () => {
    console.log('刷新模型列表');
  };

  const handleBeforeSubmit = async (data: any) => {
    console.log('高级验证:', data);
    
    // 模拟异步验证
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    // 自定义业务逻辑验证
    if (data.provider === 'CustomAI' && data.api_key.length < 25) {
      alert('CustomAI提供商要求API密钥长度至少25个字符');
      return false;
    }
    
    return true;
  };

  const handleAfterSubmit = (data: any, result: any) => {
    console.log('高级提交成功:', { data, result });
    alert('高级配置提交成功！');
  };

  const handleError = (error: string) => {
    console.error('高级操作失败:', error);
    alert(`高级操作失败: ${error}`);
  };

  return (
    <Container maxWidth="lg">
      <Box sx={{ py: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          ModelModal 高级配置示例
        </Typography>
        
        <Typography variant="body1" paragraph>
          这个示例展示了如何使用高级配置功能，包括自定义提供商、验证器、主题等。
        </Typography>

        {/* 配置控制面板 */}
        <Box sx={{ mb: 4, p: 3, border: '1px solid', borderColor: 'divider', borderRadius: 2 }}>
          <Typography variant="h6" gutterBottom>
            功能配置
          </Typography>
          
          <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 2 }}>
            <FormControlLabel
              control={
                <Switch
                  checked={config.enableModelTesting}
                  onChange={(e) => setConfig(prev => ({ ...prev, enableModelTesting: e.target.checked }))}
                />
              }
              label="启用模型测试"
            />
            <FormControlLabel
              control={
                <Switch
                  checked={config.enableHeaderConfig}
                  onChange={(e) => setConfig(prev => ({ ...prev, enableHeaderConfig: e.target.checked }))}
                />
              }
              label="启用Header配置"
            />
            <FormControlLabel
              control={
                <Switch
                  checked={config.enableApiVersion}
                  onChange={(e) => setConfig(prev => ({ ...prev, enableApiVersion: e.target.checked }))}
                />
              }
              label="启用API版本"
            />
            <FormControlLabel
              control={
                <Switch
                  checked={config.enableProviderSelection}
                  onChange={(e) => setConfig(prev => ({ ...prev, enableProviderSelection: e.target.checked }))}
                />
              }
              label="启用提供商选择"
            />
          </Box>

          <Box sx={{ mt: 2 }}>
            <FormControlLabel
              control={
                <Switch
                  checked={config.language === 'en-US'}
                  onChange={(e) => setConfig(prev => ({ ...prev, language: e.target.checked ? 'en-US' : 'zh-CN' }))}
                />
              }
              label="英文界面"
            />
            <FormControlLabel
              control={
                <Switch
                  checked={config.theme === 'custom'}
                  onChange={(e) => setConfig(prev => ({ ...prev, theme: e.target.checked ? 'custom' : 'default' }))}
                />
              }
              label="自定义主题"
            />
          </Box>
        </Box>

        <Box sx={{ display: 'flex', gap: 2, mb: 4 }}>
          <Button variant="contained" onClick={handleOpen} color="primary">
            打开高级配置模态框
          </Button>
        </Box>

        <Typography variant="h6" gutterBottom>
          高级功能说明：
        </Typography>
        <ul>
          <li>自定义模型提供商配置</li>
          <li>自定义验证规则</li>
          <li>动态主题切换</li>
          <li>功能开关控制</li>
          <li>预配置组件创建</li>
          <li>自定义CSS样式</li>
          <li>异步验证支持</li>
        </ul>

        {/* 使用预配置的组件 */}
        <CustomModelModal
          open={open}
          data={null}
          type="chat"
          onClose={handleClose}
          refresh={handleRefresh}
          modelService={modelService}
          customProviders={customProviders}
          onBeforeSubmit={handleBeforeSubmit}
          onAfterSubmit={handleAfterSubmit}
          onError={handleError}
        />
      </Box>
    </Container>
  );
};

export default AdvancedConfigExample; 