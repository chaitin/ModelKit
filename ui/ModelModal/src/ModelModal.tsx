import React, { useEffect, useState } from 'react';
import { useForm, Controller } from 'react-hook-form';
import {
  Box,
  Button,
  MenuItem,
  Stack,
  TextField,
  useTheme,
  Paper,
} from '@mui/material';
import { Info as InfoIcon } from '@mui/icons-material';

import {
  ModelModalProps,
  AddModelForm,
  ModelListItem,
  ModelService,
  ModelModalConfig,
} from './types';
import { DEFAULT_MODEL_PROVIDERS, getProvidersByType } from './constants/providers';
import { getLocaleMessage, getTitleMap } from './constants/locale';

// 简单的Card组件，避免外部依赖
const SimpleCard: React.FC<{ children: React.ReactNode; sx?: any }> = ({ children, sx }) => (
  <Paper
    sx={{
      borderRadius: '10px',
      boxShadow: 'none',
      border: 'none',
      overflow: 'hidden',
      ...sx,
    }}
  >
    {children}
  </Paper>
);

// 添加透明度到颜色的工具函数
const addOpacityToColor = (color: string, opacity: number): string => {
  if (color.startsWith('#')) {
    const hex = color.slice(1);
    const r = parseInt(hex.slice(0, 2), 16);
    const g = parseInt(hex.slice(2, 4), 16);
    const b = parseInt(hex.slice(4, 6), 16);
    return `rgba(${r}, ${g}, ${b}, ${opacity})`;
  }
  return color;
};

