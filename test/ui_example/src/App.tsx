import React, { useState } from 'react';
import {
  Box,
  Button,
  Container,
  Typography,
  Paper,
  Stack,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  SelectChangeEvent,
} from '@mui/material';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { Model} from '@yokowu/modelkit-ui';
import { localModelService } from './localService';
import { ModelModal } from '@yokowu/modelkit-ui';
import toast, { Toaster } from 'react-hot-toast';

// 创建Material-UI主题
const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
  },
});

function App() {
  const [modalOpen, setModalOpen] = useState(false);
  const [selectedType, setSelectedType] = useState<string>('llm');
  const [editingModel, setEditingModel] = useState<Model | null>(null);

  // 使用react-hot-toast的消息组件
  const messageComponent = {
    error: (content: string | Error) => {
      const message = content instanceof Error ? content.message : String(content);
      return toast.error(message);
    },
    success: (content: string) => toast.success(String(content)),
    info: (content: string) => toast(String(content)),
    warning: (content: string) => toast(String(content), { icon: '⚠️' }),
  };

  const handleOpenEditModal = () => {
    // 模拟编辑一个现有模型
    const mockModel: Model = {
      id: 'edit-example',
      model_name: 'gpt-4',
      show_name: 'GPT-4 示例',
      provider: 'OpenAI' as any,
      model_type: 'llm',
      base_url: 'https://api.openai.com/v1',
      api_key: 'sk-example-key',
      api_version: '2023-05-15',
      status: 'active' as any,
      is_active: true,
      created_at: Date.now(),
      updated_at: Date.now(),
    };
    setEditingModel(mockModel);
    setModalOpen(true);
  };

  const handleCloseModal = () => {
    setModalOpen(false);
    setEditingModel(null);
  };

  const handleRefresh = () => {
    console.log('刷新模型列表');
  };

  const handleTypeChange = (event: SelectChangeEvent) => {
    setSelectedType(event.target.value);
  };



  return (
    <>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <Toaster position="top-right" />
        <Container maxWidth="md" sx={{ py: 4 }}>
        <Paper elevation={3} sx={{ p: 4 }}>
          <Typography variant="h4" component="h1" gutterBottom align="center">
            ModelModal 示例程序
          </Typography>
          
          <Typography variant="body1" paragraph align="center" color="text.secondary">
            这是一个用于测试 ModelModal 组件的示例应用程序
          </Typography>

          <Box sx={{ mb: 3 }}>
            <FormControl fullWidth>
              <InputLabel>模型类型</InputLabel>
              <Select
                value={selectedType}
                label="模型类型"
                onChange={handleTypeChange}
              >
                <MenuItem value={'llm'}>对话模型 (LLM)</MenuItem>
                <MenuItem value={'coder'}>代码补全模型</MenuItem>
                <MenuItem value={'embedding'}>向量模型</MenuItem>
                <MenuItem value={'audio'}>音频模型</MenuItem>
                <MenuItem value={'reranker'}>重排序模型</MenuItem>
              </Select>
            </FormControl>
          </Box>

          <Stack spacing={2} direction="row" justifyContent="center">
            <Button 
              variant="contained" 
              size="large"
              onClick={() => {
                setEditingModel(null);
                setModalOpen(true);
              }}
            >
              添加新模型
            </Button>
            
            <Button 
              variant="outlined" 
              size="large"
              onClick={handleOpenEditModal}
            >
              编辑示例模型
            </Button>
            
            <Button 
              variant="text" 
              size="large"
              onClick={() => {
                messageComponent.success('测试成功消息');
                setTimeout(() => messageComponent.error('测试错误消息'), 1000);
                setTimeout(() => messageComponent.info('测试信息消息'), 2000);
                setTimeout(() => messageComponent.warning('测试警告消息'), 3000);
              }}
            >
              测试消息提醒
            </Button>
          </Stack>

          <Box sx={{ mt: 4 }}>
            <Typography variant="h6" gutterBottom>
              功能说明：
            </Typography>
            <Typography variant="body2" component="ul" sx={{ pl: 2 }}>
              <li>支持多种模型类型选择</li>
              <li>支持添加和编辑模型配置</li>
              <li>集成多个AI服务提供商</li>
              <li>实时模型验证和测试</li>
              <li>响应式设计，支持移动端</li>
            </Typography>
          </Box>
        </Paper>

        {/* SimpleModelModal 组件 */}
        <ModelModal
          open={modalOpen}
          onClose={handleCloseModal}
          refresh={handleRefresh}
          data={editingModel}
          model_type={selectedType}
          modelService={localModelService}
          language="zh-CN"
          messageComponent={messageComponent}
          is_close_model_remark={true}
        />
        

        </Container>
      </ThemeProvider>
    </>
  );
}

export default App;