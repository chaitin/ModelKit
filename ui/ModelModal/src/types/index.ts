import { ModelProvider } from "@/constants/providers";

// 基础类型定义
export type ModelType = 'chat' | 'embedding' | 'rerank' | 'coder' | 'audio';

// 模型类型常量
export enum ConstsModelType {
  ModelTypeChat = "chat",
  ModelTypeCoder = "coder",
  ModelTypeEmbedding = "embedding",
  ModelTypeAudio = "audio",
  ModelTypeRerank = "rerank",
}

// 域模型接口
export interface Model {
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
  param?: ModelParam;
  /** 提供商 */
  provider?: ConstsModelProvider;
  /** 模型显示名称 */
  show_name?: string;
  /** 状态 active:启用 inactive:禁用 */
  status?: ConstsModelStatus;
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
export interface CreateModelReq {
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
export interface ListModelReq {
  api_header?: string;
  api_key?: string;
  base_url: string;
  provider:
  | "SiliconFlow"
  | "OpenAI"
  | "Ollama"
  | "DeepSeek"
  | "Moonshot"
  | "AzureOpenAI"
  | "BaiZhiCloud"
  | "Hunyuan"
  | "BaiLian"
  | "Volcengine";
  type: "chat" | "coder" | "embedding" | "audio" | "rerank";
}

// 检查模型数据
export interface CheckModelReq {
  /** 接口地址 */
  api_base: string;
  api_header?: string;
  /** 接口密钥 */
  api_key: string;
  api_version?: string;
  /** 模型名称 */
  model_name: string;
  /** 提供商 */
  provider: ConstsModelProvider;
  type: "llm" | "coder" | "embedding" | "rerank";
}

// 更新模型数据
export interface UpdateModelReq {
  /** 接口地址 如：https://api.qwen.com */
  api_base?: string;
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
  provider:
    | "SiliconFlow"
    | "OpenAI"
    | "Ollama"
    | "DeepSeek"
    | "Moonshot"
    | "AzureOpenAI"
    | "BaiZhiCloud"
    | "Hunyuan"
    | "BaiLian"
    | "Volcengine"
    | "Other";
  /** 模型显示名称 */
  show_name?: string;
  /** 状态 active:启用 inactive:禁用 */
  status?: ConstsModelStatus;
}

// 模型服务接口
export interface ModelService {
  createModel: (data: CreateModelReq) => Promise<{ model: Model }>;
  listModel: (data: ListModelReq) => Promise<{ models: ModelListItem[] }>;
  checkModel: (data: CheckModelReq) => Promise<{ model: Model }>;
  updateModel: (data: UpdateModelReq) => Promise<{ model: Model }>;
}

export interface ModelListItem {
  model?: string;
}

// 表单数据
export interface AddModelForm {
  provider: keyof typeof ModelProvider;
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
  data: Model | null;
  type: ConstsModelType;
  onClose: () => void;
  refresh: () => void;
  modelService: ModelService;
}
