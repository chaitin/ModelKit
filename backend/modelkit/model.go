package modelkit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"github.com/chaitin/ModelKit/backend/consts"
	"github.com/chaitin/ModelKit/backend/domain"
	"github.com/chaitin/ModelKit/backend/pkg/request"
)

type ModelKit struct {
	client *http.Client
}

func NewModelKit() domain.ModelKit {
	client := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			MaxConnsPerHost:     100,
			IdleConnTimeout:     time.Second * 30,
		},
	}
	return &ModelKit{client: client}
}

func (m *ModelKit) CheckModel(ctx context.Context, req *domain.CheckModelReq) (*domain.Model, error) {
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
func (m *ModelKit) ListModel(ctx context.Context, req *domain.ListModelReq) ([]*domain.Model, error) {
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

func (m *ModelKit) PandaModelList(ctx context.Context, req *domain.GetProviderModelListReq) (*domain.GetProviderModelListResp, error) {
	switch provider := consts.ModelOwner(req.Provider); provider {
	case consts.ModelOwnerMoonshot,
		consts.ModelOwnerDeepSeek,
		consts.ModelOwnerAzureOpenAI,
		consts.ModelOwnerVolcengine,
		consts.ModelOwnerZhiPu:
		return &domain.GetProviderModelListResp{
			Models: domain.From(domain.ModelOwners[consts.ModelOwner(req.Provider)]),
		}, nil
	case consts.ModelOwnerGemini:
		client, err := genai.NewClient(ctx, option.WithAPIKey(req.APIKey))
		if err != nil {
			return nil, err
		}
		defer client.Close()

		modelsList := make([]domain.ProviderModelListItem, 0)
		modelsIter := client.ListModels(ctx)
		for {
			model, err := modelsIter.Next()
			if err != nil {
				break
			}

			if !slices.Contains(model.SupportedGenerationMethods, "generateContent") {
				continue
			}

			if !strings.Contains(model.Name, "gemini") {
				continue
			}

			name, _ := strings.CutPrefix(model.Name, "models/")
			modelsList = append(modelsList, domain.ProviderModelListItem{
				Model: name,
			})
		}

		if len(modelsList) == 0 {
			return nil, fmt.Errorf("failed to get gemini models")
		}

		return &domain.GetProviderModelListResp{
			Models: modelsList,
		}, nil

	case consts.ModelOwnerOpenAI, consts.ModelOwnerHunyuan, consts.ModelOwnerBaiLian:
		u, err := url.Parse(req.BaseURL)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, "/models")
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", req.APIKey))
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get models: %s", resp.Status)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		type OpenAIResp struct {
			Object string `json:"object"`
			Data   []struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		var models OpenAIResp
		err = json.Unmarshal(body, &models)
		if err != nil {
			return nil, err
		}
		modelsList := make([]domain.ProviderModelListItem, 0)
		for _, model := range models.Data {
			modelsList = append(modelsList, domain.ProviderModelListItem{
				Model: model.ID,
			})
		}
		return &domain.GetProviderModelListResp{
			Models: modelsList,
		}, nil
	case consts.ModelOwnerOllama:
		// get from ollama http://10.10.16.24:11434/api/tags
		u, err := url.Parse(req.BaseURL)
		if err != nil {
			return nil, err
		}
		u.Path = "/api/tags"
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}
		if req.APIHeader != "" {
			headers := getHeaderMap(req.APIHeader)
			for k, v := range headers {
				request.Header.Set(k, v)
			}
		}
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var models domain.GetProviderModelListResp
		err = json.Unmarshal(body, &models)
		if err != nil {
			return nil, err
		}
		return &models, nil
	case consts.ModelOwnerSiliconFlow, consts.ModelOwnerBaiZhiCloud:
		if req.Type == consts.ModelTypeEmbedding || req.Type == consts.ModelTypeReranker {
			if provider == consts.ModelOwnerBaiZhiCloud {
				if req.Type == consts.ModelTypeEmbedding {
					return &domain.GetProviderModelListResp{
						Models: []domain.ProviderModelListItem{
							{
								Model: "bge-m3",
							},
						},
					}, nil
				} else {
					return &domain.GetProviderModelListResp{
						Models: []domain.ProviderModelListItem{
							{
								Model: "bge-reranker-v2-m3",
							},
						},
					}, nil
				}
			}
		}
		u, err := url.Parse(req.BaseURL)
		if err != nil {
			return nil, err
		}
		u.Path = "/v1/models"
		q := u.Query()
		q.Set("type", "text")
		q.Set("sub_type", "chat")
		u.RawQuery = q.Encode()
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", req.APIKey))
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get models: %s", resp.Status)
		}
		type SiliconFlowModelResp struct {
			Object string `json:"object"`
			Data   []struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		var models SiliconFlowModelResp
		err = json.Unmarshal(body, &models)
		if err != nil {
			return nil, err
		}
		modelsList := make([]domain.ProviderModelListItem, 0, len(models.Data))
		for _, model := range models.Data {
			modelsList = append(modelsList, domain.ProviderModelListItem{
				Model: model.ID,
			})
		}
		return &domain.GetProviderModelListResp{
			Models: modelsList,
		}, nil
	default:
		return nil, fmt.Errorf("invalid provider: %s", req.Provider)
	}
}

func getHeaderMap(header string) map[string]string {
	headerMap := make(map[string]string)
	for _, h := range strings.Split(header, "\n") {
		if key, value, ok := strings.Cut(h, "="); ok {
			headerMap[key] = value
		}
	}
	return headerMap
}
