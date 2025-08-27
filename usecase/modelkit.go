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

func ModelList(ctx context.Context, req *domain.ModelListReq) (*domain.ModelListResp, error) {
	log.Printf("ModelList req: provider=%s, baseURL=%s", req.Provider, req.BaseURL)
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
		consts.ModelProviderZhiPu,
		consts.ModelProviderVolcengine:
		return &domain.ModelListResp{
			Models: domain.From(domain.ModelProviders[provider]),
		}, nil
	// 以下模型供应商需要特殊处理
	case consts.ModelProviderOllama:
		resp, err := ollamaListModel(req.BaseURL, httpClient, req.APIHeader)
		// 尝试通过替换baseURL的host为host.docker.internal解决ollama list err
		if err != nil {
			newBaseURL, replaceHostErr := baseURLReplaceHost(req.BaseURL)
			if replaceHostErr != nil {
				return &domain.ModelListResp{
					Error: err.Error(),
				}, nil
			}
			// 替换host后与原始host相同，无需继续请求
			if newBaseURL == req.BaseURL {
				return resp, nil
			}
			_, listErr := ollamaListModel(newBaseURL, httpClient, req.APIHeader)
			// 替换后可以成功请求
			if listErr == nil {
				return &domain.ModelListResp{
					Error: fmt.Errorf("请将host替换为host.docker.internal").Error(),
				}, nil
			}
		}
		// end
		return resp, err
	case consts.ModelProviderGemini:
		client, err := generativeGenai.NewClient(ctx, option.WithAPIKey(req.APIKey))
		if err != nil {
			return &domain.ModelListResp{
				Error: err.Error(),
			}, nil
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
		// openai 兼容模型
	default:
		models, err := reqModelListApi(req, httpClient, &domain.OpenAIResp{})

		if err != nil {
			return &domain.ModelListResp{
				Error: err.Error(),
			}, nil
		}
		return &domain.ModelListResp{
			Models: models,
		}, nil
	}
}

func CheckModel(ctx context.Context, req *domain.CheckModelReq) (*domain.CheckModelResp, error) {
	log.Printf("CheckModel req: provider=%s, model=%s, baseURL=%s", req.Provider, req.Model, req.BaseURL)
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
				log.Printf("Failed to close resp body: %v", closeErr)
			}
		}()
		if resp.StatusCode != http.StatusOK {
			checkResp.Error = resp.Status
			return checkResp, nil
		}
		return checkResp, nil
	}
	// end
	// end
	provider := consts.ParseModelProvider(req.Provider)

	resp, err := getChatModelGenerateChat(ctx, provider, modelType, req.BaseURL, req)
	// 其他模型供应商，尝试修复baseURL
	if err != nil && provider == consts.ModelProviderOther {
		res, err := tryFixBaseURL(ctx, req, provider, modelType)
		if err != nil {
			checkResp.Error = err.Error()
			return checkResp, nil
		}
		if res != "" {
			checkResp.Error = res
			return checkResp, nil
		}
	}
	// end
	if err != nil && provider != consts.ModelProviderOther {
		// 检查错误信息中是否包含余额相关关键词
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
	if err != nil {
		checkResp.Error = err.Error()
		return checkResp, nil
	}

	if resp.Content == "" {
		checkResp.Error = "生成内容失败"
		return checkResp, nil
	}
	checkResp.Content = resp.Content
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

		chatModel, err := gemini.NewChatModel(ctx, &gemini.Config{
			Client: client,
			Model:  model.ModelName,
			ThinkingConfig: &genai.ThinkingConfig{
				IncludeThoughts: true,
				ThinkingBudget:  nil,
			},
		})
		if err != nil {
			return nil, err
		}
		return chatModel, nil
	case consts.ModelProviderOllama:
		baseUrl, err := utils.URLRemovePath(config.BaseURL)
		if err != nil {
			return nil, err
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
			return nil, err
		}
		return chatModel, nil
	// 兼容 openai api
	default:
		chatModel, err := openai.NewChatModel(ctx, config)
		if err != nil {
			return nil, err
		}
		return chatModel, nil
	}
}

// 以下是辅助函数，用于处理模型列表和检查相关的功能

