package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	generativeGenai "github.com/google/generative-ai-go/genai"
	"github.com/ollama/ollama/api"
	"google.golang.org/api/option"
	"google.golang.org/genai"

	"github.com/chaitin/ModelKit/v2/consts"
	"github.com/chaitin/ModelKit/v2/domain"
	"github.com/chaitin/ModelKit/v2/utils"
)

type ModelKit struct {
	logger *slog.Logger
}

// NewModelKit 创建一个新的ModelKit实例
func NewModelKit(logger *slog.Logger) *ModelKit {
	return &ModelKit{
		logger: logger,
	}
}

func (m *ModelKit) ModelList(ctx context.Context, req *domain.ModelListReq) (*domain.ModelListResp, error) {
	if m.logger != nil {
		m.logger.Info("ModelList req:", req.Provider, req.BaseURL)
	} else {
		log.Printf("ModelList req: provider=%s, baseURL=%s", req.Provider, req.BaseURL)
	}
	httpClient := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			MaxConnsPerHost:     100,
			IdleConnTimeout:     time.Second * 30,
			Proxy:               http.ProxyFromEnvironment,
		},
	}
	provider := consts.ParseModelProvider(req.Provider)

	switch provider {
	// 人工返回模型列表
	case consts.ModelProviderAzureOpenAI,
		consts.ModelProviderVolcengine:
		return &domain.ModelListResp{
			Models: domain.From(domain.ModelProviders[provider]),
		}, nil
	case consts.ModelProviderGemini:
		client, err := generativeGenai.NewClient(ctx, option.WithAPIKey(req.APIKey))
		if err != nil {
			return &domain.ModelListResp{
				Error: err.Error(),
			}, nil
		}
		defer func() {
			if closeErr := client.Close(); closeErr != nil {
				if m.logger != nil {
					m.logger.Error("Failed to close gemini client: %v", slog.Any("error", closeErr))
				} else {
					log.Printf("Failed to close gemini client: %v", closeErr)
				}
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
			return &domain.ModelListResp{
				Error: fmt.Errorf("获取Gemini模型列表失败: 未找到可用模型").Error(),
			}, nil
		}

		return &domain.ModelListResp{
			Models: modelsList,
		}, nil
	case consts.ModelProviderGithub:
		models, err := reqModelListApi(req, httpClient, &domain.GithubResp{})
		if err != nil {
			return &domain.ModelListResp{
				Error: err.Error(),
			}, nil
		}
		return &domain.ModelListResp{
			Models: models,
		}, nil
	case consts.ModelProviderOllama:
		var modelListResp domain.ModelListResp
		var err error
		// 当BaseURL以/v1结尾时，使用OpenAI兼容的API调用方式
		if strings.HasSuffix(req.BaseURL, "/v1") {
			var models []domain.ModelListItem
			models, err = reqModelListApi(req, httpClient, &domain.OpenAIResp{})
			if err == nil {
				modelListResp.Models = models
			}
		} else {
			var resp *domain.ModelListResp
			resp, err = ollamaListModel(req.BaseURL, httpClient, req.APIHeader)
			if err == nil {
				modelListResp = *resp
			}
		}
		// ollama list发生错误， 尝试修复url
		if err != nil {
			msg := generateBaseURLFixSuggestion(err.Error(), req.BaseURL, provider)
			if msg == "" {
				return &domain.ModelListResp{
					Error: err.Error(),
				}, nil
			} else {
				return &domain.ModelListResp{
					Error: msg,
				}, nil
			}
		}
		// end
		return &modelListResp, err
	case consts.ModelProviderGPUStack:
		models, err := reqModelListApi(req, httpClient, &domain.GPUStackListModelResp{})
		// gpu stack list发生错误， 尝试修复url
		if err != nil {
			m.logger.Error("GPUStack list model failed", "error", err, "models: ", models)
			msg := generateBaseURLFixSuggestion(err.Error(), req.BaseURL, provider)
			if msg == "" {
				return &domain.ModelListResp{
					Error: err.Error(),
				}, nil
			} else {
				return &domain.ModelListResp{
					Error: msg,
				}, nil
			}
		}
		// end
		return &domain.ModelListResp{
			Models: models,
		}, nil
		// openai 兼容模型
	default:
		models, err := reqModelListApi(req, httpClient, &domain.OpenAIResp{})

		// ollama可能有url输入格式问题
		if err != nil && provider == consts.ModelProviderOllama {
			msg := generateBaseURLFixSuggestion(err.Error(), req.BaseURL, provider)
			if msg == "" {
				return &domain.ModelListResp{
					Error: err.Error(),
				}, nil
			} else {
				return &domain.ModelListResp{
					Error: msg,
				}, nil
			}
		} else if err != nil {
			return &domain.ModelListResp{
				Error: err.Error(),
			}, nil
		}
		return &domain.ModelListResp{
			Models: models,
		}, nil
	}
}

