import {
  postCheckModel,
  getListModel,
} from '@/api/Model';
import { ConstsModelType, DomainModel, DomainResp } from '@/api/types';

// 扩展 DomainModel 类型以包含高级设置参数
interface ExtendedDomainModel extends DomainModel {
  param?: {
    context_window?: number;
    max_tokens?: number;
    r1_enabled?: boolean;
    support_images?: boolean;
    support_computer_use?: boolean;
    support_prompt_cache?: boolean;
  };
}

import Card from '@/components/card';
import { ModelProvider } from '@/pages/model/constant';
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
} from '@mui/material';
import { message, Modal } from '@c-x/ui';
import { useEffect, useState } from 'react';
// @ts-ignore
import { Controller, useForm } from 'react-hook-form';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
// 引入 iconfont 图标
import '@/assets/fonts/iconfont.js';

// 自定义 Iconfont 图标组件
interface IconfontIconProps {
  type: string;
  style?: React.CSSProperties;
  sx?: any;
}

const IconfontIcon: React.FC<IconfontIconProps> = ({ type, style = {}, sx = {} }) => {
  // 将 sx 属性转换为 style 属性
  const sxStyle = typeof sx === 'object' ? sx : {};
  const defaultStyle: React.CSSProperties = {
    width: '18px',
    height: '18px',
    fontSize: '18px',
    color: 'currentColor',
    display: 'inline-block',
    verticalAlign: 'middle',
    ...style,
    ...sxStyle,
  };

  return (
    <svg style={defaultStyle} viewBox="0 0 1024 1024">
      <use xlinkHref={`#${type}`} />
    </svg>
  );
};

interface AddModelProps {
  open: boolean;
  data: ExtendedDomainModel | null;
  type: ConstsModelType;
  onClose: () => void;
  refresh: () => void;
}

interface AddModelForm {
  provider: keyof typeof ModelProvider;
  model: string;
  base_url: string;
  api_version: string;
  api_key: string;
  api_header_key: string;
  api_header_value: string;
  type: ConstsModelType;
  show_name: string;
  // 高级设置字段
  context_window_size: number;
  max_output_tokens: number;
  enable_r1_params: boolean;
  support_image: boolean;
  support_compute: boolean;
  support_prompt_caching: boolean;
}

const titleMap = {
  [ConstsModelType.ModelTypeChat]: '对话模型',
  [ConstsModelType.ModelTypeCoder]: '代码补全模型',
  [ConstsModelType.ModelTypeEmbedding]: '向量模型',
  [ConstsModelType.ModelTypeReranker]: '重排序模型',
};