export const ModelModal: React.FC<ModelModalProps> = ({
  open,
  data,
  type,
  onClose,
  refresh,
  modelService,
  config = {},
  onBeforeSubmit,
  onAfterSubmit,
  onError,
  customProviders,
  customValidation,
  customStyles,
}) => {
  const theme = useTheme();
  
  // 合并配置
  const mergedConfig: Required<ModelModalConfig> = {
    theme: {
      primaryColor: config.theme?.primaryColor || theme.palette.primary.main,
      secondaryColor: config.theme?.secondaryColor || theme.palette.secondary.main,
      borderRadius: config.theme?.borderRadius || '10px',
      spacing: config.theme?.spacing || 2,
    },
    locale: {
      language: config.locale?.language || 'zh-CN',
      messages: config.locale?.messages || {},
    },
    validation: {
      requiredFields: config.validation?.requiredFields || ['provider', 'model', 'base_url', 'api_key'],
      customValidators: config.validation?.customValidators || {},
    },
    features: {
      enableModelTesting: config.features?.enableModelTesting ?? true,
      enableHeaderConfig: config.features?.enableHeaderConfig ?? true,
      enableApiVersion: config.features?.enableApiVersion ?? true,
      enableProviderSelection: config.features?.enableProviderSelection ?? true,
    },
    styles: {
      modalWidth: config.styles?.modalWidth || 800,
      sidebarWidth: config.styles?.sidebarWidth || 200,
      customCSS: config.styles?.customCSS || '',
    },
  };

  // 使用自定义提供商或默认提供商
  const providers = customProviders || DEFAULT_MODEL_PROVIDERS;
  const filteredProviders = getProvidersByType(type, providers);

  const {
    formState: { errors },
    handleSubmit,
    control,
    reset,
    setValue,
    watch,
  } = useForm<AddModelForm>({
    defaultValues: {
      type,
      provider: 'BaiZhiCloud',
      base_url: filteredProviders.BaiZhiCloud?.defaultBaseUrl || '',
      model: '',
      api_version: '',
      api_key: '',
      api_header_key: '',
      api_header_value: '',
    },
  });

  const providerBrand = watch('provider');
  const [modelUserList, setModelUserList] = useState<{ model: string }[]>([]);
  const [loading, setLoading] = useState(false);
  const [modelLoading, setModelLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);

  const handleReset = () => {
    onClose();
    reset({
      type,
      provider: 'BaiZhiCloud',
      model: '',
      base_url: '',
      api_key: '',
      api_version: '',
      api_header_key: '',
      api_header_value: '',
    });
    setModelUserList([]);
    setSuccess(false);
    setLoading(false);
    setModelLoading(false);
    setError('');
    refresh();
  };

  const getModel = (value: AddModelForm) => {
    let header = '';
    if (value.api_header_key && value.api_header_value) {
      header = value.api_header_key + '=' + value.api_header_value;
    }
    
    setModelLoading(true);
    modelService
      .getModelNameList({
        type,
        api_key: value.api_key,
        base_url: value.base_url,
        provider: value.provider,
        api_header: header,
      })
      .then((res: { models: { model: string }[] }) => {
        setModelUserList(
          (res.models || []).sort((a, b) => a.model.localeCompare(b.model))
        );
        if (data && (res.models || []).find((it) => it.model === data.model)) {
          setValue('model', data.model);
        } else {
          setValue('model', res.models?.[0]?.model || '');
        }
        setSuccess(true);
      })
      .finally(() => {
        setModelLoading(false);
      });
  };

  const onSubmit = async (value: AddModelForm) => {
    // 执行前置验证
    if (onBeforeSubmit) {
      const shouldContinue = await onBeforeSubmit(value);
      if (!shouldContinue) return;
    }

    let header = '';
    if (value.api_header_key && value.api_header_value) {
      header = value.api_header_key + '=' + value.api_header_value;
    }
    
    setError('');
    setLoading(true);
    
    try {
      if (mergedConfig.features.enableModelTesting) {
        const testResult = await modelService.testModel({
          type,
          api_key: value.api_key,
          base_url: value.base_url,
          api_version: value.api_version,
          provider: value.provider,
          model: value.model,
          api_header: header,
        });
        
        if (testResult.error) {
          setError(testResult.error);
          setLoading(false);
          onError?.(testResult.error);
          return;
        }
      }

      if (data) {
        await modelService.updateModel({
          type,
          api_key: value.api_key,
          base_url: value.base_url,
          model: value.model,
          api_header: header,
          api_version: value.api_version,
          ModelName: data.id,
          provider: value.provider,
        });
        onAfterSubmit?.(value, { ModelName: data.id });
      } else {
        const result = await modelService.createModel({
          type,
          api_key: value.api_key,
          base_url: value.base_url,
          model: value.model,
          api_header: header,
          provider: value.provider,
        });
        onAfterSubmit?.(value, result);
      }
      
      handleReset();
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '操作失败';
      setError(errorMessage);
      onError?.(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const resetCurData = (value: ModelListItem) => {
    if (value.provider && value.provider !== 'Other') {
      getModel({
        api_key: value.api_key || '',
        base_url: value.base_url || '',
        model: value.model || '',
        provider: value.provider || '',
        api_version: value.api_version || '',
        api_header_key: value.api_header?.split('=')[0] || '',
        api_header_value: value.api_header?.split('=')[1] || '',
        type,
      });
    }
    reset({
      type,
      provider: value.provider || 'Other',
      model: value.model || '',
      base_url: value.base_url || '',
      api_key: value.api_key || '',
      api_version: value.api_version || '',
      api_header_key: value.api_header?.split('=')[0] || '',
      api_header_value: value.api_header?.split('=')[1] || '',
    });
  };

  useEffect(() => {
    if (open) {
      if (data) {
        resetCurData(data);
      } else {
        reset({
          type,
          provider: 'BaiZhiCloud',
          model: '',
          base_url: filteredProviders.BaiZhiCloud?.defaultBaseUrl || '',
          api_key: '',
          api_version: '',
          api_header_key: '',
          api_header_value: '',
        });
      }
    }
  }, [data, open, type, filteredProviders.BaiZhiCloud?.defaultBaseUrl]);

  if (!open) return null;

  const titleMap = getTitleMap(mergedConfig.locale.language);
  const isEdit = !!data;
  const title = isEdit
    ? `${getLocaleMessage('edit', mergedConfig.locale.language)} ${titleMap[type]}`
    : `${getLocaleMessage('add', mergedConfig.locale.language)} ${titleMap[type]}`;

  return (
    <Box
      sx={{
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: 'rgba(0, 0, 0, 0.5)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        zIndex: 1300,
      }}
      onClick={handleReset}
    >
      <Box
        sx={{
          width: mergedConfig.styles.modalWidth,
          maxHeight: '90vh',
          backgroundColor: 'background.paper',
          borderRadius: mergedConfig.styles.borderRadius,
          boxShadow: 24,
          overflow: 'hidden',
        }}
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <Box
          sx={{
            p: 3,
            borderBottom: '1px solid',
            borderColor: 'divider',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          <Box sx={{ fontSize: 18, fontWeight: 'bold' }}>{title}</Box>
          <Button onClick={handleReset} variant="text">
            {getLocaleMessage('cancel', mergedConfig.locale.language)}
          </Button>
        </Box>

        {/* Content */}
        <Box sx={{ p: 3, maxHeight: 'calc(90vh - 120px)', overflow: 'auto' }}>
          <Stack direction="row" alignItems="stretch" gap={3}>
            {/* Sidebar */}
            {mergedConfig.features.enableProviderSelection && (
              <Stack
                gap={1}
                sx={{
                  width: mergedConfig.styles.sidebarWidth,
                  flexShrink: 0,
                  bgcolor: 'background.paper2',
                  borderRadius: mergedConfig.styles.borderRadius,
                  p: 1,
                }}
              >
                <Box sx={{ fontSize: 14, lineHeight: '24px', fontWeight: 'bold', p: 1 }}>
                  {getLocaleMessage('modelProvider', mergedConfig.locale.language)}
                </Box>
                {Object.values(filteredProviders).map((it) => (
                  <Stack
                    direction="row"
                    alignItems="center"
                    gap={1.5}
                    key={it.label}
                    sx={{
                      cursor: 'pointer',
                      fontSize: 14,
                      lineHeight: '24px',
                      p: 1,
                      borderRadius: mergedConfig.styles.borderRadius,
                      fontWeight: 'bold',
                      ...(providerBrand === it.label && {
                        bgcolor: addOpacityToColor(mergedConfig.theme.primaryColor || '#1976d2', 0.1),
                        color: 'primary.main',
                      }),
                      '&:hover': {
                        color: 'primary.main',
                      },
                    }}
                    onClick={() => {
                      if (data && data.provider === it.label) {
                        resetCurData(data);
                      } else {
                        setModelUserList([]);
                        setError('');
                        setModelLoading(false);
                        setSuccess(false);
                        reset({
                          provider: it.label,
                          base_url: it.label === 'AzureOpenAI' ? '' : it.defaultBaseUrl,
                          model: '',
                          api_version: '',
                          api_key: '',
                          api_header_key: '',
                          api_header_value: '',
                        });
                      }
                    }}
                  >
                    <InfoIcon sx={{ fontSize: 18 }} />
                    {it.cn || it.label || getLocaleMessage('other', mergedConfig.locale.language)}
                  </Stack>
                ))}
              </Stack>
            )}

            {/* Main Form */}
            <Box sx={{ flex: 1 }}>
              {/* API Address */}
              <Box sx={{ fontSize: 14, lineHeight: '32px' }}>
                {getLocaleMessage('apiAddress', mergedConfig.locale.language)}{' '}
                <Box component="span" sx={{ color: 'red' }}>
                  *
                </Box>
              </Box>
              <Controller
                control={control}
                name="base_url"
                rules={{
                  required: {
                    value: true,
                    message: getLocaleMessage('urlRequired', mergedConfig.locale.language),
                  },
                }}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    disabled={!filteredProviders[providerBrand]?.urlWrite}
                    size="small"
                    placeholder={filteredProviders[providerBrand]?.defaultBaseUrl}
                    error={!!errors.base_url}
                    helperText={errors.base_url?.message}
                    onChange={(e) => {
                      field.onChange(e.target.value);
                      setModelUserList([]);
                      setValue('model', '');
                      setSuccess(false);
                    }}
                  />
                )}
              />

              {/* API Secret */}
              <Stack
                direction="row"
                alignItems="center"
                justifyContent="space-between"
                sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}
              >
                <Box>
                  {getLocaleMessage('apiSecret', mergedConfig.locale.language)}
                  {filteredProviders[providerBrand]?.secretRequired && (
                    <Box component="span" sx={{ color: 'red' }}>
                      *
                    </Box>
                  )}
                </Box>
                {filteredProviders[providerBrand]?.modelDocumentUrl && (
                  <Box
                    component="span"
                    sx={{ color: 'primary.main', cursor: 'pointer', ml: 1, textAlign: 'right' }}
                    onClick={() =>
                      window.open(filteredProviders[providerBrand]?.modelDocumentUrl, '_blank')
                    }
                  >
                    {getLocaleMessage('viewDocumentation', mergedConfig.locale.language)}
                  </Box>
                )}
              </Stack>
              <Controller
                control={control}
                name="api_key"
                rules={{
                  required: {
                    value: filteredProviders[providerBrand]?.secretRequired || false,
                    message: getLocaleMessage('secretRequired', mergedConfig.locale.language),
                  },
                }}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    size="small"
                    placeholder=""
                    error={!!errors.api_key}
                    helperText={errors.api_key?.message}
                    onChange={(e) => {
                      field.onChange(e.target.value);
                      setModelUserList([]);
                      setValue('model', '');
                      setSuccess(false);
                    }}
                  />
                )}
              />

              {/* API Version */}
              {providerBrand === 'AzureOpenAI' && mergedConfig.features.enableApiVersion && (
                <>
                  <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                    {getLocaleMessage('apiVersion', mergedConfig.locale.language)}
                  </Box>
                  <Controller
                    control={control}
                    name="api_version"
                    render={({ field }) => (
                      <TextField
                        {...field}
                        fullWidth
                        size="small"
                        placeholder="2024-10-21"
                        error={!!errors.api_version}
                        helperText={errors.api_version?.message}
                        onChange={(e) => {
                          field.onChange(e.target.value);
                          setModelUserList([]);
                          setValue('model', '');
                          setSuccess(false);
                        }}
                      />
                    )}
                  />
                </>
              )}

              {/* Model Selection */}
              {providerBrand === 'Other' ? (
                <>
                  <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                    {getLocaleMessage('modelName', mergedConfig.locale.language)}{' '}
                    <Box component="span" sx={{ color: 'red' }}>
                      *
                    </Box>
                  </Box>
                  <Controller
                    control={control}
                    name="model"
                    rules={{
                      required: getLocaleMessage('modelNameRequired', mergedConfig.locale.language),
                    }}
                    render={({ field }) => (
                      <TextField
                        {...field}
                        fullWidth
                        size="small"
                        placeholder=""
                        error={!!errors.model}
                        helperText={errors.model?.message}
                      />
                    )}
                  />
                  <Box sx={{ fontSize: 12, color: 'error.main', mt: 1 }}>
                    {getLocaleMessage('modelNameHint', mergedConfig.locale.language)}
                  </Box>
                </>
              ) : modelUserList.length === 0 ? (
                <Button
                  fullWidth
                  variant="outlined"
                  disabled={modelLoading}
                  sx={{ mt: 4 }}
                  onClick={handleSubmit(getModel)}
                >
                  {modelLoading
                    ? '加载中...'
                    : getLocaleMessage('getModelList', mergedConfig.locale.language)}
                </Button>
              ) : (
                <>
                  <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                    {getLocaleMessage('modelName', mergedConfig.locale.language)}{' '}
                    <Box component="span" sx={{ color: 'red' }}>
                      *
                    </Box>
                  </Box>
                  <Controller
                    control={control}
                    name="model"
                    render={({ field }) => (
                      <TextField
                        {...field}
                        fullWidth
                        select
                        size="small"
                        placeholder=""
                        error={!!errors.model}
                        helperText={errors.model?.message}
                      >
                        {modelUserList.map((it) => (
                          <MenuItem key={it.model} value={it.model}>
                            {it.model}
                          </MenuItem>
                        ))}
                      </TextField>
                    )}
                  />

                  {/* Custom Header */}
                  {filteredProviders[providerBrand]?.customHeader &&
                    mergedConfig.features.enableHeaderConfig && (
                      <>
                        <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                          {getLocaleMessage('header', mergedConfig.locale.language)}
                        </Box>
                        <Stack direction="row" gap={1}>
                          <Controller
                            control={control}
                            name="api_header_key"
                            render={({ field }) => (
                              <TextField
                                {...field}
                                fullWidth
                                size="small"
                                placeholder="key"
                                error={!!errors.api_header_key}
                                helperText={errors.api_header_key?.message}
                              />
                            )}
                          />
                          <Box sx={{ fontSize: 14, lineHeight: '36px' }}>=</Box>
                          <Controller
                            control={control}
                            name="api_header_value"
                            render={({ field }) => (
                              <TextField
                                {...field}
                                fullWidth
                                size="small"
                                placeholder="value"
                                error={!!errors.api_header_value}
                                helperText={errors.api_header_value?.message}
                              />
                            )}
                          />
                        </Stack>
                      </>
                    )}
                </>
              )}

              {/* Error Display */}
              {error && (
                <SimpleCard
                  sx={{
                    color: 'error.main',
                    mt: 2,
                    fontSize: 12,
                    p: 2,
                    bgcolor: 'background.paper2',
                    border: '1px solid',
                    borderColor: 'error.main',
                  }}
                >
                  {error}
                </SimpleCard>
              )}
            </Box>
          </Stack>
        </Box>

        {/* Footer */}
        <Box
          sx={{
            p: 3,
            borderTop: '1px solid',
            borderColor: 'divider',
            display: 'flex',
            justifyContent: 'flex-end',
            gap: 2,
          }}
        >
          <Button onClick={handleReset} variant="outlined">
            {getLocaleMessage('cancel', mergedConfig.locale.language)}
          </Button>
          <Button
            variant="contained"
            loading={loading}
            disabled={!success && providerBrand !== 'Other'}
            onClick={handleSubmit(onSubmit)}
          >
            {getLocaleMessage('save', mergedConfig.locale.language)}
          </Button>
        </Box>
      </Box>
    </Box>
  );
}; 