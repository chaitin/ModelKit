import { Model } from "@/types/types";
import { CLAUDE_SUPPORTED_WEBSEARCH_REGEX, DOUBAO_THINKING_MODEL_REGEX, EMBEDDING_REGEX, FUNCTION_CALLING_REGEX, GEMINI_SEARCH_REGEX, REASONING_REGEX, RERANKING_REGEX, TEXT_TO_IMAGE_REGEX, VISION_REGEX } from "./regex";
import { PERPLEXITY_SEARCH_MODELS } from "@/constants/models";

/**
 * 添加透明度到颜色
 */
export const addOpacityToColor = (color: string, opacity: number): string => {
  if (color.startsWith('#')) {
    const hex = color.slice(1);
    const r = parseInt(hex.slice(0, 2), 16);
    const g = parseInt(hex.slice(2, 4), 16);
    const b = parseInt(hex.slice(4, 6), 16);
    return `rgba(${r}, ${g}, ${b}, ${opacity})`;
  }
  return color;
};

/**
 * 验证URL格式
 */
export const isValidURL = (url: string): string => {
  try {
    const urlObj = new URL(url);

    // 1. 检查是否有 schema
    if (!urlObj.protocol) {
      return "URL 必须包含协议（如 http:// 或 https://）";
    }

    // 2. 检查是否是 localhost 或 127.0.0.1
    if (urlObj.hostname === 'localhost' || urlObj.hostname === '127.0.0.1') {
      return "请使用宿主机主机名(linux:172.17.0.1, mac/windows:host.docker.internal)";
    }

    // 3. 检查是否以斜杠结尾
    if (url.endsWith('/')) {
      return "请去掉URL末尾的斜杠";
    }

    // 4. 检查路径中是否包含 /chat/completions
    if (urlObj.pathname.includes('/chat/completions')) {
      return "请去掉URL中的/chat/completions路径";
    }

    // 5. 检查是否以 /v+数字 结尾 或者 包含/v+数字/
    const pathPattern = /\/v\d+(\/.*)?$/;
    if (!pathPattern.test(urlObj.pathname)) {
      return "模型供应商必须支持与 OpenAI 兼容的 API 格式";
    }

    return "";
  } catch {
    return "URL格式错误";
  }
};

/**
 * 验证API密钥格式
 */
export const isValidAPIKey = (key: string): boolean => {
  // 基本的API密钥验证规则
  return key.length >= 10 && /^[a-zA-Z0-9\-_]+$/.test(key);
};

/**
 * 格式化错误消息
 */
export const formatErrorMessage = (error: any): string => {
  if (typeof error === 'string') return error;
  if (error?.message) return error.message;
  if (error?.error) return error.error;
  return '未知错误';
};

/**
 * 防抖函数
 */
export const debounce = <T extends (...args: any[]) => any>(
  func: T,
  wait: number
): ((...args: Parameters<T>) => void) => {
  let timeout: NodeJS.Timeout;
  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
};

/**
 * 节流函数
 */
