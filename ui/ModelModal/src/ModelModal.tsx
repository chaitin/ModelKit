import {
  Box,
  Button,
  MenuItem,
  Stack,
  TextField,
  useTheme,
  alpha as addOpacityToColor,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Checkbox,
  FormControlLabel,
  ListSubheader,
  InputAdornment,
  IconButton,
} from '@mui/material';
import { Visibility, VisibilityOff, Search } from '@mui/icons-material';
import { Icon, message, Modal, ThemeProvider } from '@ctzhian/ui';
import Card from './components/card';
import ModelTagsWithLabel from './components/ModelTagsWithLabel';
import ModelTagFilter from './components/ModelTagFilter';
import React, { useEffect, useState } from 'react';
import { useForm, Controller } from 'react-hook-form';
import { AddModelForm, Model, ModelModalProps } from './types/types';
import { DEFAULT_MODEL_PROVIDERS } from './constants/providers';
import { getLocaleMessage } from './constants/locale';
import './assets/fonts/iconfont';
import { lightTheme } from './theme';
import { isValidURL } from './utils';
import { getModelGroup, getModelLogo } from './utils/model';

const titleMap: Record<string, string> = {
  ['llm']: '对话模型',
  ['chat']: '对话模型',
  ['coder']: '代码补全模型',
  ['code']: '代码补全模型',
  ['embedding']: '向量模型',
  ['rerank']: '重排序模型',
  ['reranker']: '重排序模型',
  ['analysis']: '分析模型',
  ['analysis-vl']: '图像分析模型',
};

