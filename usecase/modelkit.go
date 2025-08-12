package usecase

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

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	genai2 "github.com/google/generative-ai-go/genai"
	"github.com/ollama/ollama/api"
	"google.golang.org/api/option"
	"google.golang.org/genai"

	"github.com/chaitin/ModelKit/consts"
	"github.com/chaitin/ModelKit/domain"
	"github.com/chaitin/ModelKit/utils"
)

func ModelList(ctx context.Context, req *domain.ModelListReq) (*domain.ModelListResp, error) {
	switch provider := consts.ModelProvider(req.Provider); provider {
	case consts.ModelProviderMoonshot,
		consts.ModelProviderDeepSeek,
		consts.ModelProviderAzureOpenAI,
		consts.ModelProviderVolcengine,
		consts.ModelProviderZhiPu:
		return &domain.ModelListResp{
			Models: domain.From(domain.ModelProviders[provider]),
		}, nil
	case consts.ModelProviderGemini:
		client, err := genai2.NewClient(ctx, option.WithAPIKey(req.APIKey))
		if err != nil {
			return nil, err
		}
		defer client.Close()

		modelsList := make([]domain.ModelListItem, 0)
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
			modelsList = append(modelsList, domain.ModelListItem{
				Model: name,
			})
		}

		if len(modelsList) == 0 {
			return nil, fmt.Errorf("failed to get gemini models")
		}

		return &domain.ModelListResp{
			Models: modelsList,
		}, nil

	case consts.ModelProviderOpenAI, consts.ModelProviderHunyuan, consts.ModelProviderBaiLian:
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
		modelsList := make([]domain.ModelListItem, 0)
		for _, model := range models.Data {
			modelsList = append(modelsList, domain.ModelListItem{
				Model: model.ID,
			})
		}
		return &domain.ModelListResp{
			Models: modelsList,
		}, nil
	case consts.ModelProviderOllama:
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
			headers := utils.GetHeaderMap(req.APIHeader)
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
		var models domain.ModelListResp
		err = json.Unmarshal(body, &models)
		if err != nil {
			return nil, err
		}
		return &models, nil
	case consts.ModelProviderSiliconFlow, consts.ModelProviderBaiZhiCloud:
		modelType := consts.ModelType(req.Type)
		if modelType == consts.ModelTypeEmbedding || modelType == consts.ModelTypeRerank {
			if provider == consts.ModelProviderBaiZhiCloud {
				if modelType == consts.ModelTypeEmbedding {
					return &domain.ModelListResp{
						Models: []domain.ModelListItem{
							{
								Model: "bge-m3",
							},
						},
					}, nil
				} else {
					return &domain.ModelListResp{
						Models: []domain.ModelListItem{
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
		modelsList := make([]domain.ModelListItem, 0, len(models.Data))
		for _, model := range models.Data {
			modelsList = append(modelsList, domain.ModelListItem{
				Model: model.ID,
			})
		}
		return &domain.ModelListResp{
			Models: modelsList,
		}, nil
	default:
		return nil, fmt.Errorf("invalid provider: %s", provider)
	}
}

func CheckModel(ctx context.Context, req *domain.CheckModelReq) (*domain.CheckModelResp, error) {
	checkResp := &domain.CheckModelResp{}
	modelType := consts.ModelType(req.Type)
	if modelType == consts.ModelTypeEmbedding || modelType == consts.ModelTypeRerank {
		url := req.BaseURL
		reqBody := map[string]any{}
		if modelType == consts.ModelTypeEmbedding {
			reqBody = map[string]any{
				"model":           req.Model,
				"input":           "PandaWiki is a platform for creating and sharing knowledge bases.",
				"encoding_format": "float",
			}
			url = req.BaseURL + "/embeddings"
		}
		if modelType == consts.ModelTypeRerank {
			reqBody = map[string]any{
				"model": req.Model,
				"documents": []string{
					"PandaWiki is a platform for creating and sharing knowledge bases.",
					"PandaWiki is a platform for creating and sharing knowledge bases.",
					"PandaWiki is a platform for creating and sharing knowledge bases.",
				},
				"query": "PandaWiki",
			}
			url = req.BaseURL + "/rerank"
		}
		body, err := json.Marshal(reqBody)
		if err != nil {
			checkResp.Error = fmt.Sprintf("marshal request body failed: %s", err.Error())
			return checkResp, nil
		}
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			checkResp.Error = fmt.Sprintf("new request failed: %s", err.Error())
			return checkResp, nil
		}
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", req.APIKey))
		request.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			checkResp.Error = fmt.Sprintf("send request failed: %s", err.Error())
			return checkResp, nil
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			checkResp.Error = fmt.Sprintf("request failed: %s", resp.Status)
			return checkResp, nil
		}
		return checkResp, nil
	}

	chatModel, err := GetChatModel(ctx, &domain.ModelMetadata{
		Provider:   consts.ModelProvider(req.Provider),
		ModelName:  req.Model,
		APIKey:     req.APIKey,
		APIHeader:  req.APIHeader,
		BaseURL:    req.BaseURL,
		APIVersion: req.APIVersion,
		ModelType:  consts.ModelType(req.Type),
	})
	if err != nil {
		checkResp.Error = err.Error()
		return checkResp, nil
	}
	resp, err := chatModel.Generate(ctx, []*schema.Message{
		schema.SystemMessage("You are a helpful assistant."),
		schema.UserMessage("hi"),
	})
	if err != nil {
		checkResp.Error = err.Error()
		return checkResp, nil
	}
	content := resp.Content
	if content == "" {
		checkResp.Error = "generate failed"
		return checkResp, nil
	}
	checkResp.Content = content
	return checkResp, nil
}

func GetChatModel(ctx context.Context, model *domain.ModelMetadata) (model.BaseChatModel, error) {
	// config chat model
	modelProvider := consts.ModelProvider(model.Provider)
	var temperature float32 = 0.0
	config := &openai.ChatModelConfig{
		APIKey:      model.APIKey,
		BaseURL:     model.BaseURL,
		Model:       string(model.ModelName),
		Temperature: &temperature,
	}
	if modelProvider == consts.ModelProviderAzureOpenAI {
		config.ByAzure = true
		config.APIVersion = model.APIVersion
		if config.APIVersion == "" {
			config.APIVersion = "2024-10-21"
		}
	}
	if model.APIHeader != "" {
		client := utils.GetHttpClientWithAPIHeaderMap(model.APIHeader)
		if client != nil {
			config.HTTPClient = client
		}
	}
	switch modelProvider {
	case consts.ModelProviderDeepSeek:
		chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
			BaseURL:     model.BaseURL,
			APIKey:      model.APIKey,
			Model:       model.ModelName,
			Temperature: temperature,
		})
		if err != nil {
			return nil, fmt.Errorf("create chat model failed: %w", err)
		}
		return chatModel, nil
	case consts.ModelProviderGemini:
		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey: model.APIKey,
		})
		if err != nil {
			return nil, fmt.Errorf("create genai client failed: %w", err)
		}

		chatModel, err := gemini.NewChatModel(ctx, &gemini.Config{
			Client: client,
			Model:  model.ModelName,
			ThinkingConfig: &genai.ThinkingConfig{
				IncludeThoughts: true,
				ThinkingBudget:  nil,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("create chat model failed: %w", err)
		}
		return chatModel, nil
	case consts.ModelProviderOllama:
		baseUrl, err := utils.URLRemovePath(config.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("ollama url parse failed: %w", err)
		}

		chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
			BaseURL: baseUrl,
			Timeout: config.Timeout,
			Model:   config.Model,
			Options: &api.Options{
				Temperature: temperature,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("create chat model failed: %w", err)
		}
		return chatModel, nil
	default:
		chatModel, err := openai.NewChatModel(ctx, config)
		if err != nil {
			return nil, fmt.Errorf("create chat model failed: %w", err)
		}
		return chatModel, nil
	}
}
