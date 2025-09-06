import { getLowerBaseModelName } from "."
import { CLAUDE_SUPPORTED_WEBSEARCH_REGEX, DOUBAO_THINKING_MODEL_REGEX, EMBEDDING_REGEX, FUNCTION_CALLING_REGEX, GEMINI_SEARCH_REGEX, REASONING_REGEX, RERANKING_REGEX, TEXT_TO_IMAGE_REGEX, VISION_REGEX } from "./regex";
import { PERPLEXITY_SEARCH_MODELS } from "@/constants/models";

import Ai360ModelLogo from '@/assets/images/models/360.png'
import Ai360ModelLogoDark from '@/assets/images/models/360_dark.png'
import AdeptModelLogo from '@/assets/images/models/adept.png'
import AdeptModelLogoDark from '@/assets/images/models/adept_dark.png'
import Ai21ModelLogo from '@/assets/images/models/ai21.png'
import Ai21ModelLogoDark from '@/assets/images/models/ai21_dark.png'
import AimassModelLogo from '@/assets/images/models/aimass.png'
import AimassModelLogoDark from '@/assets/images/models/aimass_dark.png'
import AisingaporeModelLogo from '@/assets/images/models/aisingapore.png'
import AisingaporeModelLogoDark from '@/assets/images/models/aisingapore_dark.png'
import BaichuanModelLogo from '@/assets/images/models/baichuan.png'
import BaichuanModelLogoDark from '@/assets/images/models/baichuan_dark.png'
import BgeModelLogo from '@/assets/images/models/bge.webp'
import BigcodeModelLogo from '@/assets/images/models/bigcode.webp'
import BigcodeModelLogoDark from '@/assets/images/models/bigcode_dark.webp'
import ChatGLMModelLogo from '@/assets/images/models/chatglm.png'
import ChatGLMModelLogoDark from '@/assets/images/models/chatglm_dark.png'
import ChatGptModelLogo from '@/assets/images/models/chatgpt.jpeg'
import ClaudeModelLogo from '@/assets/images/models/claude.png'
import ClaudeModelLogoDark from '@/assets/images/models/claude_dark.png'
import CodegeexModelLogo from '@/assets/images/models/codegeex.png'
import CodegeexModelLogoDark from '@/assets/images/models/codegeex_dark.png'
import CodestralModelLogo from '@/assets/images/models/codestral.png'
import CohereModelLogo from '@/assets/images/models/cohere.png'
import CohereModelLogoDark from '@/assets/images/models/cohere_dark.png'
import CopilotModelLogo from '@/assets/images/models/copilot.png'
import CopilotModelLogoDark from '@/assets/images/models/copilot_dark.png'
import DalleModelLogo from '@/assets/images/models/dalle.png'
import DalleModelLogoDark from '@/assets/images/models/dalle_dark.png'
import DbrxModelLogo from '@/assets/images/models/dbrx.png'
import DeepSeekModelLogo from '@/assets/images/models/deepseek.png'
import DeepSeekModelLogoDark from '@/assets/images/models/deepseek_dark.png'
import DianxinModelLogo from '@/assets/images/models/dianxin.png'
import DianxinModelLogoDark from '@/assets/images/models/dianxin_dark.png'
import DoubaoModelLogo from '@/assets/images/models/doubao.png'
import DoubaoModelLogoDark from '@/assets/images/models/doubao_dark.png'
import {
  default as EmbeddingModelLogo,
  default as EmbeddingModelLogoDark
} from '@/assets/images/models/embedding.png'
import FlashaudioModelLogo from '@/assets/images/models/flashaudio.png'
import FlashaudioModelLogoDark from '@/assets/images/models/flashaudio_dark.png'
import FluxModelLogo from '@/assets/images/models/flux.png'
import FluxModelLogoDark from '@/assets/images/models/flux_dark.png'
import GeminiModelLogo from '@/assets/images/models/gemini.png'
import GeminiModelLogoDark from '@/assets/images/models/gemini_dark.png'
import GemmaModelLogo from '@/assets/images/models/gemma.png'
import GemmaModelLogoDark from '@/assets/images/models/gemma_dark.png'
import { default as GoogleModelLogo, default as GoogleModelLogoDark } from '@/assets/images/models/google.png'
import ChatGPT35ModelLogo from '@/assets/images/models/gpt_3.5.png'
import ChatGPT4ModelLogo from '@/assets/images/models/gpt_4.png'
import {
  default as ChatGPT4ModelLogoDark,
  default as ChatGPT35ModelLogoDark,
  default as ChatGptModelLogoDark,
  default as ChatGPTo1ModelLogoDark
} from '@/assets/images/models/gpt_dark.png'
import ChatGPTImageModelLogo from '@/assets/images/models/gpt_image_1.png'
import ChatGPTo1ModelLogo from '@/assets/images/models/gpt_o1.png'
import GPT5ModelLogo from '@/assets/images/models/gpt-5.png'
import GPT5ChatModelLogo from '@/assets/images/models/gpt-5-chat.png'
import GPT5MiniModelLogo from '@/assets/images/models/gpt-5-mini.png'
import GPT5NanoModelLogo from '@/assets/images/models/gpt-5-nano.png'
import GrokModelLogo from '@/assets/images/models/grok.png'
import GrokModelLogoDark from '@/assets/images/models/grok_dark.png'
import GrypheModelLogo from '@/assets/images/models/gryphe.png'
import GrypheModelLogoDark from '@/assets/images/models/gryphe_dark.png'
import HailuoModelLogo from '@/assets/images/models/hailuo.png'
import HailuoModelLogoDark from '@/assets/images/models/hailuo_dark.png'
import HuggingfaceModelLogo from '@/assets/images/models/huggingface.png'
import HuggingfaceModelLogoDark from '@/assets/images/models/huggingface_dark.png'
import HunyuanModelLogo from '@/assets/images/models/hunyuan.png'
import HunyuanModelLogoDark from '@/assets/images/models/hunyuan_dark.png'
import IbmModelLogo from '@/assets/images/models/ibm.png'
import IbmModelLogoDark from '@/assets/images/models/ibm_dark.png'
import InternlmModelLogo from '@/assets/images/models/internlm.png'
import InternlmModelLogoDark from '@/assets/images/models/internlm_dark.png'
import InternvlModelLogo from '@/assets/images/models/internvl.png'
import JinaModelLogo from '@/assets/images/models/jina.png'
import JinaModelLogoDark from '@/assets/images/models/jina_dark.png'
import KeLingModelLogo from '@/assets/images/models/keling.png'
import KeLingModelLogoDark from '@/assets/images/models/keling_dark.png'
import LlamaModelLogo from '@/assets/images/models/llama.png'
import LlamaModelLogoDark from '@/assets/images/models/llama_dark.png'
import LLavaModelLogo from '@/assets/images/models/llava.png'
import LLavaModelLogoDark from '@/assets/images/models/llava_dark.png'
import LumaModelLogo from '@/assets/images/models/luma.png'
import LumaModelLogoDark from '@/assets/images/models/luma_dark.png'
import MagicModelLogo from '@/assets/images/models/magic.png'
import MagicModelLogoDark from '@/assets/images/models/magic_dark.png'
import MediatekModelLogo from '@/assets/images/models/mediatek.png'
import MediatekModelLogoDark from '@/assets/images/models/mediatek_dark.png'
import MicrosoftModelLogo from '@/assets/images/models/microsoft.png'
import MicrosoftModelLogoDark from '@/assets/images/models/microsoft_dark.png'
import MidjourneyModelLogo from '@/assets/images/models/midjourney.png'
import MidjourneyModelLogoDark from '@/assets/images/models/midjourney_dark.png'
import {
  default as MinicpmModelLogo,
  default as MinicpmModelLogoDark
} from '@/assets/images/models/minicpm.webp'
import MinimaxModelLogo from '@/assets/images/models/minimax.png'
import MinimaxModelLogoDark from '@/assets/images/models/minimax_dark.png'
import MistralModelLogo from '@/assets/images/models/mixtral.png'
import MistralModelLogoDark from '@/assets/images/models/mixtral_dark.png'
import MoonshotModelLogo from '@/assets/images/models/moonshot.png'
import MoonshotModelLogoDark from '@/assets/images/models/moonshot_dark.png'
import {
  default as NousResearchModelLogo,
  default as NousResearchModelLogoDark
} from '@/assets/images/models/nousresearch.png'
import NvidiaModelLogo from '@/assets/images/models/nvidia.png'
import NvidiaModelLogoDark from '@/assets/images/models/nvidia_dark.png'
import PalmModelLogo from '@/assets/images/models/palm.png'
import PalmModelLogoDark from '@/assets/images/models/palm_dark.png'
import PanguModelLogo from '@/assets/images/models/pangu.svg'
import {
  default as PerplexityModelLogo,
  default as PerplexityModelLogoDark
} from '@/assets/images/models/perplexity.png'
import PixtralModelLogo from '@/assets/images/models/pixtral.png'
import PixtralModelLogoDark from '@/assets/images/models/pixtral_dark.png'
import QwenModelLogo from '@/assets/images/models/qwen.png'
import QwenModelLogoDark from '@/assets/images/models/qwen_dark.png'
import RakutenaiModelLogo from '@/assets/images/models/rakutenai.png'
import RakutenaiModelLogoDark from '@/assets/images/models/rakutenai_dark.png'
import SparkDeskModelLogo from '@/assets/images/models/sparkdesk.png'
import SparkDeskModelLogoDark from '@/assets/images/models/sparkdesk_dark.png'
import StabilityModelLogo from '@/assets/images/models/stability.png'
import StabilityModelLogoDark from '@/assets/images/models/stability_dark.png'
import StepModelLogo from '@/assets/images/models/step.png'
import StepModelLogoDark from '@/assets/images/models/step_dark.png'
import SunoModelLogo from '@/assets/images/models/suno.png'
import SunoModelLogoDark from '@/assets/images/models/suno_dark.png'
import TeleModelLogo from '@/assets/images/models/tele.png'
import TeleModelLogoDark from '@/assets/images/models/tele_dark.png'
import TokenFluxModelLogo from '@/assets/images/models/tokenflux.png'
import TokenFluxModelLogoDark from '@/assets/images/models/tokenflux_dark.png'
import UpstageModelLogo from '@/assets/images/models/upstage.png'
import UpstageModelLogoDark from '@/assets/images/models/upstage_dark.png'
import ViduModelLogo from '@/assets/images/models/vidu.png'
import ViduModelLogoDark from '@/assets/images/models/vidu_dark.png'
import VoyageModelLogo from '@/assets/images/models/voyageai.png'
import WenxinModelLogo from '@/assets/images/models/wenxin.png'
import WenxinModelLogoDark from '@/assets/images/models/wenxin_dark.png'
import XirangModelLogo from '@/assets/images/models/xirang.png'
import XirangModelLogoDark from '@/assets/images/models/xirang_dark.png'
import YiModelLogo from '@/assets/images/models/yi.png'
import YiModelLogoDark from '@/assets/images/models/yi_dark.png'
import YoudaoLogo from '@/assets/images/providers/netease-youdao.svg'
import NomicLogo from '@/assets/images/providers/nomic.png'

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
  if (!provider || !model_id || isEmbeddingModel(model_id, provider) || isRerankModel(model_id)) {
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

export const getModelGroup = (model_id: string): string => {
  // 1. 提取model_id第一个-之前的部分
  const firstPart = model_id.split('-')[0];
  
  // 2. 从头获取连续的纯字母部分
  const withoutNumbers = firstPart.match(/^[a-zA-Z]+/)?.[0] || '';
  
  // 3. 返回结果
  return withoutNumbers;
}

export function getModelLogo(modelId: string) {
  const isLight = true

  if (!modelId) {
    return undefined
  }

  const logoMap = {
    pixtral: isLight ? PixtralModelLogo : PixtralModelLogoDark,
    jina: isLight ? JinaModelLogo : JinaModelLogoDark,
    abab: isLight ? MinimaxModelLogo : MinimaxModelLogoDark,
    minimax: isLight ? MinimaxModelLogo : MinimaxModelLogoDark,
    o1: isLight ? ChatGPTo1ModelLogo : ChatGPTo1ModelLogoDark,
    o3: isLight ? ChatGPTo1ModelLogo : ChatGPTo1ModelLogoDark,
    o4: isLight ? ChatGPTo1ModelLogo : ChatGPTo1ModelLogoDark,
    'gpt-image': ChatGPTImageModelLogo,
    'gpt-3': isLight ? ChatGPT35ModelLogo : ChatGPT35ModelLogoDark,
    'gpt-4': isLight ? ChatGPT4ModelLogo : ChatGPT4ModelLogoDark,
    'gpt-5-mini': GPT5MiniModelLogo,
    'gpt-5-nano': GPT5NanoModelLogo,
    'gpt-5-chat': GPT5ChatModelLogo,
    'gpt-5': GPT5ModelLogo,
    gpts: isLight ? ChatGPT4ModelLogo : ChatGPT4ModelLogoDark,
    'gpt-oss(?:-[\\w-]+)': isLight ? ChatGptModelLogo : ChatGptModelLogoDark,
    'text-moderation': isLight ? ChatGptModelLogo : ChatGptModelLogoDark,
    'babbage-': isLight ? ChatGptModelLogo : ChatGptModelLogoDark,
    'sora-': isLight ? ChatGptModelLogo : ChatGptModelLogoDark,
    '(^|/)omni-': isLight ? ChatGptModelLogo : ChatGptModelLogoDark,
    'Embedding-V1': isLight ? WenxinModelLogo : WenxinModelLogoDark,
    'text-embedding-v': isLight ? QwenModelLogo : QwenModelLogoDark,
    'text-embedding': isLight ? ChatGptModelLogo : ChatGptModelLogoDark,
    'davinci-': isLight ? ChatGptModelLogo : ChatGptModelLogoDark,
    glm: isLight ? ChatGLMModelLogo : ChatGLMModelLogoDark,
    deepseek: isLight ? DeepSeekModelLogo : DeepSeekModelLogoDark,
    '(qwen|qwq|qwq-|qvq-)': isLight ? QwenModelLogo : QwenModelLogoDark,
    gemma: isLight ? GemmaModelLogo : GemmaModelLogoDark,
    'yi-': isLight ? YiModelLogo : YiModelLogoDark,
    llama: isLight ? LlamaModelLogo : LlamaModelLogoDark,
    mixtral: isLight ? MistralModelLogo : MistralModelLogo,
    mistral: isLight ? MistralModelLogo : MistralModelLogoDark,
    codestral: CodestralModelLogo,
    ministral: isLight ? MistralModelLogo : MistralModelLogoDark,
    magistral: isLight ? MistralModelLogo : MistralModelLogoDark,
    moonshot: isLight ? MoonshotModelLogo : MoonshotModelLogoDark,
    kimi: isLight ? MoonshotModelLogo : MoonshotModelLogoDark,
    phi: isLight ? MicrosoftModelLogo : MicrosoftModelLogoDark,
    baichuan: isLight ? BaichuanModelLogo : BaichuanModelLogoDark,
    claude: isLight ? ClaudeModelLogo : ClaudeModelLogoDark,
    gemini: isLight ? GeminiModelLogo : GeminiModelLogoDark,
    bison: isLight ? PalmModelLogo : PalmModelLogoDark,
    palm: isLight ? PalmModelLogo : PalmModelLogoDark,
    step: isLight ? StepModelLogo : StepModelLogoDark,
    hailuo: isLight ? HailuoModelLogo : HailuoModelLogoDark,
    doubao: isLight ? DoubaoModelLogo : DoubaoModelLogoDark,
    'ep-202': isLight ? DoubaoModelLogo : DoubaoModelLogoDark,
    cohere: isLight ? CohereModelLogo : CohereModelLogoDark,
    command: isLight ? CohereModelLogo : CohereModelLogoDark,
    minicpm: isLight ? MinicpmModelLogo : MinicpmModelLogoDark,
    '360': isLight ? Ai360ModelLogo : Ai360ModelLogoDark,
    aimass: isLight ? AimassModelLogo : AimassModelLogoDark,
    codegeex: isLight ? CodegeexModelLogo : CodegeexModelLogoDark,
    copilot: isLight ? CopilotModelLogo : CopilotModelLogoDark,
    creative: isLight ? CopilotModelLogo : CopilotModelLogoDark,
    balanced: isLight ? CopilotModelLogo : CopilotModelLogoDark,
    precise: isLight ? CopilotModelLogo : CopilotModelLogoDark,
    dalle: isLight ? DalleModelLogo : DalleModelLogoDark,
    'dall-e': isLight ? DalleModelLogo : DalleModelLogoDark,
    dbrx: isLight ? DbrxModelLogo : DbrxModelLogo,
    flashaudio: isLight ? FlashaudioModelLogo : FlashaudioModelLogoDark,
    flux: isLight ? FluxModelLogo : FluxModelLogoDark,
    grok: isLight ? GrokModelLogo : GrokModelLogoDark,
    hunyuan: isLight ? HunyuanModelLogo : HunyuanModelLogoDark,
    internlm: isLight ? InternlmModelLogo : InternlmModelLogoDark,
    internvl: InternvlModelLogo,
    llava: isLight ? LLavaModelLogo : LLavaModelLogoDark,
    magic: isLight ? MagicModelLogo : MagicModelLogoDark,
    midjourney: isLight ? MidjourneyModelLogo : MidjourneyModelLogoDark,
    'mj-': isLight ? MidjourneyModelLogo : MidjourneyModelLogoDark,
    'tao-': isLight ? WenxinModelLogo : WenxinModelLogoDark,
    'ernie-': isLight ? WenxinModelLogo : WenxinModelLogoDark,
    voice: isLight ? FlashaudioModelLogo : FlashaudioModelLogoDark,
    'tts-1': isLight ? ChatGptModelLogo : ChatGptModelLogoDark,
    'whisper-': isLight ? ChatGptModelLogo : ChatGptModelLogoDark,
    'stable-': isLight ? StabilityModelLogo : StabilityModelLogoDark,
    sd2: isLight ? StabilityModelLogo : StabilityModelLogoDark,
    sd3: isLight ? StabilityModelLogo : StabilityModelLogoDark,
    sdxl: isLight ? StabilityModelLogo : StabilityModelLogoDark,
    sparkdesk: isLight ? SparkDeskModelLogo : SparkDeskModelLogoDark,
    generalv: isLight ? SparkDeskModelLogo : SparkDeskModelLogoDark,
    wizardlm: isLight ? MicrosoftModelLogo : MicrosoftModelLogoDark,
    microsoft: isLight ? MicrosoftModelLogo : MicrosoftModelLogoDark,
    hermes: isLight ? NousResearchModelLogo : NousResearchModelLogoDark,
    gryphe: isLight ? GrypheModelLogo : GrypheModelLogoDark,
    suno: isLight ? SunoModelLogo : SunoModelLogoDark,
    chirp: isLight ? SunoModelLogo : SunoModelLogoDark,
    luma: isLight ? LumaModelLogo : LumaModelLogoDark,
    keling: isLight ? KeLingModelLogo : KeLingModelLogoDark,
    'vidu-': isLight ? ViduModelLogo : ViduModelLogoDark,
    ai21: isLight ? Ai21ModelLogo : Ai21ModelLogoDark,
    'jamba-': isLight ? Ai21ModelLogo : Ai21ModelLogoDark,
    mythomax: isLight ? GrypheModelLogo : GrypheModelLogoDark,
    nvidia: isLight ? NvidiaModelLogo : NvidiaModelLogoDark,
    dianxin: isLight ? DianxinModelLogo : DianxinModelLogoDark,
    tele: isLight ? TeleModelLogo : TeleModelLogoDark,
    adept: isLight ? AdeptModelLogo : AdeptModelLogoDark,
    aisingapore: isLight ? AisingaporeModelLogo : AisingaporeModelLogoDark,
    bigcode: isLight ? BigcodeModelLogo : BigcodeModelLogoDark,
    mediatek: isLight ? MediatekModelLogo : MediatekModelLogoDark,
    upstage: isLight ? UpstageModelLogo : UpstageModelLogoDark,
    rakutenai: isLight ? RakutenaiModelLogo : RakutenaiModelLogoDark,
    ibm: isLight ? IbmModelLogo : IbmModelLogoDark,
    'google/': isLight ? GoogleModelLogo : GoogleModelLogoDark,
    xirang: isLight ? XirangModelLogo : XirangModelLogoDark,
    hugging: isLight ? HuggingfaceModelLogo : HuggingfaceModelLogoDark,
    youdao: YoudaoLogo,
    embedding: isLight ? EmbeddingModelLogo : EmbeddingModelLogoDark,
    perplexity: isLight ? PerplexityModelLogo : PerplexityModelLogoDark,
    sonar: isLight ? PerplexityModelLogo : PerplexityModelLogoDark,
    'bge-': BgeModelLogo,
    'voyage-': VoyageModelLogo,
    tokenflux: isLight ? TokenFluxModelLogo : TokenFluxModelLogoDark,
    'nomic-': NomicLogo,
    'pangu-': PanguModelLogo
  }

  for (const key in logoMap) {
    const regex = new RegExp(key, 'i')
    if (regex.test(modelId)) {
      return logoMap[key]
    }
  }

  return undefined
}