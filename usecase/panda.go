package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"slices"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"github.com/chaitin/ModelKit/consts"
	"github.com/chaitin/ModelKit/domain"
)

func PandaModelList(ctx context.Context, provider string, apiKey string, baseURL string, apiHeader string, modelType string, modelSubType string) (*domain.GetProviderModelListResp, error) {
	switch provider := consts.ModelOwner(provider); provider {
	case consts.ModelOwnerMoonshot,
		consts.ModelOwnerDeepSeek,
		consts.ModelOwnerAzureOpenAI,
		consts.ModelOwnerVolcengine,
		consts.ModelOwnerZhiPu:
		return &domain.GetProviderModelListResp{
			Models: domain.From(domain.ModelOwners[provider]),
		}, nil
	case consts.ModelOwnerGemini:
		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
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
		u, err := url.Parse(baseURL)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, "/models")
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
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
		u, err := url.Parse(baseURL)
		if err != nil {
			return nil, err
		}
		u.Path = "/api/tags"
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}
		if apiHeader != "" {
			headers := getHeaderMap(apiHeader)
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
		modelType := consts.ModelType(modelType)
		if modelType == consts.ModelTypeEmbedding || modelType == consts.ModelTypeReranker {
			if provider == consts.ModelOwnerBaiZhiCloud {
				if modelType == consts.ModelTypeEmbedding {
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
		u, err := url.Parse(baseURL)
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
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
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
		return nil, fmt.Errorf("invalid provider: %s", provider)
	}
}