export const ModelModal: React.FC<ModelModalProps> = ({
  open,
  onClose,
  refresh,
  data,
  model_type = 'llm',
  modelService,
  language = 'zh-CN',
  messageComponent,
  is_close_model_remark = false,
  addingModelTutorialURL = 'https://github.com/chaitin/ModelKit/blob/main/docs/AddinModelTutorial.md',
  beforeSubmit,
  onOk,
}: ModelModalProps) => {
  const theme = useTheme();

  // 消息处理器，优先使用传入的messageComponent，否则使用@ctzhian/ui的message
  const messageHandler = messageComponent || message;

  const providers: Record<string, any> = DEFAULT_MODEL_PROVIDERS;

  const {
    formState: { errors },
    handleSubmit,
    control,
    reset,
    setValue,
    watch,
  } = useForm<AddModelForm>({
    defaultValues: {
      model_type,
      provider: 'BaiZhiCloud',
      base_url: providers['BaiZhiCloud'].defaultBaseUrl,
      model_name: '',
      api_version: '',
      api_key: '',
      api_header_key: '',
      api_header_value: '',
      show_name: '',
      // 高级设置默认值
      context_window_size: 64000,
      max_output_tokens: 8192,
      enable_r1_params: false,
      support_image: model_type === 'analysis-vl' ? true : false,
      support_compute: false,
      support_prompt_caching: false,
    },
  });

  const providerBrand = watch('provider');
  const baseUrl = watch('base_url');

  const [modelUserList, setModelUserList] = useState<{ model: string }[]>([]);
  const [filteredModelList, setFilteredModelList] = useState<
    { model: string; provider: string }[]
  >([]);
  const [modelSearchQuery, setModelSearchQuery] = useState('');

  const [loading, setLoading] = useState(false);
  const [modelLoading, setModelLoading] = useState(false);
  const [error, setError] = useState('');
  const [addModelError, setAddModelError] = useState('');
  const [success, setSuccess] = useState(false);
  const [expandAdvanced, setExpandAdvanced] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const handleReset = () => {
    onClose();
    reset({
      model_type,
      provider: 'BaiZhiCloud',
      model_name: '',
      base_url: '',
      api_key: '',
      api_version: '',
      api_header_key: '',
      api_header_value: '',
      resource_name: '',
      // 重置高级设置
      context_window_size: 64000,
      max_output_tokens: 8192,
      enable_r1_params: false,
      support_image: model_type === 'analysis-vl' ? true : false,
      support_compute: false,
      support_prompt_caching: false,
    });
    setModelUserList([]);
    setFilteredModelList([]);
    setSuccess(false);
    setLoading(false);
    setModelLoading(false);
    setError('');
    setAddModelError('');
    // 重置高级设置的展开状态
    setExpandAdvanced(false);
    refresh();
  };

  // 获取处理后的URL（与灰色显示逻辑一致）
  const getProcessedUrl = (baseUrl: string, provider: string) => {
    if (!providers[provider].urlWrite) {
      return baseUrl;
    }
    if (baseUrl.endsWith('#')) {
      return baseUrl;
    }
    const forceUseOriginalHost = () => {
      if (baseUrl.endsWith('/')) {
        baseUrl = baseUrl.slice(0, -1);
        return true;
      }
      if (/\/v\d+$/.test(baseUrl)) {
        return true;
      }
      return baseUrl.endsWith('volces.com/api/v3');
    };

    return forceUseOriginalHost() ? baseUrl : `${baseUrl}/v1`;
  };

  // 从Azure OpenAI base_url中提取resource_name
  const extractResourceNameFromUrl = (baseUrl: string): string => {
    if (!baseUrl) return '';
    // 匹配 https://<resource_name>.openai.azure.com 格式
    const match = baseUrl.match(/https:\/\/([^.]+)\.openai\.azure\.com/);
    return match ? match[1] : '';
  };

  const getModel = (value: AddModelForm) => {
    let header = '';
    if (value.api_header_key && value.api_header_value) {
      header = value.api_header_key + '=' + value.api_header_value;
    }
    setAddModelError('');
    setModelLoading(true);
    modelService
      .listModel({
        model_type,
        api_key: value.api_key,
        base_url: getProcessedUrl(value.base_url, value.provider),
        provider: value.provider as Exclude<typeof value.provider, 'Other'>,
        api_header: value.api_header || header,
      })
      .then((res) => {
        if (res.error) {
          messageHandler.error('获取模型列表失败');
          setAddModelError(res.error);
          setModelLoading(false);
        } else {
          setModelUserList(
            (res.models || [])
              .filter((item): item is { model: string } => !!item.model)
              .sort((a, b) => a.model!.localeCompare(b.model!)),
          );
          if (
            data &&
            (res.models || []).find((it) => it.model === data.model_name)
          ) {
            setValue('model_name', data.model_name!);
          } else {
            setValue('model_name', res.models?.[0]?.model || '');
          }
          setSuccess(true);
          setAddModelError('');
        }
      })
      .finally(() => {
        setModelLoading(false);
      })
      .catch((res) => {
        setModelLoading(false);
      });
  };

  const onSubmit = (value: AddModelForm) => {
    let header = '';
    if (value.api_header_key && value.api_header_value) {
      header = value.api_header_key + '=' + value.api_header_value;
    }
    setError('');
    setAddModelError('');
    setLoading(true);
    modelService
      .checkModel({
        model_type,
        model_name: value.model_name,
        api_key: value.api_key,
        base_url: getProcessedUrl(value.base_url, value.provider),
        api_version: value.api_version,
        provider: value.provider,
        api_header: value.api_header || header,
        param: {
          context_window: value.context_window_size,
          max_tokens: value.max_output_tokens,
          r1_enabled: value.enable_r1_params,
          support_images: value.support_image,
          support_computer_use: value.support_compute,
          support_prompt_cache: value.support_prompt_caching,
        },
      })
      .then((res) => {
        // 错误处理
        if (res.error) {
          messageHandler.error('模型检查失败');
          setAddModelError(res.error);
          setLoading(false);
        } else if (data) {
          modelService
            .updateModel({
              api_key: value.api_key,
              model_type,
              base_url: getProcessedUrl(value.base_url, value.provider),
              model_name: value.model_name,
              api_header: value.api_header || header,
              api_version: value.api_version,
              id: (data as any).id || '',
              provider: value.provider as Exclude<
                typeof value.provider,
                'Other'
              >,
              show_name: value.show_name,
              // 添加高级设置字段到 param 对象中
              param: {
                context_window: value.context_window_size,
                max_tokens: value.max_output_tokens,
                r1_enabled: value.enable_r1_params,
                support_images: value.support_image,
                support_computer_use: value.support_compute,
                support_prompt_cache: value.support_prompt_caching,
              },
            })
            .then((res) => {
              if (res.error) {
                messageHandler.error('修改模型失败');
                setLoading(false);
              } else {
                messageHandler.success('修改成功');
                handleReset();
              }
            })
            .finally(() => {
              setLoading(false);
            })
            .catch((res) => {
              messageHandler.error('修改模型失败');
              setLoading(false);
            });
        } else {
          modelService
            .createModel({
              model_type,
              api_key: value.api_key,
              base_url: getProcessedUrl(value.base_url, value.provider),
              model_name: value.model_name,
              api_header: value.api_header || header,
              provider: value.provider as Exclude<
                typeof value.provider,
                'Other'
              >,
              show_name: value.show_name,
              // 添加高级设置字段到 param 对象中
              param: {
                context_window: value.context_window_size,
                max_tokens: value.max_output_tokens,
                r1_enabled: value.enable_r1_params,
                support_images: value.support_image,
                support_computer_use: value.support_compute,
                support_prompt_cache: value.support_prompt_caching,
              },
            })
            .then((res) => {
              if (res.error) {
                messageHandler.error('添加模型失败');
                setLoading(false);
              } else {
                messageHandler.success('添加成功');
                handleReset();
              }
            })
            .finally(() => {
              setLoading(false);
            })
            .catch((res) => {
              messageHandler.error('添加模型失败');
              setLoading(false);
            });
        }
      })
      .catch((res) => {
        messageHandler.error('检查模型失败');
        setLoading(false);
      });
  };

  // 处理提交前的验证逻辑
  const handleOk = async (value: AddModelForm) => {
    // 如果有 beforeSubmit，先执行验证
    if (beforeSubmit) {
      try {
        // 调用 beforeSubmit，支持同步和异步返回值
        const result = beforeSubmit(value);

        // 如果返回的是 Promise，等待其 resolve
        const shouldSubmit = result instanceof Promise ? await result : result;

        // 如果 beforeSubmit 返回 false，不继续执行
        if (!shouldSubmit) {
          return;
        }
      } catch (error) {
        // 如果 beforeSubmit 抛出异常或 Promise reject，不继续执行
        console.error('beforeSubmit error:', error);
        return;
      }
    }

    // beforeSubmit 验证通过后，执行 onOk 或 onSubmit
    if (onOk) {
      onOk(value);
    } else {
      onSubmit(value);
    }
  };

  const resetCurData = (value: Model) => {
    // @ts-ignore
    if (
      value.provider &&
      !['Other', 'AzureOpenAI', 'Volcengine'].includes(value.provider)
    ) {
      getModel({
        api_key: value.api_key || '',
        base_url: value.base_url || '',
        model_name: value.model_name || '',
        provider: value.provider,
        api_version: value.api_version || '',
        api_header_key: value.api_header?.split('=')[0] || '',
        api_header: value.api_header || '',
        api_header_value: value.api_header?.split('=')[1] || '',
        model_type,
        show_name: value.show_name || '',
        resource_name:
          value.provider === 'AzureOpenAI'
            ? extractResourceNameFromUrl(value.base_url || '')
            : '',
        context_window_size: 64000,
        max_output_tokens: 8192,
        enable_r1_params: false,
        support_image: false,
        support_compute: false,
        support_prompt_caching: false,
      });
    }
    reset({
      model_type,
      provider: value.provider || 'Other',
      model_name: value.model_name || '',
      base_url: value.base_url || '',
      api_key: value.api_key || '',
      api_version: value.api_version || '',
      api_header_key: value.api_header?.split('=')[0] || '',
      api_header_value: value.api_header?.split('=')[1] || '',
      show_name: value.show_name || '',
      resource_name:
        value.provider === 'AzureOpenAI'
          ? extractResourceNameFromUrl(value.base_url || '')
          : '',
      context_window_size: value.param?.context_window || 64000,
      max_output_tokens: value.param?.max_tokens || 8192,
      enable_r1_params: value.param?.r1_enabled || false,
      support_image:
        model_type === 'analysis-vl'
          ? true
          : value.param?.support_images || false,
      support_compute: value.param?.support_computer_use || false,
      support_prompt_caching: value.param?.support_prompt_cache || false,
    });
  };

  useEffect(() => {
    if (open) {
      if (data) {
        resetCurData(data);
      } else {
        reset({
          model_type,
          provider: 'BaiZhiCloud',
          model_name: '',
          base_url: providers['BaiZhiCloud'].defaultBaseUrl,
          api_key: '',
          api_version: '',
          api_header_key: '',
          api_header_value: '',
          show_name: '',
          resource_name: '',
          // 高级设置默认值
          context_window_size: 64000,
          max_output_tokens: 8192,
          enable_r1_params: false,
          support_image: model_type === 'analysis-vl' ? true : false,
          support_compute: false,
          support_prompt_caching: false,
        });
      }
      // 确保每次打开时高级设置都是折叠的
      setExpandAdvanced(false);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [data, open]);

  return (
    <ThemeProvider theme={lightTheme}>
      <Modal
        title={
          data ? `修改${titleMap[model_type]}` : `添加${titleMap[model_type]}`
        }
        open={open}
        width={800}
        onCancel={handleReset}
        cancelText={getLocaleMessage('cancel', language)}
        okText={getLocaleMessage('save', language)}
        onOk={handleSubmit(handleOk)}
        okButtonProps={{
          loading,
          disabled:
            !success &&
            !['Other', 'AzureOpenAI', 'Volcengine'].includes(providerBrand) &&
            !(
              providerBrand === 'BaiLian' &&
              (model_type === 'embedding' ||
                model_type === 'rerank' ||
                model_type === 'reranker')
            ),
        }}
      >
        <Stack
          direction={'row'}
          alignItems={'stretch'}
          gap={3}
          sx={{ height: 500 }}
        >
          <Stack
            gap={1}
            sx={{
              width: 200,
              flexShrink: 0,
              bgcolor: 'rgb(248, 249, 250)',
              borderRadius: '10px',
              p: 1,
              height: '100%',
              display: 'flex',
              flexDirection: 'column',
            }}
          >
            <Box
              sx={{
                fontSize: 14,
                lineHeight: '24px',
                fontWeight: 'bold',
                p: 1,
                flexShrink: 0,
              }}
            >
              模型供应商
            </Box>
            <Box
              sx={{
                flex: 1,
                overflowY: 'auto',
                '&::-webkit-scrollbar': {
                  width: '6px',
                },
                '&::-webkit-scrollbar-track': {
                  background: 'transparent',
                },
                '&::-webkit-scrollbar-thumb': {
                  background: '#d0d0d0',
                  borderRadius: '3px',
                },
                '&::-webkit-scrollbar-thumb:hover': {
                  background: '#a0a0a0',
                },
              }}
            >
              <Stack gap={1}>
                {Object.values(providers)
                  .filter((it) => {
                    // 根据model_type和provider配置决定是否显示
                    switch (model_type) {
                      case 'chat':
                      case 'llm':
                        return it.chat;
                      case 'code':
                      case 'coder':
                        return it.code;
                      case 'embedding':
                        return it.embedding;
                      case 'rerank':
                      case 'reranker':
                        return it.rerank;
                      case 'analysis':
                        return it.analysis;
                      case 'analysis-vl':
                        return it.analysis_vl;
                      default:
                        return (
                          it.label === 'BaiZhiCloud' || it.label === 'Other'
                        );
                    }
                  })
                  .map((it) => (
                    <Stack
                      direction={'row'}
                      alignItems={'center'}
                      gap={1.5}
                      key={it.label}
                      sx={{
                        cursor: 'pointer',
                        fontSize: 14,
                        lineHeight: '24px',
                        p: 1,
                        borderRadius: '10px',
                        fontWeight: 'bold',
                        fontFamily: 'Gbold',
                        ...(providerBrand === it.label && {
                          bgcolor: addOpacityToColor(
                            theme.palette.primary.main,
                            0.1,
                          ),
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
                          setAddModelError('');
                          setModelLoading(false);
                          setSuccess(false);
                          reset({
                            provider:
                              it.label as keyof typeof DEFAULT_MODEL_PROVIDERS,
                            base_url: (() => {
                              if (it.label === 'AzureOpenAI') return '';
                              if (it.label === 'BaiLian') {
                                if (model_type === 'embedding') {
                                  return 'https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding';
                                }
                                if (model_type === 'rerank' || model_type === 'reranker') {
                                  return 'https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank#';
                                }
                              }
                              return it.defaultBaseUrl;
                            })(),
                            model_name: '',
                            api_version: '',
                            api_key: '',
                            api_header_key: '',
                            api_header_value: '',
                            show_name: '',
                            resource_name: '',
                            // 重置高级设置
                            context_window_size: 64000,
                            max_output_tokens: 8192,
                            enable_r1_params: false,
                            support_image:
                              model_type === 'analysis-vl' ? true : false,
                            support_compute: false,
                            support_prompt_caching: false,
                          });
                        }
                      }}
                    >
                      <Icon type={it.icon} sx={{ fontSize: 18 }} />
                      {it.cn || it.label || '其他'}
                    </Stack>
                  ))}
              </Stack>
            </Box>
          </Stack>
          <Box
            sx={{
              flex: 1,
              height: '100%',
              overflowY: 'auto',
              '&::-webkit-scrollbar': {
                width: '6px',
              },
              '&::-webkit-scrollbar-track': {
                background: 'transparent',
              },
              '&::-webkit-scrollbar-thumb': {
                background: '#d0d0d0',
                borderRadius: '3px',
              },
              '&::-webkit-scrollbar-thumb:hover': {
                background: '#a0a0a0',
              },
            }}
          >
            <Stack
              direction={'row'}
              alignItems={'center'}
              justifyContent={'space-between'}
              sx={{ fontSize: 14, lineHeight: '32px' }}
            >
              <Box>
                API 地址{' '}
                <Box component={'span'} sx={{ color: 'red' }}>
                  *
                </Box>
              </Box>
              {providers[providerBrand].addingModelTutorial && (
                <Box
                  component={'span'}
                  sx={{
                    color: 'info.main',
                    cursor: 'pointer',
                    ml: 1,
                    textAlign: 'right',
                  }}
                  onClick={() => window.open(addingModelTutorialURL, '_blank')}
                >
                  添加模型教程
                </Box>
              )}
            </Stack>
            <Controller
              control={control}
              name='base_url'
              rules={{
                required: {
                  value: true,
                  message: 'URL 不能为空',
                },
                validate: (value) => {
                  if (
                    !providers[providerBrand].urlWrite ||
                    providerBrand === 'AzureOpenAI'
                  ) {
                    return true;
                  }
                  const res = isValidURL(value);
                  return res === '' || res;
                },
              }}
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  disabled={!providers[providerBrand].urlWrite}
                  size='small'
                  placeholder={providers[providerBrand].defaultBaseUrl}
                  error={!!errors.base_url}
                  helperText={errors.base_url?.message}
                  onChange={(e) => {
                    field.onChange(e.target.value);
                    setModelUserList([]);
                    // 只有当供应商不是"Other"时才清空模型名称
                    if (providerBrand !== 'Other') {
                      setValue('model_name', '');
                    }
                    setSuccess(false);
                    setAddModelError('');
                  }}
                />
              )}
            />
            {baseUrl && providers[providerBrand].urlWrite && (
              <Stack
                direction={'row'}
                alignItems={'center'}
                justifyContent={'space-between'}
                sx={{
                  fontSize: 12,
                  color: 'text.secondary',
                  mt: 0.5,
                }}
              >
                <Box sx={{ wordBreak: 'break-all', flex: 1 }}>
                  {/* 根据模型类型显示不同的URL后缀：embedding显示/embeddings，rerank显示/rerank，其他显示/chat/completions */}
                  {baseUrl &&
                    providers[providerBrand].urlWrite &&
                    (() => {
                      const processedUrl = getProcessedUrl(
                        baseUrl,
                        providerBrand,
                      );
                      if (baseUrl.endsWith('#')) {
                        return processedUrl;
                      }
                      // 根据模型类型添加不同的后缀
                      if (model_type === 'embedding') {
                        return `${processedUrl}/embeddings`;
                      } else if (
                        model_type === 'rerank' ||
                        model_type === 'reranker'
                      ) {
                        return `${processedUrl}/rerank`;
                      } else {
                        return `${processedUrl}/chat/completions`;
                      }
                    })()}
                </Box>
                <Box sx={{ ml: 2, flexShrink: 0, fontSize: 10, opacity: 0.7 }}>
                  /结尾忽略V1版本，#结尾强制使用输入地址
                </Box>
              </Stack>
            )}
            <Stack
              direction={'row'}
              alignItems={'center'}
              justifyContent={'space-between'}
              sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}
            >
              <Box>
                API Key
                {providers[providerBrand].secretRequired && (
                  <Box component={'span'} sx={{ color: 'red' }}>
                    {' '}
                    *
                  </Box>
                )}
              </Box>
              {providers[providerBrand].modelDocumentUrl && (
                <Box
                  component={'span'}
                  sx={{
                    color: 'info.main',
                    cursor: 'pointer',
                    ml: 1,
                    textAlign: 'right',
                  }}
                  onClick={() =>
                    window.open(
                      providers[providerBrand].modelDocumentUrl,
                      '_blank',
                    )
                  }
                >
                  获取API Key
                </Box>
              )}
            </Stack>
            <Controller
              control={control}
              name='api_key'
              rules={{
                required: {
                  value: providers[providerBrand].secretRequired,
                  message: 'API Key 不能为空',
                },
              }}
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  size='small'
                  type={showPassword ? 'text' : 'password'}
                  placeholder=''
                  error={!!errors.api_key}
                  helperText={errors.api_key?.message}
                  InputProps={{
                    endAdornment: (
                      <InputAdornment position='end'>
                        <IconButton
                          aria-label='toggle password visibility'
                          onClick={() => setShowPassword(!showPassword)}
                          edge='end'
                        >
                          {showPassword ? <VisibilityOff /> : <Visibility />}
                        </IconButton>
                      </InputAdornment>
                    ),
                  }}
                  onChange={(e) => {
                    field.onChange(e.target.value);
                    setModelUserList([]);
                    // 只有当供应商不是"Other"时才清空模型名称
                    if (providerBrand !== 'Other') {
                      setValue('model_name', '');
                    }
                    setSuccess(false);
                    setAddModelError('');
                  }}
                />
              )}
            />
            {(modelUserList.length !== 0 || providerBrand === 'Other') &&
              !is_close_model_remark && (
                <>
                  <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                    模型备注
                    <Box component={'span'} sx={{ color: 'red' }}>
                      *
                    </Box>
                  </Box>
                  <Controller
                    control={control}
                    name='show_name'
                    rules={{
                      required: {
                        value: !is_close_model_remark,
                        message: '模型备注不能为空',
                      },
                    }}
                    render={({ field }) => (
                      <TextField
                        {...field}
                        fullWidth
                        size='small'
                        placeholder=''
                        error={!!errors.show_name}
                        helperText={errors.show_name?.message}
                      />
                    )}
                  />
                </>
              )}
            {providerBrand === 'AzureOpenAI' && (
              <>
                <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                  Resource Name
                  <Box component={'span'} sx={{ color: 'red' }}>
                    *
                  </Box>
                </Box>
                <Controller
                  control={control}
                  name='resource_name'
                  rules={{
                    required: {
                      value: true,
                      message: 'Resource Name 不能为空',
                    },
                  }}
                  render={({ field }) => (
                    <TextField
                      {...field}
                      fullWidth
                      size='small'
                      error={!!errors.resource_name}
                      helperText={errors.resource_name?.message}
                      onChange={(e) => {
                        field.onChange(e.target.value);
                        // 动态更新base_url
                        const resourceName = e.target.value;
                        if (resourceName) {
                          setValue(
                            'base_url',
                            `https://${resourceName}.openai.azure.com`,
                          );
                        } else {
                          setValue('base_url', '');
                        }
                        setModelUserList([]);
                        setValue('model_name', '');
                        setSuccess(false);
                        setAddModelError('');
                      }}
                    />
                  )}
                />
                <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                  API Version
                </Box>
                <Controller
                  control={control}
                  name='api_version'
                  render={({ field }) => (
                    <TextField
                      {...field}
                      fullWidth
                      size='small'
                      placeholder='2024-10-21'
                      error={!!errors.api_version}
                      helperText={errors.api_version?.message}
                      onChange={(e) => {
                        field.onChange(e.target.value);
                        setModelUserList([]);
                        setValue('model_name', '');
                        setSuccess(false);
                        setAddModelError('');
                      }}
                    />
                  )}
                />
              </>
            )}
            {['Other', 'AzureOpenAI', 'Volcengine'].includes(providerBrand) ||
            (providerBrand === 'BaiLian' &&
              (model_type === 'embedding' ||
                model_type === 'rerank' ||
                model_type === 'reranker')) ? (
              <>
                <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                  模型名称{' '}
                  <Box component={'span'} sx={{ color: 'red' }}>
                    *
                  </Box>
                </Box>
                <Controller
                  control={control}
                  name='model_name'
                  rules={{
                    required: '模型名称不能为空',
                  }}
                  render={({ field }) => (
                    <TextField
                      {...field}
                      fullWidth
                      size='small'
                      placeholder=''
                      error={!!errors.model_name}
                      helperText={errors.model_name?.message}
                    />
                  )}
                />
                <Box sx={{ fontSize: 12, color: 'error.main', mt: 1 }}>
                  需要与模型供应商提供的名称完全一致
                </Box>
              </>
            ) : modelUserList.length === 0 ? (
              <Button
                fullWidth
                variant='outlined'
                loading={modelLoading}
                sx={{
                  mt: 4,
                  borderRadius: '10px',
                  boxShadow: 'none',
                  fontFamily: `var(--font-gilory), var(--font-HarmonyOS), 'PingFang SC', 'Roboto', 'Helvetica', 'Arial', sans-serif`,
                  color: 'black',
                  borderColor: 'black',
                }}
                onClick={handleSubmit(getModel)}
              >
                获取模型列表
              </Button>
            ) : (
              <>
                <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                  模型名称{' '}
                  <Box component={'span'} sx={{ color: 'red' }}>
                    *
                  </Box>
                </Box>
                <Controller
                  control={control}
                  name='model_name'
                  render={({ field }) => (
                    <TextField
                      {...field}
                      fullWidth
                      select
                      size='small'
                      placeholder=''
                      error={!!errors.model_name}
                      helperText={errors.model_name?.message}
                      SelectProps={{
                        MenuProps: {
                          autoFocus: false,
                          MenuListProps: {
                            autoFocusItem: false,
                          },
                          PaperProps: {
                            sx: {
                              maxHeight: 450,
                              '& .MuiList-root': {
                                paddingTop: 0,
                              },
                            },
                          },
                        },
                      }}
                    >
                      {(() => {
                        // 使用筛选后的模型列表，如果没有筛选则使用原始列表
                        const modelsBase =
                          filteredModelList.length > 0
                            ? filteredModelList
                            : modelUserList.map((item) => ({
                              model: item.model,
                              provider: providerBrand,
                            }));

                        const query = modelSearchQuery.trim().toLowerCase();
                        const modelsToShow = query
                          ? modelsBase.filter((m) =>
                            m.model.toLowerCase().includes(query),
                          )
                          : modelsBase;

                        // 按组分类模型
                        const groupedModels = modelsToShow.reduce(
                          (acc, model) => {
                            const group = getModelGroup(model.model);
                            if (!acc[group]) {
                              acc[group] = [];
                            }
                            acc[group].push(model);
                            return acc;
                          },
                          {} as Record<string, typeof modelsToShow>,
                        );

                        // 渲染分组后的模型
                        const modelItems = Object.entries(groupedModels)
                          .map(([group, models]) => [
                            <ListSubheader
                              key={`header-${group}`}
                              sx={{
                                backgroundColor: 'transparent',
                                fontWeight: 'bold',
                                position: 'static',
                              }}
                            >
                              {group}
                            </ListSubheader>,
                            ...models.map((it) => (
                              <MenuItem key={it.model} value={it.model}>
                                <Box
                                  sx={{
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'space-between',
                                    width: '100%',
                                  }}
                                >
                                  <Box
                                    sx={{
                                      display: 'flex',
                                      alignItems: 'center',
                                      flex: 1,
                                    }}
                                  >
                                    {getModelLogo(it.model) && (
                                      <Box
                                        component='img'
                                        src={getModelLogo(it.model)}
                                        alt={`${it.model} logo`}
                                        sx={{
                                          width: 20,
                                          height: 20,
                                          mr: 1,
                                          borderRadius: '50%',
                                        }}
                                      />
                                    )}
                                    {it.model}
                                  </Box>
                                  <ModelTagsWithLabel
                                    model_id={it.model}
                                    provider={providerBrand}
                                    size={10}
                                    showLabel={false}
                                    showTooltip={false}
                                  />
                                </Box>
                              </MenuItem>
                            )),
                          ])
                          .flat();

                        return [
                          <ListSubheader
                            key='sticky-tools'
                            sx={{
                              position: 'sticky',
                              top: 0,
                              zIndex: 1001,
                              backgroundColor: '#ffffff !important',
                              borderBottom: '1px solid',
                              borderColor: 'divider',
                              boxShadow: '0 2px 4px rgba(0,0,0,0.06)',
                              cursor: 'default',
                              px: 1,
                              py: 1,
                            }}
                            onMouseDown={(e) => {
                              e.stopPropagation();
                            }}
                            onClick={(e) => {
                              e.stopPropagation();
                            }}
                          >
                            <Stack direction='column' spacing={1}>
                              <TextField
                                value={modelSearchQuery}
                                onChange={(e) =>
                                  setModelSearchQuery(e.target.value)
                                }
                                fullWidth
                                size='small'
                                placeholder='搜索模型名称'
                                onKeyDown={(e) => {
                                  e.stopPropagation();
                                }}
                                onMouseDown={(e) => {
                                  e.stopPropagation();
                                }}
                                onClick={(e) => {
                                  e.stopPropagation();
                                }}
                                InputProps={{
                                  startAdornment: (
                                    <InputAdornment position='start'>
                                      <Search fontSize='small' />
                                    </InputAdornment>
                                  ),
                                }}
                              />
                              <Box
                                onMouseDown={(e) => {
                                  e.stopPropagation();
                                }}
                                onClick={(e) => {
                                  e.stopPropagation();
                                }}
                              >
                                <ModelTagFilter
                                  models={modelUserList.map((item) => ({
                                    model: item.model,
                                    provider: providerBrand,
                                  }))}
                                  onFilteredModelsChange={(filteredModels) => {
                                    setFilteredModelList(filteredModels);
                                  }}
                                />
                              </Box>
                            </Stack>
                          </ListSubheader>,
                          ...modelItems,
                        ];
                      })()}
                    </TextField>
                  )}
                />
                {providers[providerBrand].customHeader && (
                  <>
                    <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                      Header
                    </Box>
                    <Stack direction={'row'} gap={1}>
                      <Controller
                        control={control}
                        name='api_header_key'
                        render={({ field }) => (
                          <TextField
                            {...field}
                            fullWidth
                            size='small'
                            placeholder='key'
                            error={!!errors.api_header_key}
                            helperText={errors.api_header_key?.message}
                          />
                        )}
                      />
                      <Box sx={{ fontSize: 14, lineHeight: '36px' }}>=</Box>
                      <Controller
                        control={control}
                        name='api_header_value'
                        render={({ field }) => (
                          <TextField
                            {...field}
                            fullWidth
                            size='small'
                            placeholder='value'
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
            {/* 高级设置部分 - 在选择了模型或者是其它供应商时显示，但不包括embedding、rerank、reranker类型 */}
            {(modelUserList.length !== 0 || providerBrand === 'Other') &&
              !['embedding', 'rerank', 'reranker'].includes(model_type) && (
                <Box sx={{ mt: 2 }}>
                  <Accordion
                    sx={{
                      boxShadow: 'none',
                      bgcolor: 'transparent',
                      '&:before': {
                        display: 'none',
                      },
                      '& .MuiAccordionSummary-root': {
                        padding: 0,
                        minHeight: 'auto',
                        '& .MuiAccordionSummary-content': {
                          margin: 0,
                        },
                      },
                      '& .MuiAccordionDetails-root': {
                        padding: 0,
                      },
                    }}
                    expanded={expandAdvanced}
                    onChange={() => setExpandAdvanced(!expandAdvanced)}
                  >
                    <AccordionSummary
                      sx={{
                        fontSize: 14,
                        fontWeight: 500,
                        lineHeight: '32px',
                        color: 'primary.main',
                        display: 'flex',
                        alignItems: 'center',
                        cursor: 'pointer',
                        '&:hover': {
                          opacity: 0.8,
                        },
                      }}
                    >
                      高级设置
                    </AccordionSummary>
                    <AccordionDetails>
                      <Stack spacing={0}>
                        {/* 复选框组 - 使用更紧凑的布局 */}
                        <Stack spacing={0} sx={{ ml: -1.2 }}>
                          <Controller
                            control={control}
                            name='support_image'
                            render={({ field }) => {
                              const isAnalysisVl = model_type === 'analysis-vl';
                              const isChecked = isAnalysisVl
                                ? true
                                : field.value;

                              return (
                                <FormControlLabel
                                  control={
                                    <Checkbox
                                      checked={isChecked}
                                      onChange={(e) => {
                                        if (!isAnalysisVl) {
                                          field.onChange(e.target.checked);
                                        }
                                      }}
                                      disabled={isAnalysisVl}
                                      size='small'
                                    />
                                  }
                                  label={
                                    <Box sx={{ fontSize: 12 }}>
                                      启用图片
                                      <Box
                                        component='span'
                                        sx={{
                                          ml: 1,
                                          color: 'text.secondary',
                                          fontSize: 11,
                                        }}
                                      >
                                        {isAnalysisVl
                                          ? '(图像分析模型默认启用图片功能)'
                                          : '(支持图片输入的模型可以启用此选项)'}
                                      </Box>
                                    </Box>
                                  }
                                  sx={{ margin: 0 }}
                                />
                              );
                            }}
                          />
                          <Controller
                            control={control}
                            name='enable_r1_params'
                            render={({ field }) => (
                              <FormControlLabel
                                control={
                                  <Checkbox
                                    checked={field.value}
                                    onChange={(e) =>
                                      field.onChange(e.target.checked)
                                    }
                                    size='small'
                                  />
                                }
                                label={
                                  <Box sx={{ fontSize: 12 }}>
                                    启用 R1 模型参数
                                    <Box
                                      component='span'
                                      sx={{
                                        ml: 1,
                                        color: 'text.secondary',
                                        fontSize: 11,
                                      }}
                                    >
                                      (使用 QWQ 等 R1
                                      系列模型时必须启用，避免出现 400 错误)
                                    </Box>
                                  </Box>
                                }
                                sx={{ margin: 0 }}
                              />
                            )}
                          />
                        </Stack>
                        <Box>
                          <Box sx={{ fontSize: 14, lineHeight: '32px' }}>
                            上下文窗口大小
                          </Box>
                          <Controller
                            control={control}
                            name='context_window_size'
                            render={({ field }) => (
                              <>
                                <TextField
                                  {...field}
                                  fullWidth
                                  size='small'
                                  placeholder='例如：16000'
                                  type='number'
                                  onChange={(e) =>
                                    field.onChange(Number(e.target.value))
                                  }
                                />
                                <Box
                                  sx={{
                                    mt: 1,
                                    display: 'flex',
                                    gap: 1,
                                    alignItems: 'center',
                                  }}
                                >
                                  {[
                                    { label: '128k', value: 128000 },
                                    { label: '256k', value: 256000 },
                                    { label: '512k', value: 512000 },
                                    { label: '1m', value: 1_000_000 },
                                  ].map((option) => (
                                    <Box
                                      key={option.label}
                                      sx={{
                                        fontSize: 12,
                                        color: 'primary.main',
                                        cursor: 'pointer',
                                        padding: '2px 4px',
                                        '&:hover': {
                                          textDecoration: 'underline',
                                        },
                                      }}
                                      onClick={() =>
                                        field.onChange(option.value)
                                      }
                                    >
                                      {option.label}
                                    </Box>
                                  ))}
                                </Box>
                              </>
                            )}
                          />
                        </Box>

                        <Box>
                          <Box sx={{ fontSize: 14, lineHeight: '32px' }}>
                            最大输出 Token
                          </Box>
                          <Controller
                            control={control}
                            name='max_output_tokens'
                            render={({ field }) => (
                              <TextField
                                {...field}
                                fullWidth
                                size='small'
                                placeholder='例如：4000'
                                type='number'
                                onChange={(e) =>
                                  field.onChange(Number(e.target.value))
                                }
                              />
                            )}
                          />
                        </Box>
                      </Stack>
                    </AccordionDetails>
                  </Accordion>
                </Box>
              )}
            {addModelError && (
              <Card
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
                {addModelError}
              </Card>
            )}
            {error && (
              <Card
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
              </Card>
            )}
          </Box>
        </Stack>
      </Modal>
    </ThemeProvider>
  );
};
