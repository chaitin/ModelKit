// 主要组件
export { ModelModal } from './ModelModal';

// 类型定义
export type { ModelType, Model, ModelParam, ConstsModelStatus, ConstsModelProvider, ModelProviderConfig, ModelProviderMap, CreateModelReq as CreateModelData, ListModelReq as GetModelNameData, CheckModelReq as CheckModelData, UpdateModelReq as UpdateModelData, ModelService, ModelListItem, AddModelForm, ModelModalProps } from './types/types';

// 导出 ConstsModelType 枚举值
export { ConstsModelType } from './types/types';

// 常量
export { DEFAULT_MODEL_PROVIDERS, getProvidersByType } from './constants/providers';
export { LOCALE_MESSAGES, getLocaleMessage, getTitleMap } from './constants/locale';
export * from './constants/theme';

// 主题类型声明会通过TypeScript自动包含

// 工具函数
export {
  addOpacityToColor,
  isValidURL,
  isValidAPIKey,
  formatErrorMessage,
  debounce,
  throttle,
  deepClone,
  generateId,
  isDevelopment,
  logger,
} from './utils';

// 默认配置
export const DEFAULT_CONFIG = {
  theme: {
    primaryColor: '#1976d2',
    secondaryColor: '#dc004e',
    borderRadius: '10px',
    spacing: 2,
  },
  locale: {
    language: 'zh-CN',
    messages: {},
  },
  validation: {
    requiredFields: ['provider', 'model', 'base_url', 'api_key'],
    customValidators: {},
  },
  features: {
    enableModelTesting: true,
    enableHeaderConfig: true,
    enableApiVersion: true,
    enableProviderSelection: true,
  },
  styles: {
    modalWidth: 800,
    sidebarWidth: 200,
    customCSS: '',
    borderRadius: '10px',
  },
};

// 创建预配置的组件
export const createModelModal = (defaultConfig: Partial<typeof DEFAULT_CONFIG> = {}) => {
  const mergedConfig = { ...DEFAULT_CONFIG, ...defaultConfig };

  return (props: Omit<import('./types/types').ModelModalProps, 'config'>) => {
    return {
      ...props,
      config: mergedConfig,
    };
  };
};