func (m *ModelKit) CheckModel(ctx context.Context, req *domain.CheckModelReq) (*domain.CheckModelResp, error) {
	if m.logger != nil {
		m.logger.Info("CheckModel req", "provider", req.Provider, "model", req.Model, "baseURL", req.BaseURL)
	} else {
		log.Printf("CheckModel req: provider=%s, model=%s, baseURL=%s", req.Provider, req.Model, req.BaseURL)
	}
	checkResp := &domain.CheckModelResp{}
	modelType := consts.ParseModelType(req.Type)

	// embedding 与 rerank 模型检查
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
		client := http.DefaultClient
		client.Transport = http.DefaultTransport
		resp, err := client.Do(request)
		if err != nil {
			checkResp.Error = err.Error()
			return checkResp, nil
		}
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				if m.logger != nil {
					m.logger.Error("Failed to close resp body: %v", slog.Any("error", closeErr))
				} else {
					log.Printf("Failed to close resp body: %v", closeErr)
				}
			}
		}()
		if resp.StatusCode != http.StatusOK {
			checkResp.Error = resp.Status
			return checkResp, nil
		}
		return checkResp, nil
	}
	// end
	provider := consts.ParseModelProvider(req.Provider)

	resp, err := m.getChatModelGenerateChat(ctx, provider, modelType, req.BaseURL, req)

	// 可编辑url的供应商，尝试修复baseURL
	if err != nil && (provider == consts.ModelProviderOther || provider == consts.ModelProviderOllama || provider == consts.ModelProviderAzureOpenAI) {
		msg := generateBaseURLFixSuggestion(err.Error(), req.BaseURL, provider)
		if msg == "" {
			checkResp.Error = err.Error()
		} else {
			checkResp.Error = msg
		}
		return checkResp, nil
	}
	// end
	// 检查错误信息中是否包含余额相关关键词
	if err != nil && provider != consts.ModelProviderOther {
		errorMsg := strings.ToLower(err.Error())
		for _, keyword := range consts.ApiKeyBalanceKeyWords {
			if strings.Contains(errorMsg, keyword) {
				checkResp.Error = "API Key余额不足"
				return checkResp, nil
			}
		}
		checkResp.Error = err.Error()
		return checkResp, nil
	}
	// end
	if err != nil {
		checkResp.Error = err.Error()
		return checkResp, nil
	}

	if resp == "" {
		checkResp.Error = "生成内容失败"
		return checkResp, nil
	}
	checkResp.Content = resp
	return checkResp, nil
}

