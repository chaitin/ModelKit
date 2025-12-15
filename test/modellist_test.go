package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/chaitin/ModelKit/v2/domain"
	"github.com/chaitin/ModelKit/v2/usecase"
)

func TestModelList_BaiZhiCloud_ByType(t *testing.T) {
	ctx := context.Background()
	mk := usecase.NewModelKit(nil)

	apiKey := strings.TrimSpace(os.Getenv("baizhiapikey"))
	if apiKey == "" {
		t.Fatalf("missing baizhiapikey")
	}

	baseURL := "https://model-square.app.baizhi.cloud/v1"
	provider := "baizhi-cloud"

	types := []string{"chat", "embedding", "rerank", "code", "analysis", "analysis-vl"}

	for _, tp := range types {
		tpLocal := tp
		name := fmt.Sprintf("provider=%s type=%s base_url=%s apiKey=present", provider, tpLocal, baseURL)
		t.Run(name, func(t *testing.T) {
			ml, err := mk.ModelList(ctx, &domain.ModelListReq{
				Provider: provider,
				BaseURL:  baseURL,
				APIKey:   apiKey,
				Type:     tpLocal,
			})
			if err != nil {
				mlJSON, _ := json.Marshal(ml)
				t.Logf("test case: %s", name)
				t.Logf("modellist_resp=%s", string(mlJSON))
				t.Fatalf("ModelList error: %v", err)
			}

			curlModels, curlErr := curlAllModels(ctx, baseURL, apiKey)
			if curlErr != nil {
				mlJSON, _ := json.Marshal(ml)
				t.Logf("test case: %s", name)
				t.Logf("modellist_resp=%s", string(mlJSON))
				t.Fatalf("curl models error: %v", curlErr)
			}

			filtered := usecase.FilterModelsByType(curlModels, &domain.ModelListReq{Provider: provider, Type: tpLocal})

			mlJSON, _ := json.Marshal(ml)
			filteredJSON, _ := json.Marshal(filtered)
			t.Logf("test case: %s", name)
			t.Logf("modellist_resp=%s", string(mlJSON))
			t.Logf("curl_filtered=%s", string(filteredJSON))

			if !sameModelSets(ml.Models, filtered) {
				t.Fatalf("result mismatch: modellist_count=%d curl_filtered_count=%d", len(ml.Models), len(filtered))
			}
		})
	}
}

func curlAllModels(ctx context.Context, baseURL, apiKey string) ([]domain.ModelListItem, error) {
	url := strings.TrimSuffix(baseURL, "/") + "/models"
	cmd := exec.CommandContext(ctx, "curl", "-s", "-H", "Authorization: Bearer "+apiKey, url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("curl failed: %v, output=%s", err, string(out))
	}
	var resp struct {
		Object string `json:"object"`
		Data   []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if uerr := json.Unmarshal(out, &resp); uerr != nil {
		return nil, fmt.Errorf("unmarshal failed: %v, body=%s", uerr, string(out))
	}
	items := make([]domain.ModelListItem, 0, len(resp.Data))
	for _, d := range resp.Data {
		items = append(items, domain.ModelListItem{Model: d.ID})
	}
	return items, nil
}

func sameModelSets(a, b []domain.ModelListItem) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]int, len(a))
	for _, it := range a {
		m[it.Model]++
	}
	for _, it := range b {
		if m[it.Model] == 0 {
			return false
		}
		m[it.Model]--
	}
	for _, v := range m {
		if v != 0 {
			return false
		}
	}
	return true
}
