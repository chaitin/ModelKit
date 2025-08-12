package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/ollama/ollama/api"
	"google.golang.org/genai"

	"github.com/chaitin/ModelKit/consts"
	"github.com/chaitin/ModelKit/domain"
	"github.com/chaitin/ModelKit/pkg/request"
	"github.com/chaitin/ModelKit/utils"
)

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
	switch provider := consts.ModelProvider(req.Provider); provider {
	case consts.ModelProviderAzureOpenAI,
		consts.ModelProviderVolcengine:
		return &domain.ModelListResp{
			Models: domain.From(domain.ModelProviders[provider]),
		}, nil
	case consts.ModelProviderOpenAI,
		consts.ModelProviderHunyuan,
		consts.ModelProviderMoonshot,
		consts.ModelProviderDeepSeek,
		consts.ModelProviderSiliconFlow,
		consts.ModelProviderBaiZhiCloud,
		consts.ModelProviderBaiLian:
		u, err := url.Parse(req.BaseURL)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, "/models")

		client := request.NewClient(u.Scheme, u.Host, httpClient.Timeout, request.WithClient(httpClient))
		query := utils.GetQuery(req)
		resp, err := request.Get[domain.OpenAIResp](
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

		var models []domain.ModelListItem
		for _, item := range resp.Data {
			models = append(models, domain.ModelListItem{Model: item.ID})
		}

		return &domain.ModelListResp{
			Models: models,
		}, nil

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

	default:
		return nil, fmt.Errorf("invalid provider: %s", req.Provider)
	}
}

func CheckModel(ctx context.Context, req *domain.CheckModelReq) (*domain.CheckModelResp, error) {
	checkResp := &domain.CheckModelResp{}
	modelType := consts.ModelType(req.Type)
	modelProvider := consts.ModelProvider(req.Provider)
	if modelType == consts.ModelTypeChat || modelType == consts.ModelTypeRerank {
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
			checkResp.Error = err.Error()
			return checkResp, nil
		}
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			checkResp.Error = err.Error()
			return checkResp, nil
		}
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", req.APIKey))
		request.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			checkResp.Error = err.Error()
			return checkResp, nil
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			checkResp.Error = fmt.Sprintf("request failed: %s", resp.Status)
			return checkResp, nil
		}
		return checkResp, nil
	}
	config := &openai.ChatModelConfig{
		APIKey:  req.APIKey,
		BaseURL: req.BaseURL,
		Model:   req.Model,
	}
	// for azure openai
	if modelProvider == consts.ModelProviderAzureOpenAI {
		config.ByAzure = true
		config.APIVersion = req.APIVersion
		if config.APIVersion == "" {
			config.APIVersion = "2024-10-21"
		}
	}
	// 阿里云百炼模型支持流式和思考功能
	if modelProvider == consts.ModelProviderBaiLian {
		config.ExtraFields = map[string]any{
			"stream":          true,
			"enable_thinking": true,
		}
	}
	if req.APIHeader != "" {
		client := utils.GetHttpClientWithAPIHeaderMap(req.APIHeader)
		if client != nil {
			config.HTTPClient = client
		}
	}
	chatModel, err := openai.NewChatModel(ctx, config)
	if err != nil {
		checkResp.Error = err.Error()
		return checkResp, nil
	}
	// 阿里云百炼模型(不支持 翻译/OCR 模型的添加)
	if modelProvider == consts.ModelProviderBaiLian {
		msgs := []*schema.Message{
			schema.SystemMessage("You are a helpful assistant."),
			schema.UserMessage("hi"),
		}
		stream, err := chatModel.Stream(ctx, msgs)
		if err != nil {
			checkResp.Error = err.Error()
			return checkResp, nil
		}
		var content string
		for {
			msg, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				checkResp.Error = err.Error()
				return checkResp, nil
			}
			if msg.Content != "" {
				content += msg.Content
			}
		}

		if content == "" {
			checkResp.Error = "generate failed"
			return checkResp, nil
		}
	} else {
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
	}

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