export const throttle = <T extends (...args: any[]) => any>(
  func: T,
  limit: number
): ((...args: Parameters<T>) => void) => {
  let inThrottle: boolean;
  return (...args: Parameters<T>) => {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
};

/**
 * 深拷贝对象
 */
export const deepClone = <T>(obj: T): T => {
  if (obj === null || typeof obj !== 'object') return obj;
  if (obj instanceof Date) return new Date(obj.getTime()) as unknown as T;
  if (obj instanceof Array) return obj.map(item => deepClone(item)) as unknown as T;
  if (typeof obj === 'object') {
    const clonedObj = {} as T;
    for (const key in obj) {
      if (obj.hasOwnProperty(key)) {
        clonedObj[key] = deepClone(obj[key]);
      }
    }
    return clonedObj;
  }
  return obj;
};

/**
 * 生成唯一ID
 */
export const generateId = (): string => {
  return Math.random().toString(36).substr(2, 9) + Date.now().toString(36);
};

/**
 * 检查是否为开发环境
 */
export const isDevelopment = (): boolean => {
  return process.env.NODE_ENV === 'development';
};

/**
 * 日志工具
 */
export const logger = {
  log: (...args: any[]) => {
    if (isDevelopment()) {
      console.log('[ModelModal]', ...args);
    }
  },
  warn: (...args: any[]) => {
    if (isDevelopment()) {
      console.warn('[ModelModal]', ...args);
    }
  },
  error: (...args: any[]) => {
    if (isDevelopment()) {
      console.error('[ModelModal]', ...args);
    }
  },
};

/**
 * 从模型 ID 中提取基础名称。
 * 例如：
 * - 'deepseek/deepseek-r1' => 'deepseek-r1'
 * - 'deepseek-ai/deepseek/deepseek-r1' => 'deepseek-r1'
 * @param {string} id 模型 ID
 * @param {string} [delimiter='/'] 分隔符，默认为 '/'
 * @returns {string} 基础名称
 */
export const getBaseModelName = (id: string, delimiter: string = '/'): string => {
  const parts = id.split(delimiter)
  return parts[parts.length - 1]
}

/**
 * 从模型 ID 中提取基础名称并转换为小写。
 * 例如：
 * - 'deepseek/DeepSeek-R1' => 'deepseek-r1'
 * - 'deepseek-ai/deepseek/DeepSeek-R1' => 'deepseek-r1'
 * @param {string} id 模型 ID
 * @param {string} [delimiter='/'] 分隔符，默认为 '/'
 * @returns {string} 小写的基础名称
 */
export const getLowerBaseModelName = (id: string, delimiter: string = '/'): string => {
  return getBaseModelName(id, delimiter).toLowerCase()
}

export function isRerankModel(model_id: string): boolean {
  const modelId = getLowerBaseModelName(model_id)
  return model_id ? RERANKING_REGEX.test(modelId) || false : false
}

export function isEmbeddingModel(model_id: string, provider: string): boolean {
  if (!model_id || isRerankModel(model_id)) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id)

  if (['anthropic'].includes(provider)) {
    return false
  }

  if (provider === 'doubao' || modelId.includes('doubao')) {
    return EMBEDDING_REGEX.test(modelId)
  }

  return EMBEDDING_REGEX.test(modelId) || false
}

export function isFunctionCallingModel(model_id: string, provider: string): boolean {
  if (!model_id || isEmbeddingModel(model_id, provider) || isRerankModel(model_id)) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id)

  if (provider === 'qiniu') {
    return ['deepseek-v3-tool', 'deepseek-v3-0324', 'qwq-32b', 'qwen2.5-72b-instruct'].includes(modelId)
  }

  if (provider === 'doubao' || modelId.includes('doubao')) {
    return FUNCTION_CALLING_REGEX.test(modelId) || FUNCTION_CALLING_REGEX.test(modelId)
  }

  if (['deepseek', 'anthropic'].includes(provider)) {
    return true
  }

  if (['kimi', 'moonshot'].includes(provider)) {
    return true
  }

  return FUNCTION_CALLING_REGEX.test(modelId)
}

export function isVisionModel(model_id: string, provider: string): boolean {
  if (!model_id || isEmbeddingModel(model_id, provider) || isRerankModel(model_id)) {
    return false
  }
  // 新添字段 copilot-vision-request 后可使用 vision
  // if (model.provider === 'copilot') {
  //   return false
  // }

  const modelId = getLowerBaseModelName(model_id)
  if (provider === 'doubao' || modelId.includes('doubao')) {
    return VISION_REGEX.test(modelId) || false
  }

  return VISION_REGEX.test(modelId) || false
}

export function isTextToImageModel(model_id: string): boolean {
  const modelId = getLowerBaseModelName(model_id)
  return TEXT_TO_IMAGE_REGEX.test(modelId)
}


export function isSupportedThinkingTokenDoubaoModel(model_id: string): boolean {
  if (!model_id) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id, '/')

  return DOUBAO_THINKING_MODEL_REGEX.test(modelId)
}

export function isClaudeReasoningModel(model_id: string): boolean {
  if (model_id) {
    return false
  }
  const modelId = getLowerBaseModelName(model_id, '/')
  return (
    modelId.includes('claude-3-7-sonnet') ||
    modelId.includes('claude-3.7-sonnet') ||
    modelId.includes('claude-sonnet-4') ||
    modelId.includes('claude-opus-4')
  )
}

