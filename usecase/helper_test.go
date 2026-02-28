package usecase

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestCheckModel_TemperaturePassed(t *testing.T) {
	testName := "TestCheckModel_TemperaturePassed_Provider=Moonshot_Model=kimi-k2.5_Temp=1"
	// 1. Setup a test server to intercept the request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read and parse request body
		var reqBody struct {
			Temperature float32 `json:"temperature"`
			Model       string  `json:"model"`
			Messages    []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
			return
		}

		// 2. Assert that Temperature is correctly passed
		expectedTemp := float32(1.0)
		if reqBody.Temperature != expectedTemp {
			t.Errorf("fail case: %s; Expected temperature %f, got %f", testName, expectedTemp, reqBody.Temperature)
		}

		// Return a valid dummy response to avoid error in CheckModel before it returns
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "chatcmpl-123",
			"object": "chat.completion",
			"created": 1677652288,
			"model": "gpt-3.5-turbo-0613",
			"choices": [{
				"index": 0,
				"message": {
					"role": "assistant",
					"content": "Hello there!"
				},
				"finish_reason": "stop"
			}],
			"usage": {
				"prompt_tokens": 9,
				"completion_tokens": 12,
				"total_tokens": 21
			}
		}`))
	}))
	defer ts.Close()

	// 3. Setup ModelKit and call CheckModel
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	mk := NewModelKit(logger)
	ctx := context.Background()

	req := &domain.CheckModelReq{
		Provider: "Moonshot", // Moonshot uses OpenAI client
		Model:    "kimi-k2.5",
		BaseURL:  ts.URL, // Point to our test server
		APIKey:   "sk-test",
		Type:     "llm",
		Param: &domain.ModelParam{
			Temperature: func() *float32 { t := float32(1.0); return &t }(),
			MaxTokens:   100,
		},
	}

	resp, err := mk.CheckModel(ctx, req)
	if err != nil {
		t.Errorf("fail case: %s; CheckModel failed: %v", testName, err)
		return
	}
	if resp != nil && resp.Error != "" {
		t.Errorf("fail case: %s; CheckModel returned error: %s", testName, resp.Error)
		return
	}
	t.Logf("pass case: %s; response: %+v", testName, resp)
}
