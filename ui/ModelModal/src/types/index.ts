// 基础类型定义
export type ModelType = 'chat' | 'embedding' | 'rerank';

// 模型提供商配置
export interface ModelProviderConfig {
  label: string;
  cn?: string;
  icon: string;
  urlWrite: boolean;
  secretRequired: boolean;
  customHeader: boolean;
  modelDocumentUrl?: string;
  defaultBaseUrl: string;
}

// 模型提供商映射
export type ModelProviderMap = Record<string, ModelProviderConfig>;

// 模型列表项
export interface ModelListItem {
  id: string;
  model: string;
  type: ModelType;
  provider: string;
  base_url: string;
  api_key: string;
  api_version?: string;
  api_header?: string;
  completion_tokens?: number;
  prompt_tokens?: number;
  total_tokens?: number;
}

// 创建模型数据
export interface CreateModelData {
  type: ModelType;
  provider: string;
  model: string;
  base_url: string;
  api_key: string;
  api_version?: string;
  api_header?: string;
}

// 获取模型列表数据
export interface GetModelNameData {
  type: ModelType;
  provider: string;
  base_url: string;
  api_key: string;
  api_header?: string;
}

// 检查模型数据
export interface CheckModelData extends CreateModelData {
  api_version: string;
}

// 更新模型数据
export interface UpdateModelData extends CheckModelData {
  id: string;
}

// 模型服务接口
export interface ModelService {
  createModel: (data: CreateModelData) => Promise<{ id: string }>;
  getModelNameList: (data: GetModelNameData) => Promise<{ models: { model: string }[] }>;
  testModel: (data: CheckModelData) => Promise<{ error: string }>;
  updateModel: (data: UpdateModelData) => Promise<void>;
}

// 表单数据
export interface AddModelForm {
  provider: string;
  model: string;
  base_url: string;
  api_version: string;
  api_key: string;
  api_header_key: string;
  api_header_value: string;
  type: ModelType;
}

// 组件配置接口
export interface ModelModalConfig {
  // 主题配置
  theme?: {
    primaryColor?: string;
    secondaryColor?: string;
    borderRadius?: string;
    spacing?: number;
  };
  
  // 本地化配置
  locale?: {
    language?: 'zh-CN' | 'en-US';
    messages?: Record<string, string>;
  };
  
  // 验证配置
  validation?: {
    requiredFields?: (keyof AddModelForm)[];
    customValidators?: Record<string, (value: any) => string | undefined>;
  };
  
  // 功能开关
  features?: {
    enableModelTesting?: boolean;
    enableHeaderConfig?: boolean;
    enableApiVersion?: boolean;
    enableProviderSelection?: boolean;
  };
  
  // 自定义样式
  styles?: {
    modalWidth?: number | string;
    sidebarWidth?: number | string;
    customCSS?: string;
    borderRadius?: string;
  };
}

// 组件属性
export interface ModelModalProps {
  open: boolean;
  data?: ModelListItem | null;
  type: ModelType;
  onClose: () => void;
  refresh: () => void;
  modelService: ModelService;
  config?: ModelModalConfig;
  
  // 自定义回调
  onBeforeSubmit?: (data: AddModelForm) => boolean | Promise<boolean>;
  onAfterSubmit?: (data: AddModelForm, result: any) => void;
  onError?: (error: string) => void;
  
  // 自定义渲染
  customProviders?: ModelProviderMap;
  customValidation?: Record<string, (value: any) => string | undefined>;
  customStyles?: Record<string, any>;
} 