/** 是否为支持思考控制的Qwen3推理模型 */
export function isSupportedThinkingTokenQwenModel(model_id:string): boolean {
  if (!model_id) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id, '/')

  if (modelId.includes('coder')) {
    return false
  }

  if (modelId.startsWith('qwen3')) {
    // instruct 是非思考模型 thinking 是思考模型，二者都不能控制思考
    if (modelId.includes('instruct') || modelId.includes('thinking')) {
      return false
    }
    return true
  }

  return [
    'qwen-plus',
    'qwen-plus-latest',
    'qwen-plus-0428',
    'qwen-plus-2025-04-28',
    'qwen-plus-0714',
    'qwen-plus-2025-07-14',
    'qwen-turbo',
    'qwen-turbo-latest',
    'qwen-turbo-0428',
    'qwen-turbo-2025-04-28',
    'qwen-turbo-0715',
    'qwen-turbo-2025-07-15'
  ].includes(modelId)
}

/** 是否为Qwen推理模型 */
export function isQwenReasoningModel(model_id: string): boolean {
  if (!model_id) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id, '/')

  if (modelId.startsWith('qwen3')) {
    if (modelId.includes('thinking')) {
      return true
    } else if (modelId.includes('instruct')) {
      return false
    }
    return true
  }

  if (isSupportedThinkingTokenQwenModel(model_id)) {
    return true
  }

  if (modelId.includes('qwq') || modelId.includes('qvq')) {
    return true
  }

  return false
}

export const isStepReasoningModel = (model_id: string): boolean => {
  if (!model_id) {
    return false
  }
  const modelId = getLowerBaseModelName(model_id, '/')
  return modelId.includes('step-3') || modelId.includes('step-r1-v-mini')
}

export const isSupportedThinkingTokenZhipuModel = (model_id: string): boolean => {
  const modelId = getLowerBaseModelName(model_id, '/')
  return modelId.includes('glm-4.5')
}

export const isZhipuReasoningModel = (model_id: string): boolean => {
  if (!model_id) {
    return false
  }
  const modelId = getLowerBaseModelName(model_id, '/')
  return isSupportedThinkingTokenZhipuModel(model_id) || modelId.includes('glm-z1')
}

export const isSupportedReasoningEffortPerplexityModel = (model_id: string): boolean => {
  const modelId = getLowerBaseModelName(model_id, '/')
  return modelId.includes('sonar-deep-research')
}

export const isPerplexityReasoningModel = (model_id: string): boolean => {
  if (!model_id) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id, '/')
  return isSupportedReasoningEffortPerplexityModel(model_id) || modelId.includes('reasoning')
}

export const isSupportedThinkingTokenHunyuanModel = (model_id: string): boolean => {
  if (!model_id) {
    return false
  }
  const modelId = getLowerBaseModelName(model_id, '/')
  return modelId.includes('hunyuan-a13b')
}

export const isHunyuanReasoningModel = (model_id: string): boolean => {
  if (!model_id) {
    return false
  }
  const modelId = getLowerBaseModelName(model_id, '/')

  return isSupportedThinkingTokenHunyuanModel(model_id) || modelId.includes('hunyuan-t1')
}

export function isSupportedReasoningEffortGrokModel(model_id: string): boolean {
  if (!model_id) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id)
  if (modelId.includes('grok-3-mini')) {
    return true
  }

  return false
}

export function isGrokReasoningModel(model_id: string): boolean {
  if (!model_id) {
    return false
  }
  const modelId = getLowerBaseModelName(model_id)
  if (isSupportedReasoningEffortGrokModel(model_id) || modelId.includes('grok-4')) {
    return true
  }

  return false
}

export const isSupportedThinkingTokenGeminiModel = (model_id: string): boolean => {
  const modelId = getLowerBaseModelName(model_id, '/')
  return modelId.includes('gemini-2.5')
}

export function isGeminiReasoningModel(model_id: string): boolean {
  if (!model_id) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id)
  if (modelId.startsWith('gemini') && modelId.includes('thinking')) {
    return true
  }

  if (isSupportedThinkingTokenGeminiModel(model_id)) {
    return true
  }

  return false
}

export const isGPT5SeriesModel = (model_id: string) => {
  const modelId = getLowerBaseModelName(model_id)
  return modelId.includes('gpt-5')
}

