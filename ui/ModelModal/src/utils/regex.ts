import {  NOT_TOOL_CALL_MODELS, TOOL_CALL_MODELS,  visionModels,  notVisionModels } from "@/constants/models"

export const VISION_CHECK = new RegExp(
  `\\b(?!(?:${ notVisionModels.join('|')})\\b)(${ visionModels.join('|')})\\b`,
  'i'
)

export const CODE_CHECK = /(?:^o3$|.*(code|claude\s+sonnet|claude\s+opus|gpt-4\.1|gpt-4o|gpt-5|gemini[\s-]+2\.5|o4-mini|kimi-k2).*)/i

export const RERANKING_CHECK = /(?:rerank|re-rank|re-ranker|re-ranking|retrieval|retriever)/i

export const EMBEDDING_CHECK =
  /(?:^text-|embed|bge-|e5-|LLM2Vec|retrieval|uae-|gte-|jina-clip|jina-embeddings|voyage-)/i

export const FUNCTION_CALLING_CHECK = new RegExp(
  `\\b(?!(?:${ NOT_TOOL_CALL_MODELS.join('|')})\\b)(?:${TOOL_CALL_MODELS.join('|')})\\b`,
  'i'
)

export const TEXT_TO_IMAGE_CHECK = /flux|diffusion|stabilityai|sd-|dall|cogview|janus|midjourney|mj-|image|gpt-image/i

export const REASONING_CHECK =
  /^(o\d+(?:-[\w-]+)?|.*\b(?:reasoning|reasoner|thinking)\b.*|.*-[rR]\d+.*|.*\bqwq(?:-[\w-]+)?\b.*|.*\bhunyuan-t1(?:-[\w-]+)?\b.*|.*\bglm-zero-preview\b.*|.*\bgrok-(?:3-mini|4)(?:-[\w-]+)?\b.*)$/i

export const DOUBAO_THINKING_MODEL_CHECK =
  /doubao-(?:1[.-]5-thinking-vision-pro|1[.-]5-thinking-pro-m|seed-1[.-]6(?:-flash)?(?!-(?:thinking)(?:-|$)))(?:-[\w-]+)*/i

export const CLAUDE_SUPPORTED_WEBSEARCH_CHECK = new RegExp(
  `\\b(?:claude-3(-|\\.)(7|5)-sonnet(?:-[\\w-]+)|claude-3(-|\\.)5-haiku(?:-[\\w-]+)|claude-sonnet-4(?:-[\\w-]+)?|claude-opus-4(?:-[\\w-]+)?)\\b`,
  'i'
)

export const GEMINI_SEARCH_CHECK = new RegExp('gemini-2\\..*', 'i')