func ollamaListModel(baseURL string, httpClient *http.Client, apiHeader string) (*domain.ModelListResp, error) {
	// get from ollama http://10.10.16.24:11434/api/tags
	u, err := url.Parse(baseURL)
	if err != nil {
		return &domain.ModelListResp{
			Error: err.Error(),
		}, nil
	}
	u.Path = "/api/tags"
	client := request.NewClient(u.Scheme, u.Host, httpClient.Timeout, request.WithClient(httpClient))

	h := request.Header{}
	if apiHeader != "" {
		headers := request.GetHeaderMap(apiHeader)
		maps.Copy(h, headers)
	}
	return request.Get[domain.ModelListResp](client, u.Path, request.WithHeader(h))
}

// 通过修复baseURL尝试修复其它供应商check err, 返回用于提示用户如何修复错误
func tryFixBaseURL(ctx context.Context, req *domain.CheckModelReq, provider consts.ModelProvider, modelType consts.ModelType) (string, error) {
	log.Println("尝试修复")
	// 尝试添加v1
	fixedBaseURL, err := baseURLAddV1(req.BaseURL)
	// baseurl解析错误，直接返回
	if err != nil {
		return "", err
	}

	// baseurl被修改，重新请求
	if fixedBaseURL != req.BaseURL {
		_, err := getChatModelGenerateChat(ctx, provider, modelType, fixedBaseURL, req)
		// 添加v1有效， 提示用户
		if err == nil {
			log.Println("添加v1有效")
			return "请在API地址末尾添加/v1", nil
		}
		log.Println("添加v1无效", err, fixedBaseURL)
	}

	log.Println("尝试替换host")
	// url末尾添加v1无效，尝试替换host为host.docker.internal
	fixedBaseURL, err = baseURLReplaceHost(req.BaseURL)
	// baseurl解析错误，直接返回
	if err != nil {
		return "", err
	}

	if fixedBaseURL != req.BaseURL {
		_, err := getChatModelGenerateChat(ctx, provider, modelType, fixedBaseURL, req)
		// 替换host有效， 提示用户
		if err == nil {
			return "请将host替换为host.docker.internal", nil
		}
		log.Println("替换host无效", err, fixedBaseURL)
	}

	// 替换host也无效，尝试添加v1与替换host
	log.Println("尝试添加v1与替换host")
	fixedBaseURL, err = baseURLAddV1(req.BaseURL)
	// baseurl解析错误，直接返回
	if err != nil {
		return "", err
	}
	fixedBaseURL, err = baseURLReplaceHost(fixedBaseURL)
	// baseurl解析错误，直接返回
	if err != nil {
		return "", err
	}
	// baseurl被修改，重新请求
	if fixedBaseURL != req.BaseURL {
		_, err := getChatModelGenerateChat(ctx, provider, modelType, fixedBaseURL, req)
		// 添加v1与替换host有效， 提示用户
		if err == nil {
			return "API地址末尾添加/v1， host替换为host.docker.internal", nil
		}
		log.Println("添加v1与替换host无效", err, fixedBaseURL)
	}
	return "", nil
}

func getChatModelGenerateChat(ctx context.Context, provider consts.ModelProvider, modelType consts.ModelType, baseURL string, req *domain.CheckModelReq) (*schema.Message, error) {
	chatModel, err := GetChatModel(ctx, &domain.ModelMetadata{
		Provider:   provider,
		ModelName:  req.Model,
		APIKey:     req.APIKey,
		APIHeader:  req.APIHeader,
		BaseURL:    baseURL,
		APIVersion: req.APIVersion,
		ModelType:  modelType,
	})
	if err != nil {
		return nil, err
	}

	return chatModel.Generate(ctx, []*schema.Message{
		schema.SystemMessage("You are a helpful assistant."),
		schema.UserMessage("hi"),
	})
}

// baseURL添加/v1
func baseURLAddV1(inputURL string) (string, error) {
	rawURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	// 没有path, 则添加/v1
	if rawURL.Path == "" {
		rawURL.Path = "/v1"
	}
	return rawURL.String(), nil
}

// baseURL的host换成host.docker.internal
func baseURLReplaceHost(inputURL string) (string, error) {
	rawURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	hostAddress := "host.docker.internal"

	if rawURL.Hostname() != hostAddress {
		if rawURL.Port() != "" {
			rawURL.Host = hostAddress + ":" + rawURL.Port()
		} else {
			rawURL.Host = hostAddress
		}
	}
	return rawURL.String(), nil
}

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
