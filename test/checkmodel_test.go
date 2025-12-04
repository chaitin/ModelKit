package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/chaitin/ModelKit/v2/consts"
	"github.com/chaitin/ModelKit/v2/domain"
	"github.com/chaitin/ModelKit/v2/usecase"
)

func formatParam(p *domain.ModelParam) string {
	if p == nil {
		return "nil"
	}
	parts := make([]string, 0, 8)
	if p.ContextWindow != 0 {
		parts = append(parts, fmt.Sprintf("context_window=%d", p.ContextWindow))
	}
	if p.MaxTokens != 0 {
		parts = append(parts, fmt.Sprintf("max_tokens=%d", p.MaxTokens))
	}
	if p.R1Enabled {
		parts = append(parts, "r1_enabled=true")
	}
	if p.SupportComputerUse {
		parts = append(parts, "support_computer_use=true")
	}
	if p.SupportImages {
		parts = append(parts, "support_images=true")
	}
	if p.SupportPromptCache {
		parts = append(parts, "support_prompt_cache=true")
	}
	if p.Temperature != nil {
		parts = append(parts, fmt.Sprintf("temperature=%.2f", *p.Temperature))
	}
	if len(parts) == 0 {
		return "empty"
	}
	return strings.Join(parts, ",")
}

func TestCheckModelCombinations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	mk := usecase.NewModelKit(nil)

	apiKeyBZ := strings.TrimSpace(os.Getenv("baizhiapikey"))
	apiKeyBZ1 := strings.TrimSpace(os.Getenv("baizhiapikey1"))
	apiKeyBZ2 := strings.TrimSpace(os.Getenv("baizhiapikey2"))
	keys := make([]string, 0, 3)
	if apiKeyBZ != "" {
		keys = append(keys, apiKeyBZ)
	}
	if apiKeyBZ1 != "" {
		keys = append(keys, apiKeyBZ1)
	}
	if apiKeyBZ2 != "" {
		keys = append(keys, apiKeyBZ2)
	}
	if len(keys) == 0 {
		t.Fatalf("missing baizhiapikeys")
	}

	listBaseURL := "https://model-square.app.baizhi.cloud/v1"
	ml, err := mk.ModelList(ctx, &domain.ModelListReq{
		Provider: string(consts.ModelProviderBaiZhiCloud),
		BaseURL:  listBaseURL,
		APIKey:   keys[0],
		Type:     "",
	})
	if err != nil {
		t.Fatalf("list model error: %v", err)
	}
	if ml == nil || len(ml.Models) == 0 {
		if ml != nil && ml.Error != "" {
			t.Skipf("list model skipped: %s", ml.Error)
			return
		}
	}

	baseURLs := []string{
		"https://model-square.app.baizhi.cloud/v1",
		"https://model-square.app.baizhi.cloud/v1#",
	}
	types := []string{"chat", "embedding", "rerank"}
	temp := float32(0.2)
	params := []*domain.ModelParam{
		nil,
		&domain.ModelParam{},
	}
	params = append(params, &domain.ModelParam{Temperature: &temp})
	params = append(params, &domain.ModelParam{SupportImages: true})

	keyPool := make(chan string, len(keys))
	for _, k := range keys {
		keyPool <- k
	}

	for _, bu := range baseURLs {
		for _, tp := range types {
			tpLocal := tp
			buLocal := bu
			mlType, err := mk.ModelList(ctx, &domain.ModelListReq{
				Provider: string(consts.ModelProviderBaiZhiCloud),
				BaseURL:  listBaseURL,
				APIKey:   keys[0],
				Type:     tpLocal,
			})
			if err != nil {
				t.Errorf("list model error: %v", err)
				continue
			}
			if mlType == nil || len(mlType.Models) == 0 {
				if mlType != nil && mlType.Error != "" {
					t.Errorf("list model error: %s", mlType.Error)
					continue
				}
				t.Errorf("no models returned for type: %s", tpLocal)
				continue
			}
			models := mlType.Models
			if len(models) > 3 {
				models = models[:3]
			}
			if tpLocal == "chat" {
				models = []domain.ModelListItem{
					{Model: "qwen-vl-max-latest"},
					{Model: "glm-4.5"},
				}
			}
			for _, m := range models {
				for idx, pm := range params {
					mLocal := m
					pmLocal := pm
					idxLocal := idx
					testName := fmt.Sprintf(
						"provider=%s base=%s model=%s type=%s apiKey=present param={%s} p=%d",
						string(consts.ModelProviderBaiZhiCloud),
						strings.TrimSuffix(buLocal, "#"),
						mLocal.Model,
						tpLocal,
						formatParam(pmLocal),
						idxLocal,
					)
					t.Run(testName, func(t *testing.T) {
						t.Parallel()
						keyLocal := <-keyPool
						defer func() { keyPool <- keyLocal }()

						subCtx, subCancel := context.WithTimeout(context.Background(), 30*time.Second)
						defer subCancel()
						callBaseURL := strings.TrimSuffix(buLocal, "#")
						resp, err := mk.CheckModel(subCtx, &domain.CheckModelReq{
							Provider: string(consts.ModelProviderBaiZhiCloud),
							Model:    mLocal.Model,
							BaseURL:  callBaseURL,
							APIKey:   keyLocal,
							Type:     tpLocal,
							Param:    pmLocal,
						})
						if mLocal.Model == "glm-4.5" && pmLocal != nil && pmLocal.SupportImages {
							if err != nil || (resp != nil && resp.Error != "") {
								t.Logf("pass case: %s; response: %+v; error: %v", testName, resp, err)
								return
							}
							t.Logf("failed case: %s", testName)
							t.Errorf("expected error for glm-4.5 with support_images=true, got success: %+v", resp)
							return
						}
						if err != nil {
							t.Logf("failed case: %s", testName)
							t.Errorf("check error: %v", err)
							return
						}
						if resp == nil {
							t.Logf("failed case: %s", testName)
							t.Error("nil response")
							return
						}
						if resp.Error != "" {
							t.Logf("failed case: %s", testName)
							t.Errorf("check response error: %s", resp.Error)
							return
						}
						if tpLocal == "embedding" {
							if !strings.HasPrefix(resp.Content, "dim is :") {
								t.Logf("failed case: %s", testName)
								t.Errorf("embedding content not expected: %q", resp.Content)
								return
							}
						} else {
							if strings.TrimSpace(resp.Content) == "" {
								t.Logf("failed case: %s", testName)
								t.Fatalf("empty content")
							}
						}
						t.Logf("pass case: %s; response: %+v", testName, resp)
					})
				}
			}
		}
	}
}
