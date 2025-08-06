/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export enum ConstsModelType {
  ModelTypeChat = "chat",
  ModelTypeCoder = "coder",
  ModelTypeEmbedding = "embedding",
  ModelTypeRerank = "reranker",
  ModelTypeVision = "vision",
  ModelTypeFunctionCall = "function_call",
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
  ModelProviderGemini = "Gemini",
  ModelProviderZhiPu = "ZhiPu",
}

export interface DomainCheckModelReq {
  /** 接口密钥 */
  api_key: string;
  /** 模型名称 */
  model_name: string;
  /** 提供商 */
  owner: ConstsModelProvider;
  /** 模型类型 */
  sub_type: ConstsModelType;
}

export interface DomainModel {
  /** 创建时间 */
  created?: number;
  /** 模型的名字 */
  id?: string;
  /** 模型类型 */
  model_type?: ConstsModelType;
  /** 总是model */
  object?: string;
  /** 提供商 */
  owned_by?: ConstsModelProvider;
}

export interface DomainResp {
  code?: number;
  data?: unknown;
  message?: string;
}

export interface GetListModelParams {
  /** 提供商 */
  owned_by?:
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
    | "Gemini"
    | "ZhiPu";
  /** 模型类型 */
  sub_type?:
    | "chat"
    | "coder"
    | "embedding"
    | "reranker"
    | "vision"
    | "function_call";
}
