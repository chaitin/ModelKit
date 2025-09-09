package usecase

import (
	"github.com/chaitin/ModelKit/v2/consts"
	"github.com/chaitin/ModelKit/v2/domain"
)

type OpenAI struct{}

func NewOpenAI() domain.IModelProvider[domain.ModelMetadata] {
	return &OpenAI{}
}

func (o *OpenAI) ListModel(subType string, provider string) ([]domain.ModelMetadata, error) {
	// 将 subType 和 provider 转换为对应类型
	modelSubType := consts.ModelType(subType)
	modelProvider := consts.ModelProvider(provider)

	// 如果没有请求参数或参数为空，返回全体模型
	if modelProvider == "" && modelSubType == "" {
		result := make([]domain.ModelMetadata, len(domain.Models))
		copy(result, domain.Models)
		return result, nil
	}

	var models []domain.ModelMetadata

	// 只有 Owner 参数
	if modelProvider != "" && modelSubType == "" {
		if owner, exists := domain.ModelProviders[modelProvider]; exists {
			models = owner.Models
		}
	} else if modelProvider == "" && modelSubType != "" {
		// 只有 Type 参数
		if typeModels, exists := domain.TypeModelMap[modelSubType]; exists {
			models = typeModels
		}
	} else {
		// 同时有 Owner 和 Type 参数，需要取交集
		ownerModels, ownerExists := domain.ModelProviders[modelProvider]
		typeModels, typeExists := domain.TypeModelMap[modelSubType]

		if ownerExists && typeExists {
			// 构建一个map用于快速查找
			typeModelMap := make(map[string]bool)
			for _, model := range typeModels {
				typeModelMap[model.ModelName] = true
			}

			// 找出交集
			for _, model := range ownerModels.Models {
				if typeModelMap[model.ModelName] {
					models = append(models, model)
				}
			}
		}
	}

	return models, nil
}