export function isSupportedReasoningEffortOpenAIModel(model_id: string): boolean {
  const modelId = getLowerBaseModelName(model_id)
  return (
    (modelId.includes('o1') && !(modelId.includes('o1-preview') || modelId.includes('o1-mini'))) ||
    modelId.includes('o3') ||
    modelId.includes('o4') ||
    modelId.includes('gpt-oss') ||
    (isGPT5SeriesModel(model_id) && !modelId.includes('chat'))
  )
}

export function isOpenAIReasoningModel(model_id: string): boolean {
  const modelId = getLowerBaseModelName(model_id, '/')
  return isSupportedReasoningEffortOpenAIModel(model_id) || modelId.includes('o1')
}

export function isReasoningModel(model_id: string, provider: string): boolean {
  if (!model_id || isEmbeddingModel(model_id, provider) || isRerankModel(model_id) || isTextToImageModel(model_id)) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id)

  if (provider === 'doubao' || modelId.includes('doubao')) {
    return (
      REASONING_REGEX.test(modelId) ||
      isSupportedThinkingTokenDoubaoModel(model_id) ||
      false
    )
  }

  if (
    isClaudeReasoningModel(model_id) ||
    isOpenAIReasoningModel(model_id) ||
    isGeminiReasoningModel(model_id) ||
    isQwenReasoningModel(model_id) ||
    isGrokReasoningModel(model_id) ||
    isHunyuanReasoningModel(model_id) ||
    isPerplexityReasoningModel(model_id) ||
    isZhipuReasoningModel(model_id) ||
    isStepReasoningModel(model_id) ||
    modelId.includes('magistral') ||
    modelId.includes('minimax-m1') ||
    modelId.includes('pangu-pro-moe')
  ) {
    return true
  }

  return REASONING_REGEX.test(modelId) || false
}
export const isAnthropicModel = (model_id: string): boolean => {
  if (!model_id) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id)
  return modelId.startsWith('claude')
}

export function isOpenAIWebSearchModel(model_id: string): boolean {
  const modelId = getLowerBaseModelName(model_id)

  return (
    modelId.includes('gpt-4o-search-preview') ||
    modelId.includes('gpt-4o-mini-search-preview') ||
    (modelId.includes('gpt-4.1') && !modelId.includes('gpt-4.1-nano')) ||
    (modelId.includes('gpt-4o') && !modelId.includes('gpt-4o-image')) ||
    modelId.includes('o3') ||
    modelId.includes('o4') ||
    (modelId.includes('gpt-5') && !modelId.includes('chat'))
  )
}

export function isWebSearchModel(model_id: string, provider: string): boolean {
  if (!model_id || isEmbeddingModel(model_id, provider) || isRerankModel(model_id)) {
    return false
  }

  if (!provider) {
    return false
  }

  const isEmbedding = isEmbeddingModel(model_id, provider)

  if (isEmbedding) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id, '/')

  // 不管哪个供应商都判断了
  if (isAnthropicModel(model_id)) {
    return CLAUDE_SUPPORTED_WEBSEARCH_REGEX.test(modelId)
  }


  if (provider === 'perplexity') {
    return PERPLEXITY_SEARCH_MODELS.includes(modelId)
  }

  if (provider === 'aihubmix') {
    // modelId 不以-search结尾
    if (!modelId.endsWith('-search') && GEMINI_SEARCH_REGEX.test(modelId)) {
      return true
    }

    if (isOpenAIWebSearchModel(model_id)) {
      return true
    }

    return false
  }

  if (provider === 'openai') {
    if (GEMINI_SEARCH_REGEX.test(modelId) || isOpenAIWebSearchModel(model_id)) {
      return true
    }
  }

  if (provider === 'gemini' || provider === 'vertexai') {
    return GEMINI_SEARCH_REGEX.test(modelId)
  }

  if (provider === 'hunyuan') {
    return modelId !== 'hunyuan-lite'
  }

  if (provider === 'zhipu') {
    return modelId?.startsWith('glm-4-')
  }

  if (provider === 'dashscope') {
    const models = ['qwen-turbo', 'qwen-max', 'qwen-plus', 'qwq']
    // matches id like qwen-max-0919, qwen-max-latest
    return models.some((i) => modelId.startsWith(i))
  }

  if (provider === 'openrouter') {
    return true
  }

  if (provider === 'grok') {
    return true
  }

  return false
}