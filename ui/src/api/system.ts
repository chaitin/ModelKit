import request, { ContentType } from './httpClient';
// 暂时注释掉类型导入，直接使用any类型
// import type { 
//   SystemSettings, 
//   RequestParams,
//   ApiResponse 
// } from '@/types';

// 获取系统设置
export const getSystemSettings = (params: any = {}) =>
  request<any>({
    path: '/system/settings',
    method: 'GET',
    ...params,
  });

// 更新系统设置
export const updateSystemSettings = (
  param: any,
  params: any = {},
) =>
  request<any>({
    path: '/system/settings',
    method: 'PUT',
    body: param,
    type: ContentType.Json,
    ...params,
  });