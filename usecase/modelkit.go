package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/http"
	"net/url"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	generativeGenai "github.com/google/generative-ai-go/genai"
	"github.com/ollama/ollama/api"
	"google.golang.org/api/option"
	"google.golang.org/genai"

	"github.com/chaitin/ModelKit/consts"
	"github.com/chaitin/ModelKit/domain"
	"github.com/chaitin/ModelKit/pkg/request"
	"github.com/chaitin/ModelKit/utils"
)

// reqModelListApi 获取OpenAI兼容API的模型列表
// 使用泛型和接口抽象来支持不同供应商的响应格式
func reqModelListApi[T domain.ModelResponseParser](req *domain.ModelListReq, httpClient *http.Client, responseType T) ([]domain.ModelListItem, error) {
	u, err := url.Parse(req.BaseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "/models")

	client := request.NewClient(u.Scheme, u.Host, httpClient.Timeout, request.WithClient(httpClient))
	query, err := utils.GetQuery(req)
	if err != nil {
		return nil, err
	}
	resp, err := request.Get[T](
		client, u.Path,
		request.WithHeader(
			request.Header{
				"Authorization": fmt.Sprintf("Bearer %s", req.APIKey),
			},
		),
		request.WithQuery(query),
	)
	if err != nil {
		return nil, err
	}

	return (*resp).ParseModels(), nil
}

func ModelList(ctx context.Context, req *domain.ModelListReq) (*domain.ModelListResp, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			MaxConnsPerHost:     100,
			IdleConnTimeout:     time.Second * 30,
		},
	}
	provider := consts.ParseModelProvider(req.Provider)

	switch provider {
	// 人工返回模型列表
	case consts.ModelProviderAzureOpenAI,
		consts.ModelProviderZhiPu,
		consts.ModelProviderVolcengine:
		return &domain.ModelListResp{
			Models: domain.From(domain.ModelProviders[provider]),
		}, nil
	// 以下模型供应商需要特殊处理
	case consts.ModelProviderOllama:
		// get from ollama http://10.10.16.24:11434/api/tags
		u, err := url.Parse(req.BaseURL)
		if err != nil {
			return nil, err
		}
		u.Path = "/api/tags"
		client := request.NewClient(u.Scheme, u.Host, httpClient.Timeout, request.WithClient(httpClient))

		h := request.Header{}
		if req.APIHeader != "" {
			headers := request.GetHeaderMap(req.APIHeader)
			maps.Copy(h, headers)
		}

		return request.Get[domain.ModelListResp](client, u.Path, request.WithHeader(h))
	case consts.ModelProviderGemini:
		client, err := generativeGenai.NewClient(ctx, option.WithAPIKey(req.APIKey))
		if err != nil {
			return nil, err
		}
		defer func() {
			if closeErr := client.Close(); closeErr != nil {
				log.Printf("Failed to close gemini client: %v", closeErr)
			}
		}()

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
	case consts.ModelProviderGithub:
		models, err := reqModelListApi(req, httpClient, &domain.GithubResp{})
		if err != nil {
			return nil, err
		}
		return &domain.ModelListResp{
			Models: models,
		}, nil
		// openai 兼容模型
	default:
		models, err := reqModelListApi(req, httpClient, &domain.OpenAIResp{})

		if err != nil {
			return nil, err
		}
		return &domain.ModelListResp{
			Models: models,
		}, nil
	}
}

func CheckModel(ctx context.Context, req *domain.CheckModelReq) (*domain.CheckModelResp, error) {
	checkResp := &domain.CheckModelResp{}
	modelType := consts.ParseModelType(req.Type)

	if modelType == consts.ModelTypeEmbedding || modelType == consts.ModelTypeRerank {
		url := req.BaseURL
		reqBody := map[string]any{}
		if modelType == consts.ModelTypeEmbedding {
			reqBody = map[string]any{
				"model":           req.Model,
				"input":           "ModelKit 一个轻量级工具库，提供 AI 模型发现与 API 密钥验证功能，助你快速集成各大模型供应商能力。",
				"encoding_format": "float",
			}
			url = req.BaseURL + "/embeddings"
		}
		if modelType == consts.ModelTypeRerank {
			reqBody = map[string]any{
				"model": req.Model,
				"documents": []string{
					"ModelKit 是一个轻量级工具库，提供 AI 模型发现与 API 密钥验证功能，助你快速集成各大模型供应商能力。",
					"ModelKit 是一个轻量级工具库，提供 AI 模型发现与 API 密钥验证功能，助你快速集成各大模型供应商能力。",
					"ModelKit 是一个轻量级工具库，提供 AI 模型发现与 API 密钥验证功能，助你快速集成各大模型供应商能力。",
				},
				"query": "ModelKit",
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
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				log.Printf("Failed to close resp body: %v", closeErr)
			}
		}()
		if resp.StatusCode != http.StatusOK {
			checkResp.Error = fmt.Sprintf("request failed: %s", resp.Status)
			return checkResp, nil
		}
		return checkResp, nil
	}
	provider := consts.ParseModelProvider(req.Provider)
	chatModel, err := GetChatModel(ctx, &domain.ModelMetadata{
		Provider:   provider,
		ModelName:  req.Model,
		APIKey:     req.APIKey,
		APIHeader:  req.APIHeader,
		BaseURL:    req.BaseURL,
		APIVersion: req.APIVersion,
		ModelType:  modelType,
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
	modelProvider := model.Provider
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
	// 兼容 openai api
	default:
		chatModel, err := openai.NewChatModel(ctx, config)
		if err != nil {
			return nil, fmt.Errorf("create chat model failed: %w", err)
		}
		return chatModel, nil
	}
}
