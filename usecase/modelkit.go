package usecase

import (
	"net/http"
	"strings"
	"time"

	"github.com/chaitin/ModelKit/domain"
	"github.com/chaitin/ModelKit/pkg/request"
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

func getHeaderMap(header string) map[string]string {
	headerMap := make(map[string]string)
	for _, h := range strings.Split(header, "\n") {
		if key, value, ok := strings.Cut(h, "="); ok {
			headerMap[key] = value
		}
	}
	return headerMap
}

// func (m *ModelKit) CheckModel(ctx context.Context, req *domain.CheckModelReq) (*domain.Model, error) {
// 	if req.ubType == consts.ModelTypeEmbedding || req.SubType == consts.ModelTypeRerank {
// 		url := domain.ModelProviders[req.Owner].APIBase
// 		reqBody := map[string]any{}
// 		if req.SubType == consts.ModelTypeEmbedding {
// 			reqBody = map[string]any{
// 				"model":           req.ModelID,
// 				"input":           "ModelKit一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
// 				"encoding_format": "float",
// 			}
// 			url = domain.ModelProviders[req.Owner].APIBase + "/embeddings"
// 		}
// 		if req.SubType == consts.ModelTypeRerank {
// 			reqBody = map[string]any{
// 				"model": req.ModelID,
// 				"documents": []string{
// 					"ModelKit一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
// 					"ModelKit一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
// 					"ModelKit一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
// 				},
// 				"query": "ModelKit",
// 			}
// 			url = domain.ModelProviders[req.Owner].APIBase + "/rerank"
// 		}
// 		body, err := json.Marshal(reqBody)
// 		if err != nil {
// 			return nil, err
// 		}
// 		request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
// 		if err != nil {
// 			return nil, err
// 		}
// 		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", req.APIKey))
// 		request.Header.Set("Content-Type", "application/json")
// 		resp, err := http.DefaultClient.Do(request)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer resp.Body.Close()
// 		if resp.StatusCode != http.StatusOK {
// 			return nil, fmt.Errorf("request failed: %s", resp.Status)
// 		}
// 		return &domain.Model{}, nil
// 	}
// 	config := &openai.ChatModelConfig{
// 		APIKey:  req.APIKey,
// 		BaseURL: domain.ModelProviders[req.Owner].APIBase,
// 		Model:   string(req.ModelID),
// 	}
// 	fmt.Println("BaseURL", domain.ModelProviders[req.Owner].APIBase)
// 	// for azure openai
// 	if req.Owner == consts.ModelProviderAzureOpenAI {
// 		config.ByAzure = true
// 		config.APIVersion = domain.ModelProviders[req.Owner].APIVersion
// 		if config.APIVersion == "" {
// 			config.APIVersion = "2024-10-21"
// 		}
// 	}
// 	// end
// 	if domain.ModelProviders[req.Owner].APIHeader != "" {
// 		client := getHttpClientWithAPIHeaderMap(domain.ModelProviders[req.Owner].APIHeader)
// 		if client != nil {
// 			config.HTTPClient = client
// 		}
// 	}
// 	chatModel, err := openai.NewChatModel(ctx, config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp, err := chatModel.Generate(ctx, []*schema.Message{
// 		schema.SystemMessage("You are a helpful assistant."),
// 		schema.UserMessage("hi"),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	content := resp.Content
// 	if content == "" {
// 		return nil, fmt.Errorf("generate failed")
// 	}
// 	return &domain.Model{
// 		ModelType: req.SubType,
// 		OwnedBy:   req.Owner,
// 		ID:        req.ModelID,
// 	}, nil
// }
