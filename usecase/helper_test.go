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

func TestBuildOpenAIChatConfig_IgnoreTemperatureForReasoningModels(t *testing.T) {
	temp := float32(1.0)
	testCases := []struct {
		name      string
		modelName string
		ignored   bool
	}{
		{name: "o1", modelName: "o1-preview", ignored: true},
		{name: "o3", modelName: "o3-mini", ignored: true},
		{name: "o4", modelName: "o4-mini", ignored: true},
		{name: "gpt-5", modelName: "gpt-5-mini", ignored: true},
		{name: "regular model", modelName: "gpt-4.1-mini", ignored: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			md := &domain.ModelMetadata{
				ModelName:   tc.modelName,
				Temperature: &temp,
			}

			cfg := buildOpenAIChatConfig(md)

			if tc.ignored {
				if cfg.Temperature != nil {
					t.Fatalf("expected temperature to be ignored for model %q, got %v", tc.modelName, *cfg.Temperature)
				}
				return
			}

			if cfg.Temperature == nil {
				t.Fatalf("expected temperature to be preserved for model %q", tc.modelName)
			}
			if *cfg.Temperature != temp {
				t.Fatalf("expected temperature %v for model %q, got %v", temp, tc.modelName, *cfg.Temperature)
			}
		})
	}
}

// 测试视觉模型正则识别逻辑
// 1. 测试官方已确认支持图片输入的新模型能识别为视觉模型
// 2. 测试未确认支持图片输入或明确为文本/代码模型的模型不会误判为视觉模型
func TestIsVisionModel_NewModelRegex(t *testing.T) {
	provider := string(consts.ModelProviderOther)
	testCases := []struct {
		name      string
		modelName string
		want      bool
	}{
		{name: "gpt-5.5", modelName: "gpt-5.5", want: true},
		{name: "gpt-5.4", modelName: "gpt-5.4", want: true},
		{name: "gpt-5.4-mini", modelName: "gpt-5.4-mini", want: true},
		{name: "minimax-m2.7", modelName: "minimax-m2.7", want: false},
		{name: "minimax-m2.5", modelName: "minimax-m2.5", want: false},
		{name: "glm-5.1", modelName: "glm-5.1", want: false},
		{name: "glm-5", modelName: "glm-5", want: false},
		{name: "glm-4.7", modelName: "glm-4.7", want: false},
		{name: "kimi-k2.5", modelName: "kimi-k2.5", want: true},
		{name: "kimi-k2.6", modelName: "kimi-k2.6", want: true},
		{name: "qwen3.6-plus", modelName: "qwen3.6-plus", want: true},
		{name: "qwen3.5-plus", modelName: "qwen3.5-plus", want: true},
		{name: "qwen3.5-flash", modelName: "qwen3.5-flash", want: true},
		{name: "qwen-flash", modelName: "qwen-flash", want: true},
		{name: "deepseek-v4-pro", modelName: "deepseek-v4-pro", want: false},
		{name: "deepseek-v4-flash", modelName: "deepseek-v4-flash", want: false},
		{name: "mimo-v2.5-pro", modelName: "mimo-v2.5-pro", want: false},
		{name: "mimo-v2.5", modelName: "mimo-v2.5", want: true},
		{name: "mimo-v2-flash", modelName: "mimo-v2-flash", want: false},
		{name: "glm-4.5-air", modelName: "glm-4.5-air", want: false},
		{name: "qwen3.7-max", modelName: "qwen3.7-max", want: true},
		{name: "qwen3.6-max", modelName: "qwen3.6-max", want: true},
		{name: "doubao-seed-2.0-code", modelName: "doubao-seed-2.0-code", want: true},
		{name: "doubao-seed-2.0-pro", modelName: "doubao-seed-2.0-pro", want: true},
		{name: "hy3", modelName: "hy3", want: false},
		{name: "deepseek-v3.2", modelName: "deepseek-v3.2", want: false},
		{name: "minimax-m3", modelName: "minimax-m3", want: true},
		{name: "qwen3.7-plus", modelName: "qwen3.7-plus", want: true},
		{name: "qwen3.6-flash", modelName: "qwen3.6-flash", want: true},
		{name: "kimi-k2.7-code", modelName: "kimi-k2.7-code", want: true},
		{name: "glm-5.2", modelName: "glm-5.2", want: false},
	}

	for _, tc := range testCases {
		tc := tc
		// 测试单个模型是否按官方视觉能力判断
		t.Run(tc.name, func(t *testing.T) {
			if got := isVisionModel(tc.modelName, provider); got != tc.want {
				t.Fatalf("isVisionModel(%q) = %v, want %v", tc.modelName, got, tc.want)
			}
		})
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
