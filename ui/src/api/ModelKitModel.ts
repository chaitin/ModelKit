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

import request, { ContentType, RequestParams } from "./httpClient";
import {
  DomainCheckModelReq,
  DomainModel,
  DomainResp,
  GetListModelParams,
} from "./types";

/**
 * @description 检查模型
 *
 * @tags ModelKitModel
 * @name PostCheckModel
 * @summary 检查模型
 * @request POST:/api/v1/model/modelkit/check
 * @response `200` `(DomainResp & {
    data?: DomainModel,

})` OK
 */

export const postCheckModel = (
  model: DomainCheckModelReq,
  params: RequestParams = {},
) =>
  request<
    DomainResp & {
      data?: DomainModel;
    }
  >({
    path: `/api/v1/model/modelkit/check`,
    method: "POST",
    body: model,
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 根据筛选条件获取模型列表，支持分页，按创建时间降序排列
 *
 * @tags ModelKitModel
 * @name GetListModel
 * @summary 获取模型列表
 * @request GET:/api/v1/model/modelkit/models
 * @response `200` `(DomainResp & {
    data?: (DomainModel)[],

})` OK
 */

export const getListModel = (
  query: GetListModelParams,
  params: RequestParams = {},
) =>
  request<
    DomainResp & {
      data?: DomainModel[];
    }
  >({
    path: `/api/v1/model/modelkit/models`,
    method: "GET",
    query: query,
    type: ContentType.Json,
    format: "json",
    ...params,
  });
