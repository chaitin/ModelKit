# ModelKit

[![Go Version](https://img.shields.io/badge/Go-1.24.0+-blue.svg)](https://golang.org)
[![React Version](https://img.shields.io/badge/React-19.0.0+-blue.svg)](https://reactjs.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

ModelKit 是一个强大的AI模型管理平台，支持多种AI服务提供商，提供统一的模型管理、配置验证服务。

## 🚀 功能特性

- **多模型提供商支持**: 支持 OpenAI、Ollama、DeepSeek、SiliconFlow、Moonshot、Azure OpenAI、百智云、腾讯混元、百炼、火山引擎、Gemini、智谱等主流AI服务商
- **模型类型管理**: 支持聊天模型、嵌入模型、重排序模型、视觉模型、代码模型、函数调用等多种模型类型
- **配置验证**: 提供模型配置的实时验证功能，确保API配置正确性
- **现代化Web界面**: 基于React 19和Material-UI构建的响应式用户界面
- **国际化支持**: 内置中英文多语言支持
- **可复用组件**: 提供开箱即用的ModelModal组件，支持在其他项目中快速集成

## 使用方式
### 后端
``` bash
    // 1. 引入ModelKit
	import (
		modelkit "github.com/chaitin/ModelKit/usecase"
	)
    // 2. 调用ModelKit提供的函数即可
    modelkitRes, err := modelkit.CheckModel(...)
    modelkitRes, err := modelkit.ListModel(...)
```
### 前端
``` bash
    // 1. 引入ModelKit
    import { ModelModal, Model, ModelService, ConstsModelType as ModelKitType, ModelListItem } from '@yokowu/modelkit-ui';
    // 2.创建符合ModelService接口的服务实现
    const modelService: ModelService = {
    createModel: async (params) => {
      const response = await postCreateModel(params as unknown as DomainCreateModelReq);
      return { model: response as unknown as Model };
    },
    listModel: async (params) => {
      const response = await getGetProviderModelList(params as unknown as GetGetProviderModelListParams);
      return { models: response?.models || [] };
    },
    checkModel: async (params) => {
      const response = await postCheckModel(params as unknown as DomainCheckModelReq);
      return { model: response as unknown as Model };
    },
    updateModel: async (params) => {
      const response = await putUpdateModel(params as unknown as DomainUpdateModelReq);
      return { model: response as unknown as Model };
    }
    // 3. 使用ModelModal组件
        <ModelModal
        open={open}
        onClose={() => {
          setOpen(false);
          setEditData(null);
        }}
        refresh={refreshModel}
        data={editData as Model | null}
        type={modelType}
        modelService={modelService}
        language="zh-CN"
      />
```