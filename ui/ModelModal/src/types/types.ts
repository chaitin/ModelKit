// 域模型接口
export interface Model {
  /** 接口地址 如：https://api.qwen.com */
  base_url?: string;
  /** 接口头 如：Authorization: Bearer sk-xxxx */
  api_header?: string;
  /** 接口密钥 如：sk-xxxx */
  api_key?: string;
  /** 接口版本 如：2023-05-15 */
  api_version?: string;
  /** 创建时间 */
  created_at?: number;
  /** 模型ID */
  id?: string;
  /** 输入token数 */
  input?: number;
  /** 是否启用 */
  is_active?: boolean;
  /** 是否内部模型 */
  is_internal?: boolean;
  /** 模型名称 如: deepseek-v3 */
  model_name?: string;
  /** 模型类型 llm:对话模型 coder:代码模型 */
  model_type?: string;
  /** 输出token数 */
  output?: number;
  /** 高级参数 */
  param?: ModelParam;
  /** 提供商 */
  provider?: string;
  /** 模型显示名称 */
  show_name?: string;
  /** 状态 active:启用 inactive:禁用 */
  status?: string;
  /** 更新时间 */
  updated_at?: number;
}

export interface ModelParam {
  context_window?: number;
  max_tokens?: number;
  r1_enabled?: boolean;
  support_computer_use?: boolean;
  support_images?: boolean;
  support_prompt_cache?: boolean;
}

// 模型提供商配置
export interface ModelProviderConfig {
  label: string;
  cn?: string;
  icon: string;
  urlWrite: boolean;
  secretRequired: boolean;
  customHeader: boolean;
  chat: boolean;
  code: boolean;
  embedding: boolean;
  rerank: boolean;
  analysis: boolean;
  analysis_vl: boolean;
  modelDocumentUrl?: string;
  addingModelTutorial?: string;
  defaultBaseUrl: string;
}

// 模型提供商映射
export type ModelProviderMap = Record<string, ModelProviderConfig>;

// 创建模型数据
export interface CreateModelReq {
  base_url?: string;
  api_header?: string;
  api_key?: string;
  api_version?: string;
  model_name?: string;
  model_type?: string;
  param?: ModelParam;
  provider?: string;
  show_name?: string;
}

// 获取模型列表数据
export interface ListModelReq {
  api_header?: string;
  api_key?: string;
  base_url?: string;
  provider?: string;
  model_type?: string;
}

// 检查模型数据
export interface CheckModelReq {
  base_url?: string;
  api_header?: string;
  api_key?: string;
  api_version?: string;
  model_name?: string;
  provider?: string;
  model_type?: string;
  param?: ModelParam;
}

// 更新模型数据
export interface UpdateModelReq {
  /** 接口地址 如：https://api.qwen.com */
  base_url?: string;
  api_header?: string;
  /** 接口密钥 如：sk-xxxx */
  api_key?: string;
  api_version?: string;
  /** 模型ID */
  id?: string;
  /** 模型名称 */
  model_name?: string;
  /** 高级参数 */
  param?: ModelParam;
  /** 提供商 */
  provider?: string;
  /** 模型显示名称 */
  show_name?: string;
  /** 状态 active:启用 inactive:禁用 */
  status?: string;
  model_type?: string;
}

// 模型服务接口
export interface ModelService {
  createModel: (
    data: CreateModelReq,
  ) => Promise<{ model: Model; error?: string }>;
  listModel: (
    data: ListModelReq,
  ) => Promise<{ models: ModelListItem[]; error?: string }>;
  checkModel: (
    data: CheckModelReq,
  ) => Promise<{ model: Model; error?: string }>;
  updateModel: (
    data: UpdateModelReq,
  ) => Promise<{ model: Model; error?: string }>;
}

export interface ModelListItem {
  model?: string;
}

// 表单数据
export interface AddModelForm {
  provider: string;
  model_name: string;
  base_url: string;
  api_version: string;
  api_key: string;
  api_header_key: string;
  api_header_value: string;
  model_type: string;
  show_name: string;
  api_header: string;
  resource_name: string; // Azure OpenAI resource name
  // 高级设置字段
  context_window_size: number;
  max_output_tokens: number;
  enable_r1_params: boolean;
  support_image: boolean;
  support_compute: boolean;
  support_prompt_caching: boolean;
}

// 消息组件接口
export interface MessageComponent {
  error: (content: string) => void;
  success: (content: string) => void;
  info?: (content: string) => void;
  warning?: (content: string) => void;
}

export interface ModelModalProps {
  open: boolean;
  data: Model | null;
  model_type: string;
  onClose: () => void;
  refresh: () => void;
  modelService: ModelService;
  language?: 'zh-CN' | 'en-US';
  messageComponent?: MessageComponent;
  is_close_model_remark?: boolean;
  addingModelTutorialURL?: string;
  beforeSubmit?: (data: AddModelForm) => boolean | Promise<boolean>;
  onOk?: (data: AddModelForm) => void;
  loading?: boolean;
}
