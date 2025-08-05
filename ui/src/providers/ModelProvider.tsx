import React, { createContext, useContext, useState, useEffect } from 'react';
import { DomainModel, ConstsModelType } from '@/api/types';
import { getListModel } from '@/api/Model';

interface ModelContextType {
  coderModel: DomainModel[];
  llmModel: DomainModel[];
  isConfigModel: boolean;
  modelLoading: boolean;
  refreshModel: () => void;
}

const ModelContext = createContext<ModelContextType | undefined>(undefined);

export const useModelContext = () => {
  const context = useContext(ModelContext);
  if (!context) {
    throw new Error('useModelContext must be used within a ModelProvider');
  }
  return context;
};

interface ModelProviderProps {
  children: React.ReactNode;
}

export const ModelProvider: React.FC<ModelProviderProps> = ({ children }) => {
  const [coderModel, setCoderModel] = useState<DomainModel[]>([]);
  const [llmModel, setLlmModel] = useState<DomainModel[]>([]);
  const [modelLoading, setModelLoading] = useState(false);

  const fetchModels = async () => {
    setModelLoading(true);
    try {
      // 使用现有的getListModel接口获取模型列表
      const models = await getListModel({ sub_type: ConstsModelType.ModelTypeChat }) || [];
      // 确保 models 是数组类型
      const modelArray = Array.isArray(models) ? models : [];

      // 根据模型类型分离数据
      const llmModels = modelArray.filter((model: any) => model.sub_type === ConstsModelType.ModelTypeChat);
      const coderModels = modelArray.filter((model: any) => model.sub_type === ConstsModelType.ModelTypeCoder);

      setLlmModel(llmModels);
      setCoderModel(coderModels);
    } catch (error) {
      console.error('Failed to fetch models:', error);
      // 如果API调用失败，使用空数据
      setLlmModel([]);
      setCoderModel([]);
    } finally {
      setModelLoading(false);
    }
  };

  useEffect(() => {
    fetchModels();
  }, []);

  const refreshModel = () => {
    fetchModels();
  };

  const value: ModelContextType = {
    coderModel,
    llmModel,
    modelLoading,
    refreshModel,
    isConfigModel: false
  };

  return (
    <ModelContext.Provider value={value}>
      {children}
    </ModelContext.Provider>
  );
};