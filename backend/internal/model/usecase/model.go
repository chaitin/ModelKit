package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"

	"github.com/chaitin/ModelKit/backend/config"
	"github.com/chaitin/ModelKit/backend/consts"
	"github.com/chaitin/ModelKit/backend/db"
	"github.com/chaitin/ModelKit/backend/domain"
	"github.com/chaitin/ModelKit/backend/pkg/cvt"
	"github.com/chaitin/ModelKit/backend/pkg/request"
)

type ModelUsecase struct {
	logger *slog.Logger
	repo   domain.ModelRepo
	cfg    *config.Config
	client *http.Client
}

func NewModelUsecase(
	logger *slog.Logger,
	repo domain.ModelRepo,
	cfg *config.Config,
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
	return &ModelUsecase{repo: repo, cfg: cfg, logger: logger, client: client}
}

func (m *ModelUsecase) CheckModel(ctx context.Context, req *domain.CheckModelReq) (*domain.Model, error) {
	if req.Type == consts.ModelTypeEmbedding || req.Type == consts.ModelTypeReranker {
		url := req.APIBase
		reqBody := map[string]any{}
		if req.Type == consts.ModelTypeEmbedding {
			reqBody = map[string]any{
				"model":           req.ModelName,
				"input":           "ModelKit一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
				"encoding_format": "float",
			}
			url = req.APIBase + "/embeddings"
		}
		if req.Type == consts.ModelTypeReranker {
			reqBody = map[string]any{
				"model": req.ModelName,
				"documents": []string{
					"ModelKit一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
					"ModelKit一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
					"ModelKit一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
				},
				"query": "ModelKit",
			}
			url = req.APIBase + "/rerank"
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
		BaseURL: req.APIBase,
		Model:   string(req.ModelName),
	}
	// for azure openai
	if req.Provider == consts.ModelProviderAzureOpenAI {
		config.ByAzure = true
		config.APIVersion = req.APIVersion
		if config.APIVersion == "" {
			config.APIVersion = "2024-10-21"
		}
	}
	if req.APIHeader != "" {
		client := getHttpClientWithAPIHeaderMap(req.APIHeader)
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
		ModelType: req.Type,
		Provider:  req.Provider,
		ModelName: req.ModelName,
		APIBase:   req.APIBase,
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

// Update implements domain.ModelUsecase.
func (m *ModelUsecase) UpdateModel(ctx context.Context, req *domain.UpdateModelReq) (*domain.Model, error) {
	model, err := m.repo.UpdateModel(ctx, req.ID, func(tx *db.Tx, old *db.Model, up *db.ModelUpdateOne) error {
		if req.APIKey != nil {
			up.SetAPIKey(*req.APIKey)
		}
		if req.APIVersion != nil {
			up.SetAPIVersion(*req.APIVersion)
		}
		if req.APIHeader != nil {
			up.SetAPIHeader(*req.APIHeader)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return cvt.From(model, &domain.Model{}), nil
}

// GetModel implements domain.ModelUsecase.
func (m *ModelUsecase) GetModel(ctx context.Context, req *domain.GetModelReq) (*domain.Model, error) {
	model, err := m.repo.GetModel(ctx, req.ModelName, req.Provider)
	if err != nil {
		return nil, err
	}
	return cvt.From(model, &domain.Model{}), nil
}

// ListModel implements domain.ModelUsecase.
func (m *ModelUsecase) ListModel(ctx context.Context, req *domain.ListModelReq) ([]*domain.Model, error) {
	models, err := m.repo.ListModel(ctx, req)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Model, len(models))
	for i, model := range models {
		result[i] = cvt.From(model, &domain.Model{})
	}
	return result, nil
}