func (m *ModelKit) GetChatModel(ctx context.Context, model *domain.ModelMetadata) (model.BaseChatModel, error) {
	// config chat model
	modelProvider := model.Provider

	// 使用高级参数中的温度值，如果没有设置则使用默认值0.0
	var temperature float32 = 0.0
	if model.Temperature != nil {
		temperature = *model.Temperature
	}

	config := &openai.ChatModelConfig{
		APIKey:      model.APIKey,
		BaseURL:     model.BaseURL,
		Model:       string(model.ModelName),
		Temperature: &temperature,
	}

	// 添加高级参数支持
	if model.MaxTokens != nil {
		config.MaxTokens = model.MaxTokens
	}
	if model.TopP != nil {
		config.TopP = model.TopP
	}
	if len(model.Stop) > 0 {
		config.Stop = model.Stop
	}
	if model.PresencePenalty != nil {
		config.PresencePenalty = model.PresencePenalty
	}
	if model.FrequencyPenalty != nil {
		config.FrequencyPenalty = model.FrequencyPenalty
	}
	if model.ResponseFormat != nil {
		config.ResponseFormat = model.ResponseFormat
	}
	if model.Seed != nil {
		config.Seed = model.Seed
	}
	if model.LogitBias != nil {
		config.LogitBias = model.LogitBias
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
		deepseekConfig := &deepseek.ChatModelConfig{
			BaseURL:     model.BaseURL,
			APIKey:      model.APIKey,
			Model:       model.ModelName,
			Temperature: temperature,
		}

		// 添加 DeepSeek 支持的高级参数
		if model.MaxTokens != nil {
			deepseekConfig.MaxTokens = *model.MaxTokens
		}
		if model.TopP != nil {
			deepseekConfig.TopP = *model.TopP
		}
		if len(model.Stop) > 0 {
			deepseekConfig.Stop = model.Stop
		}
		if model.PresencePenalty != nil {
			deepseekConfig.PresencePenalty = *model.PresencePenalty
		}
		if model.FrequencyPenalty != nil {
			deepseekConfig.FrequencyPenalty = *model.FrequencyPenalty
		}
		// ResponseFormat, Seed, LogitBias 在 DeepSeek 配置中不支持，跳过

		chatModel, err := deepseek.NewChatModel(ctx, deepseekConfig)
		if err != nil {
			return nil, err
		}
		return chatModel, nil
	case consts.ModelProviderGemini:
		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey: model.APIKey,
		})
		if err != nil {
			return nil, err
		}

		geminiConfig := &gemini.Config{
			Client: client,
			Model:  model.ModelName,
			ThinkingConfig: &genai.ThinkingConfig{
				IncludeThoughts: true,
				ThinkingBudget:  nil,
			},
		}

		// 添加 Gemini 支持的高级参数
		if model.MaxTokens != nil {
			geminiConfig.MaxTokens = model.MaxTokens
		}
		if model.Temperature != nil {
			geminiConfig.Temperature = model.Temperature
		}
		if model.TopP != nil {
			geminiConfig.TopP = model.TopP
		}
		chatModel, err := gemini.NewChatModel(ctx, geminiConfig)
		if err != nil {
			return nil, err
		}
		return chatModel, nil
	case consts.ModelProviderOllama:
		// 当BaseURL以/v1结尾时，使用OpenAI兼容的API调用方式
		if strings.HasSuffix(model.BaseURL, "/v1") {
			chatModel, err := openai.NewChatModel(ctx, config)
			if err != nil {
				return nil, err
			}
			return chatModel, nil
		} else {
			baseUrl, err := utils.URLRemovePath(config.BaseURL)
			if err != nil {
				return nil, err
			}

			ollamaOptions := &api.Options{
				Temperature: temperature,
			}

			// 添加 Ollama 支持的高级参数
			if model.TopP != nil {
				ollamaOptions.TopP = *model.TopP
			}
			if len(model.Stop) > 0 {
				ollamaOptions.Stop = model.Stop
			}
			if model.PresencePenalty != nil {
				ollamaOptions.PresencePenalty = *model.PresencePenalty
			}
			if model.FrequencyPenalty != nil {
				ollamaOptions.FrequencyPenalty = *model.FrequencyPenalty
			}
			if model.Seed != nil {
				ollamaOptions.Seed = *model.Seed
			}

			chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
				BaseURL: baseUrl,
				Timeout: config.Timeout,
				Model:   config.Model,
				Options: ollamaOptions,
			})
			if err != nil {
				return nil, err
			}
			return chatModel, nil
		}

		// 兼容 openai api
	default:
		chatModel, err := openai.NewChatModel(ctx, config)
		if err != nil {
			return nil, err
		}
		return chatModel, nil
	}
}