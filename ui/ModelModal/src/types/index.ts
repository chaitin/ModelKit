// 基础类型定义
export type ModelType = 'chat' | 'embedding' | 'rerank';

// 模型类型常量
export enum ConstsModelType {
  ModelTypeLLM = 'chat',
  ModelTypeEmbedding = 'embedding',
  ModelTypeRerank = 'rerank',
}

// 域模型接口
export interface DomainModel {
  /** 接口地址 如：https://api.qwen.com */
  api_base?: string;
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
  model_type?: ConstsModelType;
  /** 输出token数 */
  output?: number;
  /** 高级参数 */
  param?: DomainModelParam;
  /** 提供商 */
  provider?: ConstsModelProvider;
  /** 模型显示名称 */
  show_name?: string;
  /** 状态 active:启用 inactive:禁用 */
  status?: ConstsModelStatus;
  /** 更新时间 */
  updated_at?: number;
}

export interface DomainModelParam {
  context_window?: number;
  max_tokens?: number;
  r1_enabled?: boolean;
  support_computer_use?: boolean;
  support_images?: boolean;
  support_prompt_cache?: boolean;
}

export enum ConstsModelStatus {
  ModelStatusActive = "active",
  ModelStatusInactive = "inactive",
}

export enum ConstsModelProvider {
  ModelProviderSiliconFlow = "SiliconFlow",
  ModelProviderOpenAI = "OpenAI",
  ModelProviderOllama = "Ollama",
  ModelProviderDeepSeek = "DeepSeek",
  ModelProviderMoonshot = "Moonshot",
  ModelProviderAzureOpenAI = "AzureOpenAI",
  ModelProviderBaiZhiCloud = "BaiZhiCloud",
  ModelProviderHunyuan = "Hunyuan",
  ModelProviderBaiLian = "BaiLian",
  ModelProviderVolcengine = "Volcengine",
}

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

// 创建模型数据
export interface CreateModelData {
  type: ModelType;
  provider: string;
  model: string;
  base_url: string;
  api_key: string;
  api_version?: string;
  api_header?: string;
  show_name?: string;
  param?: {
    context_window?: number;
    max_tokens?: number;
    r1_enabled?: boolean;
    support_images?: boolean;
    support_computer_use?: boolean;
    support_prompt_cache?: boolean;
  };
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
  ModelName: string;
}

// 模型服务接口
export interface ModelService {
  createModel: (data: CreateModelData) => Promise<{ model: DomainModel }>;
  listModel: (data: GetModelNameData) => Promise<{ models: ModelListItem[] }>;
  checkModel: (data: CheckModelData) => Promise<{ model: DomainModel }>;
  updateModel: (data: UpdateModelData) => Promise<{ model: DomainModel }>;
}

export interface ModelListItem {
  model?: string;
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
  show_name: string;
  // 高级设置字段
  context_window_size: number;
  max_output_tokens: number;
  enable_r1_params: boolean;
  support_image: boolean;
  support_compute: boolean;
  support_prompt_caching: boolean;
}

export interface ModelModalProps {
  open: boolean;
  data: DomainModel | null;
  type: ConstsModelType;
  onClose: () => void;
  refresh: () => void;
  modelService: ModelService;
}
