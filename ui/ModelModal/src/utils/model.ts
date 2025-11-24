import { getLowerBaseModelName } from "."
import { NOT_TOOL_CALL_MODELS, notVisionModels, PERPLEXITY_MODELS, TOOL_CALL_MODELS, visionModels } from "@/constants/models";

// Import all model logos using barrel exports
import {
  // Light theme logos
  Ai360ModelLogo, AdeptModelLogo, Ai21ModelLogo, AimassModelLogo, AisingaporeModelLogo,
  BaichuanModelLogo, BgeModelLogo, BigcodeModelLogo, ChatGLMModelLogo, ChatGptModelLogo,
  ClaudeModelLogo, CodegeexModelLogo, CodestralModelLogo, CohereModelLogo, CopilotModelLogo,
  DalleModelLogo, DbrxModelLogo, DeepSeekModelLogo, DianxinModelLogo, DoubaoModelLogo,
  EmbeddingModelLogo, FlashaudioModelLogo, FluxModelLogo, GeminiModelLogo, GemmaModelLogo,
  GoogleModelLogo, ChatGPT35ModelLogo, ChatGPT4ModelLogo, ChatGPTImageModelLogo, ChatGPTo1ModelLogo,
  GPT5ModelLogo, GPT5ChatModelLogo, GPT5MiniModelLogo, GPT5NanoModelLogo, GrokModelLogo,
  GrypheModelLogo, HailuoModelLogo, HuggingfaceModelLogo, HunyuanModelLogo, IbmModelLogo,
  InternlmModelLogo, InternvlModelLogo, JinaModelLogo, KeLingModelLogo, LlamaModelLogo,
  LLavaModelLogo, LumaModelLogo, MagicModelLogo, MediatekModelLogo, MicrosoftModelLogo,
  MidjourneyModelLogo, MinicpmModelLogo, MinimaxModelLogo, MistralModelLogo, MoonshotModelLogo,
  NousResearchModelLogo, NvidiaModelLogo, PalmModelLogo, PanguModelLogo, PerplexityModelLogo,
  PixtralModelLogo, QwenModelLogo, RakutenaiModelLogo, SparkDeskModelLogo, StabilityModelLogo,
  StepModelLogo, SunoModelLogo, TeleModelLogo, TokenFluxModelLogo, UpstageModelLogo,
  ViduModelLogo, VoyageModelLogo, WenxinModelLogo, XirangModelLogo, YiModelLogo,
} from '@/assets/images/models'

import { YoudaoLogo, NomicLogo } from '@/assets/images/providers'

export function isRerankModel(model_id: string): boolean {
  const modelId = getLowerBaseModelName(model_id)
  return model_id ? /(?:rerank|re-rank|re-ranker|re-ranking|retrieval|retriever)/i.test(modelId) || false : false
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
    return /(?:^text-|embed|bge-|e5-|LLM2Vec|retrieval|uae-|gte-|jina-clip|jina-embeddings|voyage-)/i.test(modelId)
  }

  return /(?:^text-|embed|bge-|e5-|LLM2Vec|retrieval|uae-|gte-|jina-clip|jina-embeddings|voyage-)/i.test(modelId) || false
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
    return new RegExp(
      `\\b(?!(?:${NOT_TOOL_CALL_MODELS.join('|')})\\b)(?:${TOOL_CALL_MODELS.join('|')})\\b`,
      'i'
    ).test(modelId)
  }

  if (['deepseek', 'anthropic'].includes(provider)) {
    return true
  }

  if (['kimi', 'moonshot'].includes(provider)) {
    return true
  }

  return new RegExp(
    `\\b(?!(?:${NOT_TOOL_CALL_MODELS.join('|')})\\b)(?:${TOOL_CALL_MODELS.join('|')})\\b`,
    'i'
  ).test(modelId)
}

export function isVisionModel(model_id: string, provider: string): boolean {
  if (!model_id || isEmbeddingModel(model_id, provider) || isRerankModel(model_id)) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id)
  return new RegExp(
    `\\b(?!(?:${notVisionModels.join('|')})\\b)(${visionModels.join('|')})\\b`,
    'i'
  ).test(modelId) || false
}

