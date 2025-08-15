import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Box,
  Typography,
  Alert,
  CircularProgress,
  Grid,
  Divider,
} from '@mui/material';
import {
  Model,
  ModelService,
  ConstsModelType,
  ConstsModelProvider,
} from '../../ui/ModelModal/src/types/types';

interface SimpleModelModalProps {
  open: boolean;
  onClose: () => void;
  refresh: () => void;
  data: Model | null;
  type: ConstsModelType;
  modelService: ModelService;
  language?: 'zh-CN' | 'en-US';
}

export const SimpleModelModal: React.FC<SimpleModelModalProps> = ({
  open,
  onClose,
  refresh,
  data,
  type,
  modelService,
  language = 'zh-CN',
}) => {
  const [formData, setFormData] = useState({
    model_name: '',
    show_name: '',
    provider: ConstsModelProvider.ModelProviderOpenAI,
    api_base: '',
    api_key: '',
    api_version: '',
    api_header: '',
  });

  const [loading, setLoading] = useState(false);
  const [testing, setTesting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  // 当模态框打开或数据变化时，初始化表单
  useEffect(() => {
    if (open) {
      if (data) {
        // 编辑模式
        setFormData({
          model_name: data.model_name || '',
          show_name: data.show_name || '',
          provider: data.provider || ConstsModelProvider.ModelProviderOpenAI,
          api_base: data.api_base || '',
          api_key: data.api_key || '',
          api_version: data.api_version || '',
          api_header: data.api_header || '',
        });
      } else {
        // 新增模式
        setFormData({
          model_name: '',
          show_name: '',
          provider: ConstsModelProvider.ModelProviderOpenAI,
          api_base: '',
          api_key: '',
          api_version: '',
          api_header: '',
        });
      }
      setError(null);
      setSuccess(null);
    }
  }, [open, data]);

  const handleInputChange = (field: string, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    setError(null);
    setSuccess(null);
  };

  const handleSubmit = async () => {
    if (!formData.model_name || !formData.api_key || !formData.api_base) {
      setError('请填写必填字段：模型名称、API密钥和API地址');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      if (data) {
        // 更新模型
        await modelService.updateModel({
          id: data.id,
          model_name: formData.model_name,
          show_name: formData.show_name,
          provider: formData.provider as any,
          api_base: formData.api_base,
          api_key: formData.api_key,
          api_version: formData.api_version,
          api_header: formData.api_header,
        });
        setSuccess('模型更新成功！');
      } else {
        // 创建新模型
        await modelService.createModel({
          model_name: formData.model_name,
          show_name: formData.show_name,
          provider: formData.provider as any,
          model_type: type,
          api_base: formData.api_base,
          api_key: formData.api_key,
          api_version: formData.api_version,
          api_header: formData.api_header,
        });
        setSuccess('模型创建成功！');
      }
      
      refresh();
      setTimeout(() => {
        onClose();
      }, 1500);
    } catch (err: any) {
      setError(err.message || '操作失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const handleTestConnection = async () => {
    if (!formData.model_name || !formData.api_key || !formData.api_base) {
      setError('请先填写模型名称、API密钥和API地址');
      return;
    }

    setTesting(true);
    setError(null);
    setSuccess(null);

    try {
      await modelService.checkModel({
        model_name: formData.model_name,
        provider: formData.provider,
        type: type as any,
        api_base: formData.api_base,
        api_key: formData.api_key,
        api_version: formData.api_version,
        api_header: formData.api_header,
      });
      setSuccess('连接测试成功！');
    } catch (err: any) {
      setError(`连接测试失败: ${err.message}`);
    } finally {
      setTesting(false);
    }
  };

  const providerOptions = [
    { value: ConstsModelProvider.ModelProviderOpenAI, label: 'OpenAI' },
    { value: ConstsModelProvider.ModelProviderOllama, label: 'Ollama' },
    { value: ConstsModelProvider.ModelProviderDeepSeek, label: 'DeepSeek' },
    { value: ConstsModelProvider.ModelProviderMoonshot, label: 'Moonshot' },
    { value: ConstsModelProvider.ModelProviderAzureOpenAI, label: 'Azure OpenAI' },
    { value: ConstsModelProvider.ModelProviderBaiZhiCloud, label: '百智云' },
    { value: ConstsModelProvider.ModelProviderHunyuan, label: '腾讯混元' },
    { value: ConstsModelProvider.ModelProviderBaiLian, label: '阿里百炼' },
    { value: ConstsModelProvider.ModelProviderVolcengine, label: '火山引擎' },
    { value: ConstsModelProvider.ModelProviderSiliconFlow, label: 'SiliconFlow' },
  ];

  const getTypeLabel = (type: ConstsModelType) => {
    const typeMap = {
      [ConstsModelType.ModelTypeLLM]: '对话模型',
      [ConstsModelType.ModelTypeCoder]: '代码补全模型',
      [ConstsModelType.ModelTypeEmbedding]: '向量模型',
      [ConstsModelType.ModelTypeAudio]: '音频模型',
      [ConstsModelType.ModelTypeReranker]: '重排序模型',
    };
    return typeMap[type] || type;
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        {data ? '编辑模型' : '添加新模型'} - {getTypeLabel(type)}
      </DialogTitle>
      
      <DialogContent>
        <Box sx={{ pt: 2 }}>
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}
          
          {success && (
            <Alert severity="success" sx={{ mb: 2 }}>
              {success}
            </Alert>
          )}

          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>服务提供商</InputLabel>
                <Select
                  value={formData.provider}
                  label="服务提供商"
                  onChange={(e) => handleInputChange('provider', e.target.value)}
                >
                  {providerOptions.map((option) => (
                    <MenuItem key={option.value} value={option.value}>
                      {option.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="模型名称"
                value={formData.model_name}
                onChange={(e) => handleInputChange('model_name', e.target.value)}
                required
                placeholder="例如: gpt-4, claude-3-sonnet"
              />
            </Grid>

            <Grid item xs={12}>
              <TextField
                fullWidth
                label="显示名称"
                value={formData.show_name}
                onChange={(e) => handleInputChange('show_name', e.target.value)}
                placeholder="模型的友好显示名称"
              />
            </Grid>

            <Grid item xs={12}>
              <TextField
                fullWidth
                label="API地址"
                value={formData.api_base}
                onChange={(e) => handleInputChange('api_base', e.target.value)}
                required
                placeholder="例如: https://api.openai.com/v1"
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="API密钥"
                type="password"
                value={formData.api_key}
                onChange={(e) => handleInputChange('api_key', e.target.value)}
                required
                placeholder="sk-..."
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="API版本"
                value={formData.api_version}
                onChange={(e) => handleInputChange('api_version', e.target.value)}
                placeholder="例如: 2023-05-15"
              />
            </Grid>

            <Grid item xs={12}>
              <TextField
                fullWidth
                label="自定义请求头"
                value={formData.api_header}
                onChange={(e) => handleInputChange('api_header', e.target.value)}
                placeholder="例如: Authorization: Bearer token"
                helperText="可选：自定义HTTP请求头"
              />
            </Grid>
          </Grid>

          <Divider sx={{ my: 3 }} />

          <Box sx={{ display: 'flex', gap: 2, justifyContent: 'center' }}>
            <Button
              variant="outlined"
              onClick={handleTestConnection}
              disabled={testing || loading}
              startIcon={testing ? <CircularProgress size={16} /> : null}
            >
              {testing ? '测试中...' : '测试连接'}
            </Button>
          </Box>
        </Box>
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose} disabled={loading}>
          取消
        </Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          disabled={loading || testing}
          startIcon={loading ? <CircularProgress size={16} /> : null}
        >
          {loading ? '保存中...' : (data ? '更新' : '创建')}
        </Button>
      </DialogActions>
    </Dialog>
  );
};