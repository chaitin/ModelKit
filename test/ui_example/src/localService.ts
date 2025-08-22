import {
  Model,
  ModelService,
  CreateModelReq,
  ListModelReq,
  CheckModelReq,
  UpdateModelReq,
  ModelListItem,
} from '@yokowu/modelkit-ui';


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
  error?: string;
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
      base_url: data.base_url,
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

  async listModel(data: ListModelReq): Promise<{ models: ModelListItem[]; error?: string }> {
      const queryParams = new URLSearchParams();
      if (data.provider) queryParams.append('provider', data.provider);
      if (data.model_type) queryParams.append('model_type', data.model_type);
      if (data.base_url) queryParams.append('base_url', data.base_url);
      if (data.api_key) queryParams.append('api_key', data.api_key);
      if (data.api_header) queryParams.append('api_header', data.api_header);
      
      const url = `/listmodel${queryParams.toString() ? '?' + queryParams.toString() : ''}`;
      const response = await this.request<ModelListResponse>(url, {
        method: 'GET',
      });
      console.log('listModel response:', response);
      return { models: response.models, error: response.error };
  }

  async checkModel(data: CheckModelReq): Promise<{ model: Model; error?: string }> {
      const queryParams = new URLSearchParams();
      if (data.provider) queryParams.append('provider', data.provider);
      if (data.model_name) queryParams.append('model_name', data.model_name);
      if (data.base_url) queryParams.append('base_url', data.base_url);
      if (data.api_key) queryParams.append('api_key', data.api_key);
      if (data.api_header) queryParams.append('api_header', data.api_header);
      if (data.model_type) queryParams.append('model_type', data.model_type);
      const url = `/checkmodel${queryParams.toString() ? '?' + queryParams.toString() : ''}`;
      console.log('checkModel url:', url);
      const response = await this.request<CheckModelResponse>(url, {
        method: 'GET',
      });
      console.log('checkModel response:', response);
      return { model: response.model, error: response.error };
  }

  async updateModel(data: UpdateModelReq): Promise<{ model: Model }> {
    const requestData = {
      id: data.id,
      model_name: data.model_name,
      show_name: data.show_name,
      provider: data.provider,
      base_url: data.base_url,
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