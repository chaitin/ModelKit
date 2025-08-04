package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"

	"github.com/chaitin/ModelKit/backend/consts"
	"github.com/chaitin/ModelKit/backend/domain"
	"github.com/chaitin/ModelKit/backend/pkg/request"
)

type ModelUsecase struct {
	client *http.Client
}

func NewModelUsecase(
) domain.ModelUsecase {
	client := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			MaxConnsPerHost:     100,
			IdleConnTimeout:     time.Second * 30,
		},
	}
	return &ModelUsecase{client: client}
}

func (m *ModelUsecase) CheckModel(ctx context.Context, req *domain.CheckModelReq) (*domain.Model, error) {
	if req.SubType == consts.ModelTypeEmbedding || req.SubType == consts.ModelTypeReranker {
		url := domain.ModelOwners[req.Owner].APIBase
		reqBody := map[string]any{}
		if req.SubType == consts.ModelTypeEmbedding {
			reqBody = map[string]any{
				"model":           req.ModelID,
				"input":           "ModelKit一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
				"encoding_format": "float",
			}
			url = domain.ModelOwners[req.Owner].APIBase + "/embeddings"
		}
		if req.SubType == consts.ModelTypeReranker {
			reqBody = map[string]any{
				"model": req.ModelID,
				"documents": []string{
					"ModelKit一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
					"ModelKit一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
					"ModelKit一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
				},
				"query": "ModelKit",
			}
			url = domain.ModelOwners[req.Owner].APIBase + "/rerank"
		}
		body, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", req.APIKey))
		request.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("request failed: %s", resp.Status)
		}
		return &domain.Model{}, nil
	}
	config := &openai.ChatModelConfig{
		APIKey:  req.APIKey,
		BaseURL: domain.ModelOwners[req.Owner].APIBase,
		Model:   string(req.ModelID),
	}
	fmt.Println("BaseURL", domain.ModelOwners[req.Owner].APIBase)
	// for azure openai
	if req.Owner == consts.ModelOwnerAzureOpenAI {
		config.ByAzure = true
		config.APIVersion = domain.ModelOwners[req.Owner].APIVersion
		if config.APIVersion == "" {
			config.APIVersion = "2024-10-21"
		}
	}
	// end
	if domain.ModelOwners[req.Owner].APIHeader != "" {
		client := getHttpClientWithAPIHeaderMap(domain.ModelOwners[req.Owner].APIHeader)
		if client != nil {
			config.HTTPClient = client
		}
	}
	chatModel, err := openai.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	resp, err := chatModel.Generate(ctx, []*schema.Message{
		schema.SystemMessage("You are a helpful assistant."),
		schema.UserMessage("hi"),
	})
	if err != nil {
		return nil, err
	}
	content := resp.Content
	if content == "" {
		return nil, fmt.Errorf("generate failed")
	}
	return &domain.Model{
		ModelType: req.SubType,
		OwnedBy:   req.Owner,
		ID:        req.ModelID,
	}, nil
}

type headerTransport struct {
	headers map[string]string
	base    http.RoundTripper
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Set(k, v)
	}
	return t.base.RoundTrip(req)
}

func getHttpClientWithAPIHeaderMap(header string) *http.Client {
	headerMap := request.GetHeaderMap(header)
	if len(headerMap) > 0 {
		// create http client with custom transport for headers
		client := &http.Client{
			Timeout: 0,
		}
		// Wrap the transport to add headers
		client.Transport = &headerTransport{
			headers: headerMap,
			base:    http.DefaultTransport,
		}
		return client
	}
	return nil
}

// ListModel implements domain.ModelUsecase.
func (m *ModelUsecase) ListModel(ctx context.Context, req *domain.ListModelReq) ([]*domain.Model, error) {
	// 如果没有请求参数或参数为空，返回全体模型
	if req == nil || (req.OwnedBy == "" && req.SubType == "") {
		result := make([]*domain.Model, len(domain.Models))
		for i := range domain.Models {
			result[i] = &domain.Models[i]
		}
		return result, nil
	}

	var models []*domain.Model

	// 只有 Owner 参数
	if req.OwnedBy != "" && req.SubType == "" {
		if owner, exists := domain.ModelOwners[req.OwnedBy]; exists {
			models = owner.Models
		}
	} else if req.OwnedBy == "" && req.SubType != "" {
		// 只有 Type 参数
		if typeModels, exists := domain.TypeModelMap[req.SubType]; exists {
			models = typeModels
		}
	} else {
		// 同时有 Owner 和 Type 参数，需要取交集
		ownerModels, ownerExists := domain.ModelOwners[req.OwnedBy]
		typeModels, typeExists := domain.TypeModelMap[req.SubType]

		if ownerExists && typeExists {
			// 构建一个map用于快速查找
			typeModelMap := make(map[string]bool)
			for _, model := range typeModels {
				typeModelMap[model.ID] = true
			}

			// 找出交集
			for _, model := range ownerModels.Models {
				if typeModelMap[model.ID] {
					models = append(models, model)
				}
			}
		}
	}

	return models, nil
}