export function isCodeModel(model_id: string, provider: string): boolean {
  if (!model_id || isEmbeddingModel(model_id, provider) || isRerankModel(model_id)) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id)
  return /(?:^o3$|.*(code|claude\s+sonnet|claude\s+opus|gpt-4\.1|gpt-4o|gpt-5|gemini[\s-]+2\.5|o4-mini|kimi-k2).*)/i.test(modelId) || false
}

export function isReasoningModel(model_id: string, provider: string): boolean {
  if (!model_id || isEmbeddingModel(model_id, provider) || isRerankModel(model_id)) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id)

  // Check if it's a text-to-image model (merged from isTextToImageModel)
  if (/flux|diffusion|stabilityai|sd-|dall|cogview|janus|midjourney|mj-|image|gpt-image/i.test(modelId)) {
    return false
  }

  if (provider === 'doubao' || modelId.includes('doubao')) {
    return (
      /^(o\d+(?:-[\w-]+)?|.*\b(?:reasoning|reasoner|thinking)\b.*|.*-[rR]\d+.*|.*\bqwq(?:-[\w-]+)?\b.*|.*\bhunyuan-t1(?:-[\w-]+)?\b.*|.*\bglm-zero-preview\b.*|.*\bgrok-(?:3-mini|4)(?:-[\w-]+)?\b.*)$/i.test(modelId) ||
      /doubao-(?:1[.-]5-thinking-vision-pro|1[.-]5-thinking-pro-m|seed-1[.-]6(?:-flash)?(?!-(?:thinking)(?:-|$)))(?:-[\w-]+)*/i.test(getLowerBaseModelName(model_id, '/')) ||
      false
    )
  }

  // Claude reasoning model check (merged from isClaudeReasoningModel)
  // Note: Fixed the bug in original function - it should check !model_id, not model_id
  const claudeModelId = getLowerBaseModelName(model_id, '/')
  const isClaudeReasoning = (
    claudeModelId.includes('claude-3-7-sonnet') ||
    claudeModelId.includes('claude-3.7-sonnet') ||
    claudeModelId.includes('claude-sonnet-4') ||
    claudeModelId.includes('claude-opus-4')
  )

  // OpenAI reasoning model check (merged from isOpenAIReasoningModel and isSupportedReasoningEffortOpenAIModel)
  const openaiModelId = getLowerBaseModelName(model_id, '/')
  const isGPT5Series = modelId.includes('gpt-5')
  const isSupportedReasoningEffortOpenAI = (
    (modelId.includes('o1') && !(modelId.includes('o1-preview') || modelId.includes('o1-mini'))) ||
    modelId.includes('o3') ||
    modelId.includes('o4') ||
    modelId.includes('gpt-oss') ||
    (isGPT5Series && !modelId.includes('chat'))
  )
  const isOpenAIReasoning = isSupportedReasoningEffortOpenAI || openaiModelId.includes('o1')

  // Gemini reasoning model check (merged from isGeminiReasoningModel and isSupportedThinkingTokenGeminiModel)
  const geminiModelId = getLowerBaseModelName(model_id, '/')
  const isSupportedThinkingTokenGemini = geminiModelId.includes('gemini-2.5')
  const isGeminiReasoning = (
    (modelId.startsWith('gemini') && modelId.includes('thinking')) ||
    isSupportedThinkingTokenGemini
  )

  // Qwen reasoning model check (merged from isQwenReasoningModel and isSupportedThinkingTokenQwenModel)
  const qwenModelId = getLowerBaseModelName(model_id, '/')
  let isQwenReasoning = false
  if (qwenModelId.startsWith('qwen3')) {
    if (qwenModelId.includes('thinking')) {
      isQwenReasoning = true
    } else if (qwenModelId.includes('instruct')) {
      isQwenReasoning = false
    } else {
      isQwenReasoning = true
    }
  } else {
    // Check isSupportedThinkingTokenQwenModel logic
    const isSupportedThinkingTokenQwen = !qwenModelId.includes('coder') && [
      'qwen-plus', 'qwen-plus-latest', 'qwen-plus-0428', 'qwen-plus-2025-04-28',
      'qwen-plus-0714', 'qwen-plus-2025-07-14', 'qwen-turbo', 'qwen-turbo-latest',
      'qwen-turbo-0428', 'qwen-turbo-2025-04-28', 'qwen-turbo-0715', 'qwen-turbo-2025-07-15'
    ].includes(qwenModelId)

    isQwenReasoning = isSupportedThinkingTokenQwen || qwenModelId.includes('qwq') || qwenModelId.includes('qvq')
  }

  // Grok reasoning model check (merged from isGrokReasoningModel and isSupportedReasoningEffortGrokModel)
  const isSupportedReasoningEffortGrok = modelId.includes('grok-3-mini')
  const isGrokReasoning = isSupportedReasoningEffortGrok || modelId.includes('grok-4')

  // Hunyuan reasoning model check (merged from isHunyuanReasoningModel and isSupportedThinkingTokenHunyuanModel)
  const hunyuanModelId = getLowerBaseModelName(model_id, '/')
  const isSupportedThinkingTokenHunyuan = hunyuanModelId.includes('hunyuan-a13b')
  const isHunyuanReasoning = isSupportedThinkingTokenHunyuan || hunyuanModelId.includes('hunyuan-t1')

  // Perplexity reasoning model check (merged from isPerplexityReasoningModel and isSupportedReasoningEffortPerplexityModel)
  const perplexityModelId = getLowerBaseModelName(model_id, '/')
  const isSupportedReasoningEffortPerplexity = perplexityModelId.includes('sonar-deep-research')
  const isPerplexityReasoning = isSupportedReasoningEffortPerplexity || perplexityModelId.includes('reasoning')

  // Zhipu reasoning model check (merged from isZhipuReasoningModel and isSupportedThinkingTokenZhipuModel)
  const zhipuModelId = getLowerBaseModelName(model_id, '/')
  const isSupportedThinkingTokenZhipu = zhipuModelId.includes('glm-4.5')
  const isZhipuReasoning = isSupportedThinkingTokenZhipu || zhipuModelId.includes('glm-z1')

  // Step reasoning model check (merged from isStepReasoningModel)
  const stepModelId = getLowerBaseModelName(model_id, '/')
  const isStepReasoning = stepModelId.includes('step-3') || stepModelId.includes('step-r1-v-mini')

  if (
    isClaudeReasoning ||
    isOpenAIReasoning ||
    isGeminiReasoning ||
    isQwenReasoning ||
    isGrokReasoning ||
    isHunyuanReasoning ||
    isPerplexityReasoning ||
    isZhipuReasoning ||
    isStepReasoning ||
    modelId.includes('magistral') ||
    modelId.includes('minimax-m1') ||
    modelId.includes('pangu-pro-moe')
  ) {
    return true
  }

  return /^(o\d+(?:-[\w-]+)?|.*\b(?:reasoning|reasoner|thinking)\b.*|.*-[rR]\d+.*|.*\bqwq(?:-[\w-]+)?\b.*|.*\bhunyuan-t1(?:-[\w-]+)?\b.*|.*\bglm-zero-preview\b.*|.*\bgrok-(?:3-mini|4)(?:-[\w-]+)?\b.*)$/i.test(modelId) || false
}

