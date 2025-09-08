import { FUNCTION_CALLING_EXCLUDED_MODELS, FUNCTION_CALLING_MODELS, visionAllowedModels, visionExcludedModels } from "@/constants/models"

// 检查模型是否支持视觉功能
export const VISION_REGEX = new RegExp(
  `\\b(?!(?:${visionExcludedModels.join('|')})\\b)(${visionAllowedModels.join('|')})\\b`,
  'i'
)

// 检查模型是否支持代码功能
export const CODE_REGEX = /(?:^o3$|.*(code|claude\s+sonnet|claude\s+opus|gpt-4\.1|gpt-4o|gpt-5|gemini[\s-]+2\.5|o4-mini|kimi-k2).*)/i

// 重排模型
export const RERANKING_REGEX = /(?:rerank|re-rank|re-ranker|re-ranking|retrieval|retriever)/i

// 嵌入模型
export const EMBEDDING_REGEX =
  /(?:^text-|embed|bge-|e5-|LLM2Vec|retrieval|uae-|gte-|jina-clip|jina-embeddings|voyage-)/i

// 函数调用模型
export const FUNCTION_CALLING_REGEX = new RegExp(
  `\\b(?!(?:${FUNCTION_CALLING_EXCLUDED_MODELS.join('|')})\\b)(?:${FUNCTION_CALLING_MODELS.join('|')})\\b`,
  'i'
)

// Text to image models
export const TEXT_TO_IMAGE_REGEX = /flux|diffusion|stabilityai|sd-|dall|cogview|janus|midjourney|mj-|image|gpt-image/i

// 推理模型
export const REASONING_REGEX =
  /^(o\d+(?:-[\w-]+)?|.*\b(?:reasoning|reasoner|thinking)\b.*|.*-[rR]\d+.*|.*\bqwq(?:-[\w-]+)?\b.*|.*\bhunyuan-t1(?:-[\w-]+)?\b.*|.*\bglm-zero-preview\b.*|.*\bgrok-(?:3-mini|4)(?:-[\w-]+)?\b.*)$/i

  // Doubao 支持思考模式的模型正则
export const DOUBAO_THINKING_MODEL_REGEX =
  /doubao-(?:1[.-]5-thinking-vision-pro|1[.-]5-thinking-pro-m|seed-1[.-]6(?:-flash)?(?!-(?:thinking)(?:-|$)))(?:-[\w-]+)*/i

// 支持 auto 的 Doubao 模型 doubao-seed-1.6-xxx doubao-seed-1-6-xxx  doubao-1-5-thinking-pro-m-xxx
export const DOUBAO_THINKING_AUTO_MODEL_REGEX =
  /doubao-(1-5-thinking-pro-m|seed-1[.-]6)(?!-(?:flash|thinking)(?:-|$))(?:-[\w-]+)*/i

  // 支持网络搜索的Claude模型
export const CLAUDE_SUPPORTED_WEBSEARCH_REGEX = new RegExp(
  `\\b(?:claude-3(-|\\.)(7|5)-sonnet(?:-[\\w-]+)|claude-3(-|\\.)5-haiku(?:-[\\w-]+)|claude-sonnet-4(?:-[\\w-]+)?|claude-opus-4(?:-[\\w-]+)?)\\b`,
  'i'
)

export const GEMINI_SEARCH_REGEX = new RegExp('gemini-2\\..*', 'i')
