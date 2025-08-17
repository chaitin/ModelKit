import {
  Model,
  ModelService,
  CreateModelReq,
  ListModelReq,
  CheckModelReq,
  UpdateModelReq,
  ModelListItem,
  ConstsModelProvider,
} from '../../ui/ModelModal/src/types/types';

// 本地Go服务的基础URL
const BASE_URL = 'http://localhost:8080/api/v1/modelkit';

// API响应格式
interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
}

// 模型列表响应
interface ModelListResponse {
  models: ModelListItem[];
}

// 模型检查响应
interface CheckModelResponse {
  model: Model;
  error?: string;
}

export class LocalModelService implements ModelService {
  private async request<T>(url: string, options: RequestInit = {}): Promise<T> {
    const response = await fetch(`${BASE_URL}${url}`, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result: ApiResponse<T> = await response.json();
    
    if (!result.success) {
      throw new Error(result.message || 'API request failed');
    }

    return result.data as T;
  }

  async createModel(data: CreateModelReq): Promise<{ model: Model }> {
    // 转换为Go服务期望的格式
    const requestData = {
      model_name: data.model_name,
      show_name: data.show_name,
      provider: data.provider,
      model_type: data.model_type,
      api_base: data.api_base,
      api_key: data.api_key,
      api_version: data.api_version,
      api_header: data.api_header,
      param: data.param,
    };

    const result = await this.request<{ model: Model }>('/models', {
      method: 'POST',
      body: JSON.stringify(requestData),
    });

    return result;
  }

  async listModel(data: ListModelReq): Promise<{ models: ModelListItem[] }> {
    try {
      const queryParams = new URLSearchParams();
      if (data.provider) queryParams.append('provider', data.provider);
      if (data.type) queryParams.append('type', data.type);
      if (data.base_url) queryParams.append('base_url', data.base_url);
      if (data.api_key) queryParams.append('api_key', data.api_key);
      if (data.api_header) queryParams.append('api_header', data.api_header);
      
      const url = `/listmodel${queryParams.toString() ? '?' + queryParams.toString() : ''}`;
      const response = await this.request<ModelListResponse>(url, {
        method: 'GET',
      });
      return { models: response.models };
    } catch (error) {
      console.error('Failed to list models:', error);
      throw error;
    }
  }

  async checkModel(data: CheckModelReq): Promise<{ model: Model }> {
    try {
      const queryParams = new URLSearchParams();
      if (data.provider) queryParams.append('provider', data.provider);
      if (data.model_name) queryParams.append('model', data.model_name);
      if (data.api_base) queryParams.append('base_url', data.api_base);
      if (data.api_key) queryParams.append('api_key', data.api_key);
      if (data.api_header) queryParams.append('api_header', data.api_header);
      if (data.api_version) queryParams.append('api_version', data.api_version);
      if (data.type) queryParams.append('type', data.type);
      
      const url = `/checkmodel${queryParams.toString() ? '?' + queryParams.toString() : ''}`;
      const response = await this.request<CheckModelResponse>(url, {
        method: 'GET',
      });
      
      if (response.error) {
        throw new Error(response.error);
      }
      
      return { model: response.model };
    } catch (error) {
      console.error('Failed to check model:', error);
      throw error;
    }
  }

  async updateModel(data: UpdateModelReq): Promise<{ model: Model }> {
    const requestData = {
      id: data.id,
      model_name: data.model_name,
      show_name: data.show_name,
      provider: data.provider,
      api_base: data.api_base,
      api_key: data.api_key,
      api_version: data.api_version,
      api_header: data.api_header,
      param: data.param,
      status: data.status,
    };

    const result = await this.request<{ model: Model }>(`/models/${data.id}`, {
      method: 'PUT',
      body: JSON.stringify(requestData),
    });

    return result;
  }
}

// 导出服务实例
export const localModelService = new LocalModelService();