const ModelAdd = ({
  open,
  onClose,
  refresh,
  data,
  type = ConstsModelType.ModelTypeChat,
}: AddModelProps) => {
  const theme = useTheme();

  const providers: Record<string, any> = ModelProvider;

  const {
    formState: { errors },
    handleSubmit,
    control,
    reset,
    setValue,
    watch,
  } = useForm<AddModelForm>({
    defaultValues: {
      type: ConstsModelType.ModelTypeChat,
      provider: 'BaiZhiCloud',
      base_url: providers['BaiZhiCloud'].defaultBaseUrl,
      model: '',
      api_version: '',
      api_key: '',
      api_header_key: '',
      api_header_value: '',
      show_name: '',
      // 高级设置默认值
      context_window_size: 64000,
      max_output_tokens: 8192,
      enable_r1_params: false,
      support_image: false,
      support_compute: false,
      support_prompt_caching: false,
    },
  });

  const providerBrand = watch('provider');

  const [modelUserList, setModelUserList] = useState<{ model: string }[]>([]);

  const [loading, setLoading] = useState(false);
  const [modelLoading, setModelLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [checkSuccess, setCheckSuccess] = useState(false);
  const [checkFailed, setCheckFailed] = useState(false);
  const [expandAdvanced, setExpandAdvanced] = useState(false);
  const [modelListInfo, setModelListInfo] = useState<DomainModel[]>([]);

  const handleReset = () => {
    onClose();
    reset({
      type: ConstsModelType.ModelTypeChat,
      provider: 'BaiZhiCloud',
      model: '',
      base_url: '',
      api_key: '',
      api_version: '',
      api_header_key: '',
      api_header_value: '',
      // 重置高级设置
      context_window_size: 64000,
      max_output_tokens: 8192,
      enable_r1_params: false,
      support_image: false,
      support_compute: false,
      support_prompt_caching: false,
    });
    setModelUserList([]);
    setSuccess(false);
    setCheckSuccess(false); // 重置检查成功状态
    setCheckFailed(false); // 重置检查失败状态
    setLoading(false);
    setModelLoading(false);
    setError('');
    setModelListInfo([]);
    // 重置高级设置的展开状态
    setExpandAdvanced(false);
    refresh();
  };

  const getModel = (value: AddModelForm) => {
    setModelLoading(true);
    
    // 调用获取模型列表接口
    getListModel({
      owned_by: value.provider as any, // 提供商
      sub_type: value.type as any, // 模型类型
    })
      .then((res: any) => {
        // 尝试多种数据解析方式
        let models: DomainModel[] = [];
        
        if (res.data && Array.isArray(res.data)) {
          models = res.data;
        } else if (Array.isArray(res)) {
          models = res;
        } else if (res.data) {
          models = [res.data]; // 单个对象包装成数组
        } else {
          models = [];
        }
        
        // 存储模型详细信息用于显示
        setModelListInfo(models);
        
        setModelUserList(
          models
            .filter((item: DomainModel) => !!item.id)
            .map((item: DomainModel) => ({ model: item.id! }))
            .sort((a: { model: string }, b: { model: string }) => a.model.localeCompare(b.model))
        );
        
        // 如果是编辑模式且找到匹配的模型，则选中它
        if (data && models.find((it: DomainModel) => it.id === data.id)) {
          setValue('model', data.id!);
        } else {
          // 否则选择第一个模型
          setValue('model', models[0]?.id || '');
        }
        setSuccess(true);
      })
      .catch((error) => {
        message.error('获取模型列表失败');
        setModelUserList([]);
        setModelListInfo([]);
      })
      .finally(() => {
        setModelLoading(false);
      });
  };

  const onSubmit = (value: AddModelForm) => {
    let header = '';
    if (value.api_header_key && value.api_header_value) {
      header = value.api_header_key + '=' + value.api_header_value;
    }
    setError('');
    setCheckSuccess(false); // 开始检查时重置成功状态
    setCheckFailed(false); // 开始检查时重置失败状态
    setLoading(true);
    postCheckModel({
      api_key: value.api_key,
      model_name: value.model,
      // @ts-ignore
      owner: value.provider,
      // @ts-ignore  
      sub_type: value.type,
    })
      .then((res) => {
        // 清除之前的错误信息
        setError('');
        // 设置检查成功状态，改变按钮样式
        setCheckSuccess(true);
        setCheckFailed(false);
        setLoading(false);
      })
      .catch((error) => {
        // 设置检查失败状态，改变按钮样式为红色
        setCheckSuccess(false);
        setCheckFailed(true);
        setError('');
        
        setLoading(false);
      });
  };

  const resetCurData = (value: DomainModel) => {
    // @ts-ignore
    if (value.provider && value.provider !== 'Other') {
      getModel({
        api_key:  '',
        base_url:  '',
        model: value.id || '',
        provider: "Other",
        api_version:  '',
        api_header_key: '',
        api_header_value: '',
        type: value.model_type || ConstsModelType.ModelTypeChat,
        show_name: '',
        context_window_size: 64000,
        max_output_tokens: 8192,
        enable_r1_params: false,
        support_image: false,
        support_compute: false,
        support_prompt_caching: false
      });
    }
    reset({
        api_key:  '',
        base_url:  '',
        model: value.id || '',
        provider: "Other",
        api_version:  '',
        api_header_key: '',
        api_header_value: '',
        type: value.model_type || ConstsModelType.ModelTypeChat,
        show_name: '',
        context_window_size: 64000,
        max_output_tokens: 8192,
        enable_r1_params: false,
        support_image: false,
        support_compute: false,
        support_prompt_caching: false
    });
  };

  useEffect(() => {
    if (open) {
      if (data) {
        resetCurData(data);
      } else {
        reset({
          type: ConstsModelType.ModelTypeChat,
          provider: 'BaiZhiCloud',
          model: '',
          base_url: providers['BaiZhiCloud'].defaultBaseUrl,
          api_key: '',
          api_version: '',
          api_header_key: '',
          api_header_value: '',
          show_name: '',
          // 高级设置默认值
          context_window_size: 64000,
          max_output_tokens: 8192,
          enable_r1_params: false,
          support_image: false,
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
    <Modal
      open={open}
      width={800}
      onCancel={handleReset}
      okText={checkSuccess ? '成功！' : checkFailed ? '检查失败' : 'Check'}
      onOk={handleSubmit(onSubmit)}
      okButtonProps={{
        loading,
        disabled: false,
        style: checkSuccess 
          ? { backgroundColor: '#4caf50', borderColor: '#4caf50', color: 'white' }
          : checkFailed
          ? { backgroundColor: '#f44336', borderColor: '#f44336', color: 'white' }
          : undefined,
      }}
      cancelButtonProps={{ style: { display: 'none' } }}
    >
      <Stack direction={'row'} alignItems={'stretch'} gap={3}>
        <Stack
          gap={1}
          sx={{
            width: 200,
            flexShrink: 0,
            bgcolor: 'background.paper2',
            borderRadius: '10px',
            p: 1,
          }}
        >
          <Box
            sx={{ fontSize: 14, lineHeight: '24px', fontWeight: 'bold', p: 1 }}
          >
            模型供应商
          </Box>
          {Object.values(providers).map((it) => (
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
                  bgcolor: addOpacityToColor(theme.palette.primary.main, 0.1),
                  color: 'primary.main',
                }),
                '&:hover': {
                  color: 'primary.main',
                },
              }}
              onClick={() => {
                if (data && data.owned_by === it.label) {
                  resetCurData(data);
                } else {
                  setModelUserList([]);
                  setModelListInfo([]);
                  setError('');
                  setModelLoading(false);
                  setSuccess(false);
                  setCheckSuccess(false); // 重置检查成功状态
                  setCheckFailed(false); // 重置检查失败状态
                  reset({
                    type: ConstsModelType.ModelTypeChat,
                    provider: it.label as keyof typeof ModelProvider,
                    base_url:
                      it.label === 'AzureOpenAI' ? '' : it.defaultBaseUrl,
                    model: '',
                    api_version: '',
                    api_key: '',
                    api_header_key: '',
                    api_header_value: '',
                    show_name: '',
                    // 重置高级设置
                    context_window_size: 64000,
                    max_output_tokens: 8192,
                    enable_r1_params: false,
                    support_image: false,
                    support_compute: false,
                    support_prompt_caching: false,
                  });
                }
              }}
            >
              <IconfontIcon type={it.icon} sx={{ fontSize: 18 }} />
              {it.cn || it.label || '其他'}
            </Stack>
          ))}
        </Stack>
        <Box sx={{ flex: 1 }}>
          <Box sx={{ fontSize: 14, lineHeight: '32px' }}>
            API 地址{' '}
            <Box component={'span'} sx={{ color: 'red' }}>
              *
            </Box>
          </Box>
          <Controller
            control={control}
            name='base_url'
            rules={{
              required: {
                value: true,
                message: 'URL 不能为空',
              },
            }}
            render={({ field }: { field: any }) => (
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
                  setModelListInfo([]);
                  setValue('model', '');
                  setSuccess(false);
                  setCheckSuccess(false); // 重置检查成功状态
                  setCheckFailed(false); // 重置检查失败状态
                }}
              />
            )}
          />
          <Stack
            direction={'row'}
            alignItems={'center'}
            justifyContent={'space-between'}
            sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}
          >
            <Box>
              API Secret
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
                    '_blank'
                  )
                }
              >
                查看文档
              </Box>
            )}
          </Stack>
          <Controller
            control={control}
            name='api_key'
            render={({ field }: { field: any }) => (
              <TextField
                {...field}
                fullWidth
                size='small'
                placeholder=''
                error={!!errors.api_key}
                helperText={errors.api_key?.message}
                onChange={(e) => {
                  field.onChange(e.target.value);
                  setCheckSuccess(false); // 重置检查成功状态
                  setCheckFailed(false); // 重置检查失败状态
                  setCheckFailed(false); // 重置检查失败状态
                }}
              />
            )}
          />
          <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
            模型类型{' '}
            <Box component={'span'} sx={{ color: 'red' }}>
              *
            </Box>
          </Box>
          <Controller
            control={control}
            name='type'
            rules={{
              required: {
                value: true,
                message: '模型类型不能为空',
              },
            }}
            render={({ field }: { field: any }) => (
              <TextField
                {...field}
                fullWidth
                select
                size='small'
                error={!!errors.type}
                helperText={errors.type?.message}
                onChange={(e) => {
                  field.onChange(e.target.value);
                  setModelUserList([]);
                  setModelListInfo([]);
                  setValue('model', '');
                  setSuccess(false);
                  setCheckSuccess(false); // 重置检查成功状态
                  setCheckFailed(false); // 重置检查失败状态
                }}
              >
                <MenuItem value={ConstsModelType.ModelTypeChat}>对话模型</MenuItem>
                <MenuItem value={ConstsModelType.ModelTypeCoder}>代码补全模型</MenuItem>
                <MenuItem value={ConstsModelType.ModelTypeEmbedding}>向量模型</MenuItem>
                <MenuItem value={ConstsModelType.ModelTypeReranker}>重排序模型</MenuItem>
                <MenuItem value={ConstsModelType.ModelTypeVision}>视觉模型</MenuItem>
                <MenuItem value={ConstsModelType.ModelTypeFunctionCall}>函数调用模型</MenuItem>
              </TextField>
            )}
          />
          {providerBrand === 'AzureOpenAI' && (
            <>
              <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                API Version
              </Box>
              <Controller
                control={control}
                name='api_version'
                render={({ field }: { field: any }) => (
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
                      setValue('model', '');
                      setSuccess(false);
                      setCheckSuccess(false); // 重置检查成功状态
                      setCheckFailed(false); // 重置检查失败状态
                    }}
                  />
                )}
              />
            </>
          )}
          {providerBrand === 'Other' ? (
            <>
              <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                模型名称{' '}
                <Box component={'span'} sx={{ color: 'red' }}>
                  *
                </Box>
              </Box>
              <Controller
                control={control}
                name='model'
                rules={{
                  required: '模型名称不能为空',
                }}
                render={({ field }: { field: any }) => (
                  <TextField
                    {...field}
                    fullWidth
                    size='small'
                    placeholder=''
                    error={!!errors.model}
                    helperText={errors.model?.message}
                  />
                )}
              />
              <Box sx={{ fontSize: 12, color: 'error.main', mt: 1 }}>
                需要与模型供应商提供的名称完全一致，不要随便填写
              </Box>
            </>
          ) : modelUserList.length === 0 ? (
            <Button
              fullWidth
              variant='outlined'
              loading={modelLoading}
              sx={{ mt: 4 }}
              onClick={() => {
                handleSubmit(getModel)();
              }}
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
                name='model'
                render={({ field }: { field: any }) => (
                  <TextField
                    {...field}
                    fullWidth
                    select
                    size='small'
                    placeholder=''
                    error={!!errors.model}
                    helperText={errors.model?.message}
                    onChange={(e) => {
                      field.onChange(e.target.value);
                      setCheckSuccess(false); // 切换模型时重置检查成功状态
                      setCheckFailed(false); // 切换模型时重置检查失败状态
                    }}
                  >
                    {modelUserList.map((it) => (
                      <MenuItem key={it.model} value={it.model}>
                        {it.model}
                      </MenuItem>
                    ))}
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
                      render={({ field }: { field: any }) => (
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
                      render={({ field }: { field: any }) => (
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

        </Box>
      </Stack>
    </Modal>
  );
};

export default ModelAdd;