export function isWebSearchModel(model_id: string, provider: string): boolean {
  if (!provider || !model_id || isEmbeddingModel(model_id, provider) || isRerankModel(model_id)) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id, '/')

  // Anthropic model check (merged from isAnthropicModel)
  if (modelId.startsWith('claude')) {
    return new RegExp(
      `\\b(?:claude-3(-|\\.)(7|5)-sonnet(?:-[\\w-]+)|claude-3(-|\\.)5-haiku(?:-[\\w-]+)|claude-sonnet-4(?:-[\\w-]+)?|claude-opus-4(?:-[\\w-]+)?)\\b`,
      'i'
    ).test(modelId)
  }

  if (provider === 'grok') {
    return true
  }

  if (provider === 'dashscope') {
    const models = ['qwen-turbo', 'qwen-max', 'qwen-plus', 'qwq']
    return models.some((i) => modelId.startsWith(i))
  }

  if (provider === 'gemini' || provider === 'vertexai') {
    return new RegExp('gemini-2\\..*', 'i').test(modelId)
  }

  if (provider === 'openai') {
    // OpenAI web search model check (merged from isOpenAIWebSearchModel)
    if (
      new RegExp('gemini-2\\..*', 'i').test(modelId) ||
      modelId.includes('gpt-4o-search-preview') ||
      modelId.includes('gpt-4o-mini-search-preview') ||
      (modelId.includes('gpt-4.1') && !modelId.includes('gpt-4.1-nano')) ||
      (modelId.includes('gpt-4o') && !modelId.includes('gpt-4o-image')) ||
      modelId.includes('o3') ||
      modelId.includes('o4') ||
      (modelId.includes('gpt-5') && !modelId.includes('chat'))
    ) {
      return true
    }
  }

  if (provider === 'zhipu') {
    return modelId?.startsWith('glm-4-')
  }

  if (provider === 'openrouter') {
    return true
  }

  if (provider === 'hunyuan') {
    return modelId !== 'hunyuan-lite'
  }

  if (provider === 'perplexity') {
    return PERPLEXITY_MODELS.includes(modelId)
  }

  if (provider === 'aihubmix') {
    if (!modelId.endsWith('-search') && new RegExp('gemini-2\\..*', 'i').test(modelId)) {
      return true
    }

    // OpenAI web search model check (merged from isOpenAIWebSearchModel)
    if (
      modelId.includes('gpt-4o-search-preview') ||
      modelId.includes('gpt-4o-mini-search-preview') ||
      (modelId.includes('gpt-4.1') && !modelId.includes('gpt-4.1-nano')) ||
      (modelId.includes('gpt-4o') && !modelId.includes('gpt-4o-image')) ||
      modelId.includes('o3') ||
      modelId.includes('o4') ||
      (modelId.includes('gpt-5') && !modelId.includes('chat'))
    ) {
      return true
    }

    return false
  }

  return false
}

