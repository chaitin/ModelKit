package usecase

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	arkEmb "github.com/cloudwego/eino-ext/components/embedding/ark"
	ollamaEmb "github.com/cloudwego/eino-ext/components/embedding/ollama"
	openaiEmb "github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/model"

	bailianEmb "github.com/chaitin/ModelKit/v2/components/embedder/bailian"
	baaiReranker "github.com/chaitin/ModelKit/v2/components/reranker/baai"
	bailianReranker "github.com/chaitin/ModelKit/v2/components/reranker/bailian"
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
	case consts.ModelProviderAzureOpenAI, consts.ModelProviderVolcengine:
		return m.listStaticProvider(req, provider)
	case consts.ModelProviderGemini:
		return m.listGemini(ctx, req)
	case consts.ModelProviderGithub:
		return m.listGithub(req, httpClient)
	case consts.ModelProviderOllama:
		return m.listOllama(req, httpClient)
	case consts.ModelProviderGPUStack:
		return m.listGPUStack(req, httpClient)
	default:
		return m.listOpenAI(req, httpClient, provider)
	}
}

func (m *ModelKit) CheckModel(ctx context.Context, req *domain.CheckModelReq) (*domain.CheckModelResp, error) {
	if m.logger != nil {
		m.logger.Info("CheckModel req", "provider", req.Provider, "model", req.Model, "baseURL", req.BaseURL)
	} else {
		log.Printf("CheckModel req: provider=%s, model=%s, baseURL=%s", req.Provider, req.Model, req.BaseURL)
	}
	modelType := consts.ParseModelType(req.Type)
	switch modelType {
	case consts.ModelTypeEmbedding:
		return m.checkEmbeddingModel(ctx, req)
	case consts.ModelTypeRerank:
		return m.checkRerankModel(ctx, req)
	default:
		return m.checkChatModel(ctx, req)
	}
}

func (m *ModelKit) GetChatModel(ctx context.Context, md *domain.ModelMetadata) (model.BaseChatModel, error) {
	switch md.Provider {
	case consts.ModelProviderDeepSeek:
		return newDeepseekChatModel(ctx, md)
	case consts.ModelProviderGemini:
		return newGeminiChatModel(ctx, md)
	case consts.ModelProviderOllama:
		return newOllamaChatModel(ctx, md)
	default:
		cfg := buildOpenAIChatConfig(md)
		return openai.NewChatModel(ctx, cfg)
	}
}

func (m *ModelKit) GetEmbedder(ctx context.Context, model *domain.ModelMetadata) (embedding.Embedder, error) {
	// dimensions := consts.DefaultDimensions
	cfg := &openaiEmb.EmbeddingConfig{
		APIKey:     model.APIKey,
		Model:      model.ModelName,
		BaseURL:    model.BaseURL,
		Dimensions: model.EmbedderParam.Dimension,
		// Dimensions: &dimensions,
	}

	switch model.Provider {
	case consts.ModelProviderBaiLian:
		return bailianEmb.NewEmbedder(ctx, &bailianEmb.EmbeddingConfig{
			APIKey:         model.APIKey,
			Model:          model.ModelName,
			BaseURL:        model.BaseURL,
			Dimension:      model.EmbedderParam.Dimension,
			TextType:       model.EmbedderParam.TextType,
			OutputType:     model.EmbedderParam.OutputType,
			EncodingFormat: model.EmbedderParam.EncodingFormat,
			Instruct:       model.EmbedderParam.Instruct,
		})
	case consts.ModelProviderAzureOpenAI:
		cfg.ByAzure = true
		cfg.APIVersion = model.APIVersion
		if cfg.APIVersion == "" {
			cfg.APIVersion = "2024-10-21"
		}
		return openaiEmb.NewEmbedder(ctx, cfg)
	case consts.ModelProviderOllama:
		if strings.HasSuffix(model.BaseURL, "/v1") {
			return openaiEmb.NewEmbedder(ctx, cfg)
		}
		baseUrl, err := utils.URLRemovePath(model.BaseURL)
		if err != nil {
			return nil, err
		}
		return ollamaEmb.NewEmbedder(ctx, &ollamaEmb.EmbeddingConfig{
			BaseURL: baseUrl,
			Model:   model.ModelName,
		})
	case consts.ModelProviderVolcengine:
		return arkEmb.NewEmbedder(ctx, &arkEmb.EmbeddingConfig{
			APIKey:  model.APIKey,
			Model:   model.ModelName,
			BaseURL: model.BaseURL,
		})
	case consts.ModelProviderGemini:
		return nil, fmt.Errorf("该提供商暂不支持向量模型")
	default:
		return openaiEmb.NewEmbedder(ctx, cfg)
	}
}

func (m *ModelKit) GetReranker(ctx context.Context, model *domain.ModelMetadata) (domain.Reranker, error) {
	if model.BaseURL == "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank#" {
		model.Provider = consts.ModelProviderBaiLian
	}

	switch model.Provider {
	case consts.ModelProviderBaiLian:
		return bailianReranker.NewReranker(ctx, bailianReranker.RerankerConfig{
			Model:   model.ModelName,
			BaseUrl: model.BaseURL,
			APIKey:  model.APIKey,
		}), nil
	default:
		return baaiReranker.NewReranker(ctx, baaiReranker.RerankerConfig{
			Model:   model.ModelName,
			BaseUrl: model.BaseURL,
			APIKey:  model.APIKey,
		}), nil
	}
}

func (m *ModelKit) UseEmbedder(ctx context.Context, e embedding.Embedder, texts []string) (*domain.EmbeddingsResponse, error) {

	if de, ok := e.(interface {
		EmbedStringsExt(context.Context, []string, ...embedding.Option) (*domain.EmbeddingsResponse, error)
	}); ok {
		return de.EmbedStringsExt(ctx, texts)
	}

	dense, err := e.EmbedStrings(ctx, texts)
	if err != nil {
		return nil, err
	}

	out := &domain.EmbeddingsResponse{
		Embeddings: make([]domain.EmbeddingItem, 0, len(dense)),
		Usage:      domain.EmbeddingUsage{TotalTokens: 0},
	}
	for i := range dense {
		out.Embeddings = append(out.Embeddings, domain.EmbeddingItem{
			Embedding: dense[i],
			TextIndex: i,
		})
	}
	return out, nil
}
