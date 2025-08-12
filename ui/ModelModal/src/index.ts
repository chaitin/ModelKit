// 主要组件
export { ModelModal } from './ModelModal';

// 类型定义
export type {
  ModelModalProps,
  ModelModalConfig,
  ModelService,
  ModelListItem,
  CreateModelData,
  GetModelNameData,
  CheckModelData,
  UpdateModelData,
  AddModelForm,
  ModelType,
  ModelProviderConfig,
  ModelProviderMap,
} from './types';

// 默认服务实现
export { DefaultModelService, createDefaultModelService } from './services/DefaultModelService';

// 常量
export { DEFAULT_MODEL_PROVIDERS, getProvidersByType } from './constants/providers';
export { LOCALE_MESSAGES, getLocaleMessage, getTitleMap } from './constants/locale';

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
    language: 'zh-CN' as const,
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
} as const;

// 创建预配置的组件
export const createModelModal = (defaultConfig: Partial<typeof DEFAULT_CONFIG> = {}) => {
  const mergedConfig = { ...DEFAULT_CONFIG, ...defaultConfig };
  
  return (props: Omit<import('./types').ModelModalProps, 'config'>) => {
    return {
      ...props,
      config: mergedConfig,
    };
  };
}; 