export function isAnalysisModel(model_id: string, provider: string): boolean {
  if (!model_id || isEmbeddingModel(model_id, provider) || isRerankModel(model_id)) {
    return false
  }

  const modelId = getLowerBaseModelName(model_id, '/')

  const bMatch = modelId.match(/(\d+(?:\.\d+)?)\s*(?:b|bn)\b/i)
  if (bMatch) {
    const n = parseFloat(bMatch[1])
    if (!isNaN(n) && n <= 7) {
      return true
    }
    return false
  }

  const mMatch = modelId.match(/(\d+(?:\.\d+)?)\s*m\b/i)
  if (mMatch) {
    const n = parseFloat(mMatch[1])
    if (!isNaN(n) && n > 0 && n <= 7000) {
      return true
    }
  }

  return false
}

export function getModelLogo(modelId: string) {
  if (!modelId) {
    return undefined
  }

  const logoMap = {
    pixtral: PixtralModelLogo,
    jina: JinaModelLogo,
    abab: MinimaxModelLogo,
    minimax: MinimaxModelLogo,
    o1: ChatGPTo1ModelLogo,
    o3: ChatGPTo1ModelLogo,
    o4: ChatGPTo1ModelLogo,
    'gpt-image': ChatGPTImageModelLogo,
    'gpt-3': ChatGPT35ModelLogo,
    'gpt-4': ChatGPT4ModelLogo,
    'gpt-5-mini': GPT5MiniModelLogo,
    'gpt-5-nano': GPT5NanoModelLogo,
    'gpt-5-chat': GPT5ChatModelLogo,
    'gpt-5': GPT5ModelLogo,
    gpts: ChatGPT4ModelLogo,
    'gpt-oss(?:-[\\w-]+)': ChatGptModelLogo,
    'text-moderation': ChatGptModelLogo,
    'babbage-': ChatGptModelLogo,
    'sora-': ChatGptModelLogo,
    '(^|/)omni-': ChatGptModelLogo,
    'Embedding-V1': WenxinModelLogo,
    'text-embedding-v': QwenModelLogo,
    'text-embedding': ChatGptModelLogo,
    'davinci-': ChatGptModelLogo,
    glm: ChatGLMModelLogo,
    deepseek: DeepSeekModelLogo,
    '(qwen|qwq|qwq-|qvq-)': QwenModelLogo,
    gemma: GemmaModelLogo,
    'yi-': YiModelLogo,
    llama: LlamaModelLogo,
    mixtral: MistralModelLogo,
    mistral: MistralModelLogo,
    codestral: CodestralModelLogo,
    ministral: MistralModelLogo,
    magistral: MistralModelLogo,
    moonshot: MoonshotModelLogo,
    kimi: MoonshotModelLogo,
    phi: MicrosoftModelLogo,
    baichuan: BaichuanModelLogo,
    claude: ClaudeModelLogo,
    gemini: GeminiModelLogo,
    bison: PalmModelLogo,
    palm: PalmModelLogo,
    step: StepModelLogo,
    hailuo: HailuoModelLogo,
    doubao: DoubaoModelLogo,
    'ep-202': DoubaoModelLogo,
    cohere: CohereModelLogo,
    command: CohereModelLogo,
    minicpm: MinicpmModelLogo,
    '360': Ai360ModelLogo,
    aimass: AimassModelLogo,
    codegeex: CodegeexModelLogo,
    copilot: CopilotModelLogo,
    creative: CopilotModelLogo,
    balanced: CopilotModelLogo,
    precise: CopilotModelLogo,
    dalle: DalleModelLogo,
    'dall-e': DalleModelLogo,
    dbrx: DbrxModelLogo,
    flashaudio: FlashaudioModelLogo,
    flux: FluxModelLogo,
    grok: GrokModelLogo,
    hunyuan: HunyuanModelLogo,
    internlm: InternlmModelLogo,
    internvl: InternvlModelLogo,
    llava: LLavaModelLogo,
    magic: MagicModelLogo,
    midjourney: MidjourneyModelLogo,
    'mj-': MidjourneyModelLogo,
    'tao-': WenxinModelLogo,
    'ernie-': WenxinModelLogo,
    voice: FlashaudioModelLogo,
    'tts-1': ChatGptModelLogo,
    'whisper-': ChatGptModelLogo,
    'stable-': StabilityModelLogo,
    sd2: StabilityModelLogo,
    sd3: StabilityModelLogo,
    sdxl: StabilityModelLogo,
    sparkdesk: SparkDeskModelLogo,
    generalv: SparkDeskModelLogo,
    wizardlm: MicrosoftModelLogo,
    microsoft: MicrosoftModelLogo,
    hermes: NousResearchModelLogo,
    gryphe: GrypheModelLogo,
    suno: SunoModelLogo,
    chirp: SunoModelLogo,
    luma: LumaModelLogo,
    keling: KeLingModelLogo,
    'vidu-': ViduModelLogo,
    ai21: Ai21ModelLogo,
    'jamba-': Ai21ModelLogo,
    mythomax: GrypheModelLogo,
    nvidia: NvidiaModelLogo,
    dianxin: DianxinModelLogo,
    tele: TeleModelLogo,
    adept: AdeptModelLogo,
    aisingapore: AisingaporeModelLogo,
    bigcode: BigcodeModelLogo,
    mediatek: MediatekModelLogo,
    upstage: UpstageModelLogo,
    rakutenai: RakutenaiModelLogo,
    ibm: IbmModelLogo,
    'google/': GoogleModelLogo,
    xirang: XirangModelLogo,
    hugging: HuggingfaceModelLogo,
    youdao: YoudaoLogo,
    embedding: EmbeddingModelLogo,
    perplexity: PerplexityModelLogo,
    sonar: PerplexityModelLogo,
    'bge-': BgeModelLogo,
    'voyage-': VoyageModelLogo,
    tokenflux: TokenFluxModelLogo,
    'nomic-': NomicLogo,
    'pangu-': PanguModelLogo
  }

  for (const key in logoMap) {
    const regex = new RegExp(key, 'i')
    if (regex.test(modelId)) {
      return logoMap[key as keyof typeof logoMap]
    }
  }

  return undefined
}

export const getModelGroup = (id: string, provider?: string): string => {
  const str = id.toLowerCase()

  let delimiters = ['/', ' ', ':']
  let anotherDelimiters = ['-', '_']

  if (provider && ['aihubmix', 'silicon', 'ocoolai', 'o3', 'dmxapi'].includes(provider.toLowerCase())) {
    delimiters = ['/', ' ', '-', '_', ':']
    anotherDelimiters = []
  }

  for (const dlmtr of delimiters) {
    if (str.includes(dlmtr)) {
      return str.split(dlmtr)[0]
    }
  }

  for (const dlmtr of anotherDelimiters) {
    if (str.includes(dlmtr)) {
      const parts = str.split(dlmtr)
      return parts.length > 1 ? parts[0] + '-' + parts[1] : parts[0]
    }
  }

  return str
}
