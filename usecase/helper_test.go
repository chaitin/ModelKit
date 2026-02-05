package usecase

import (
	"testing"

	"github.com/chaitin/ModelKit/v2/consts"
	"github.com/chaitin/ModelKit/v2/domain"
)

func TestBuildOpenAIChatConfig_Azure(t *testing.T) {
	md := &domain.ModelMetadata{
		Provider:  consts.ModelProviderAzureOpenAI,
		ModelName: "gpt-4.1-mini",
		BaseURL:   "https://example.openai.azure.com",
		APIKey:    "test-key",
	}

	cfg := buildOpenAIChatConfig(md)

	if !cfg.ByAzure {
		t.Error("Expected ByAzure to be true")
	}

	if cfg.AzureModelMapperFunc == nil {
		t.Error("Expected AzureModelMapperFunc to be set")
	} else {
		mapped := cfg.AzureModelMapperFunc("gpt-4.1-mini")
		if mapped != "gpt-4.1-mini" {
			t.Errorf("Expected mapped model to be 'gpt-4.1-mini', got '%s'", mapped)
		}
	}
}
