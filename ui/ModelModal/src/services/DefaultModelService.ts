import {
  ModelService,
  CreateModelData,
  GetModelNameData,
  CheckModelData,
  UpdateModelData,
} from '../types';

/**
 * 默认的ModelService实现
 * 可以通过继承此类来创建自定义实现
 */
export class DefaultModelService implements ModelService {
  private baseURL: string;
  private headers: Record<string, string>;

  constructor(baseURL: string = '', headers: Record<string, string> = {}) {
    this.baseURL = baseURL;
    this.headers = {
      'Content-Type': 'application/json',
      ...headers,
    };
  }

  /**
   * 创建模型
   */
  async createModel(data: CreateModelData): Promise<{ id: string }> {
    const response = await fetch(`${this.baseURL}/api/models`, {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error(`创建模型失败: ${response.statusText}`);
    }

    return response.json();
  }

  /**
   * 获取模型名称列表
   */
  async getModelNameList(data: GetModelNameData): Promise<{ models: { model: string }[] }> {
    const response = await fetch(`${this.baseURL}/api/models/list`, {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error(`获取模型列表失败: ${response.statusText}`);
    }

    return response.json();
  }

  /**
   * 测试模型
   */
  async testModel(data: CheckModelData): Promise<{ error: string }> {
    const response = await fetch(`${this.baseURL}/api/models/test`, {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error(`测试模型失败: ${response.statusText}`);
    }

    return response.json();
  }

  /**
   * 更新模型
   */
  async updateModel(data: UpdateModelData): Promise<void> {
    const response = await fetch(`${this.baseURL}/api/models/${data.id}`, {
      method: 'PUT',
      headers: this.headers,
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error(`更新模型失败: ${response.statusText}`);
    }
  }

  /**
   * 设置基础URL
   */
  setBaseURL(url: string): void {
    this.baseURL = url;
  }

  /**
   * 设置请求头
   */
  setHeaders(headers: Record<string, string>): void {
    this.headers = { ...this.headers, ...headers };
  }

  /**
   * 添加认证头
   */
  setAuthToken(token: string): void {
    this.headers.Authorization = `Bearer ${token}`;
  }
}

/**
 * 创建默认的ModelService实例
 */
export const createDefaultModelService = (
  baseURL: string = '',
  headers: Record<string, string> = {}
): DefaultModelService => {
  return new DefaultModelService(baseURL, headers);
}; 