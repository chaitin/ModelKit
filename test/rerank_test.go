package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/chaitin/ModelKit/v2/consts"
	"github.com/chaitin/ModelKit/v2/domain"
	"github.com/chaitin/ModelKit/v2/usecase"
)

type rtFunc func(*http.Request) (*http.Response, error)

func (f rtFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestBaiLianRerank(t *testing.T) {
	if os.Getenv("bailianapikey") == "" {
		t.Skip("missing bailianapikey")
	}
	ctx := context.Background()
	mk := usecase.NewModelKit(nil)

	t.Run("default", func(t *testing.T) {
		rk, err := mk.GetReranker(ctx, &domain.ModelMetadata{
			Provider:  consts.ModelProviderBaiLian,
			ModelName: "qwen3-rerank",
			BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank",
			APIKey:    os.Getenv("bailianapikey"),
		})
		if err != nil {
			t.Fatalf("GetReranker: %v", err)
		}
		docs := []string{"火鸡面", "测试文本", "向量模型", "查询", "文档"}
		res, err := rk.Rerank(ctx, domain.RerankRequest{Documents: docs, Query: "测试", ReturnDocuments: false})
		if err != nil {
			t.Skipf("Rerank default: %v", err)
			return
		}
		if len(res.Results) == 0 {
			t.Fatalf("empty results")
		}
		t.Logf("[bailian default] docs=%d", len(docs))
		logResultsHead(t, res.Results)
	})

	t.Run("top_n", func(t *testing.T) {
		rk, err := mk.GetReranker(ctx, &domain.ModelMetadata{
			Provider:  consts.ModelProviderBaiLian,
			ModelName: "qwen3-rerank",
			BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank",
			APIKey:    os.Getenv("bailianapikey"),
		})
		if err != nil {
			t.Fatalf("GetReranker: %v", err)
		}
		docs := []string{"火鸡面", "测试文本", "向量模型", "查询", "文档"}
		for _, n := range []int{1, 3, 9} {
			nn := n
			t.Run(fmt.Sprintf("%d", nn), func(t *testing.T) {
				res, err := rk.Rerank(ctx, domain.RerankRequest{N: &nn, Documents: docs, Query: "测试", ReturnDocuments: false})
				if err != nil {
					t.Skipf("Rerank top_n=%d: %v", nn, err)
					return
				}
				want := nn
				if want > len(docs) {
					want = len(docs)
				}
				if want <= 0 {
					want = 1
				}
				if len(res.Results) != want {
					t.Fatalf("results length mismatch: got %d want %d", len(res.Results), want)
				}
				t.Logf("[bailian top_n=%d]", nn)
				logResultsHead(t, res.Results)
			})
		}
	})

	t.Run("return_documents", func(t *testing.T) {
		rk, err := mk.GetReranker(ctx, &domain.ModelMetadata{
			Provider:  consts.ModelProviderBaiLian,
			ModelName: "qwen3-rerank",
			BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank",
			APIKey:    os.Getenv("bailianapikey"),
		})
		if err != nil {
			t.Fatalf("GetReranker: %v", err)
		}
		docs := []string{"火鸡面", "测试文本", "向量模型"}
		nn := 2
		res, err := rk.Rerank(ctx, domain.RerankRequest{N: &nn, Documents: docs, Query: "食物", ReturnDocuments: true})
		if err != nil {
			t.Skipf("Rerank return_documents: %v", err)
			return
		}
		if len(res.Results) == 0 {
			t.Fatalf("empty results")
		}
		t.Logf("[bailian return_documents] n=%d", nn)
		logResultsHead(t, res.Results)
	})
}

func TestBaiZhiCloudRerank(t *testing.T) {
	if os.Getenv("baizhiapikey") == "" {
		t.Skip("missing baizhiapikey")
	}
	ctx := context.Background()
	mk := usecase.NewModelKit(nil)
	rk, err := mk.GetReranker(ctx, &domain.ModelMetadata{
		Provider:  consts.ModelProviderBaiZhiCloud,
		ModelName: "bge-m3",
		BaseURL:   "https://model-square.app.baizhi.cloud/v1/rerank",
		APIKey:    os.Getenv("baizhiapikey"),
	})
	if err != nil {
		t.Fatalf("GetReranker: %v", err)
	}
	docs := []string{"风哀", "渚回", "羽光"}
	nn := 2
	res, err := rk.Rerank(ctx, domain.RerankRequest{N: &nn, Documents: docs, Query: "检索", ReturnDocuments: false})
	if err != nil {
		t.Skipf("Rerank: %v", err)
		return
	}
	if len(res.Results) == 0 {
		t.Fatalf("empty results")
	}
	t.Logf("[baizhi default] docs=%d", len(docs))
	logResultsHead(t, res.Results)
}

func TestRerankerCombinations(t *testing.T) {
	type blReq struct {
		Model string `json:"model"`
		Input struct {
			Query     string   `json:"query"`
			Documents []string `json:"documents"`
		} `json:"input"`
		Parameters struct {
			ReturnDocuments bool `json:"return_documents"`
			TopN            int  `json:"top_n"`
		} `json:"parameters"`
	}
	type blResp struct {
		Output struct {
			Results []struct {
				Index          int     `json:"index"`
				RelevanceScore float64 `json:"relevance_score"`
				Document       *struct {
					Text string `json:"text"`
				} `json:"document,omitempty"`
			} `json:"results"`
		} `json:"output"`
		Usage *struct {
			PromptTokens int `json:"prompt_tokens"`
			TotalTokens  int `json:"total_tokens"`
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage,omitempty"`
	}
	type baaiReq struct {
		Model           string   `json:"model"`
		Query           string   `json:"query"`
		Documents       []string `json:"documents"`
		ReturnDocuments bool     `json:"return_documents"`
	}
	type baaiResp struct {
		Model   string `json:"model"`
		Results []struct {
			Index          int     `json:"index"`
			RelevanceScore float64 `json:"relavance_score"`
			Document       *struct {
				Text string `json:"text"`
			} `json:"document,omitempty"`
		} `json:"results"`
	}
	type mockTransport struct {
		lastReq  *http.Request
		lastBody []byte
		build    func(*http.Request, []byte) ([]byte, int)
	}
	mt := &mockTransport{}
	mtPtr := mt
	var tr http.RoundTripper = rtFunc(func(req *http.Request) (*http.Response, error) {
		b, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(b))
		mtPtr.lastReq = req
		mtPtr.lastBody = b
		body, code := mtPtr.build(req, b)
		return &http.Response{
			StatusCode: code,
			Status:     fmt.Sprintf("%d %s", code, http.StatusText(code)),
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
			Request:    req,
		}, nil
	})

	prevTr := http.DefaultTransport
	http.DefaultTransport = tr
	defer func() { http.DefaultTransport = prevTr }()

	ctx := context.Background()
	mk := usecase.NewModelKit(nil)

	t.Run("bailian_matrix", func(t *testing.T) {
		baseURLs := []string{
			"",
			"https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank",
			"https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank#",
		}
		for _, bu := range baseURLs {
			for _, n := range []int{0, 3, 9} {
				for _, rd := range []bool{false, true} {
					name := fmt.Sprintf("bailian/%d/%t", n, rd)
					t.Run(name, func(t *testing.T) {
						nn := n
						rdv := rd
						mt.build = func(r *http.Request, b []byte) ([]byte, int) {
							var req blReq
							_ = json.Unmarshal(b, &req)
							docs := req.Input.Documents
							dlen := len(docs)
							topn := req.Parameters.TopN
							if topn <= 0 {
								topn = 1
							}
							if topn > dlen {
								topn = dlen
							}
							resp := blResp{}
							resp.Output.Results = make([]struct {
								Index          int     `json:"index"`
								RelevanceScore float64 `json:"relevance_score"`
								Document       *struct {
									Text string `json:"text"`
								} `json:"document,omitempty"`
							}, dlen)
							for i := 0; i < dlen; i++ {
								resp.Output.Results[i].Index = i
								resp.Output.Results[i].RelevanceScore = float64(i%7) * 0.1
								if req.Parameters.ReturnDocuments {
									resp.Output.Results[i].Document = &struct {
										Text string `json:"text"`
									}{Text: docs[i]}
								}
							}
							resp.Usage = &struct {
								PromptTokens int `json:"prompt_tokens"`
								TotalTokens  int `json:"total_tokens"`
								InputTokens  int `json:"input_tokens"`
								OutputTokens int `json:"output_tokens"`
							}{PromptTokens: 1, TotalTokens: 2, InputTokens: 1, OutputTokens: 1}
							out, _ := json.Marshal(resp)
							return out, http.StatusOK
						}
						rk, err := mk.GetReranker(ctx, &domain.ModelMetadata{
							Provider:  consts.ModelProviderBaiLian,
							ModelName: "qwen3-rerank",
							BaseURL:   bu,
							APIKey:    "test",
						})
						if err != nil {
							t.Fatalf("GetReranker: %v", err)
						}
						docs := []string{"a", "b", "c", "d", "e"}
						res, err := rk.Rerank(ctx, domain.RerankRequest{N: &nn, Documents: docs, Query: "q", ReturnDocuments: rdv})
						if err != nil {
							t.Fatalf("Rerank: %v", err)
						}
						var reqBody blReq
						_ = json.Unmarshal(mtPtr.lastBody, &reqBody)
						wantEndpoint := bu
						if wantEndpoint == "" {
							wantEndpoint = "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank"
						}
						if strings.HasSuffix(wantEndpoint, "#") {
							wantEndpoint = strings.TrimSuffix(wantEndpoint, "#")
						}
						if mtPtr.lastReq.URL.String() != wantEndpoint {
							t.Fatalf("endpoint mismatch: got %s want %s", mtPtr.lastReq.URL.String(), wantEndpoint)
						}
						wantTopN := nn
						if wantTopN <= 0 {
							wantTopN = 1
						}
						if wantTopN > len(docs) {
							wantTopN = len(docs)
						}
						if reqBody.Parameters.TopN != wantTopN {
							t.Fatalf("top_n mismatch: got %d want %d", reqBody.Parameters.TopN, wantTopN)
						}
						if reqBody.Parameters.ReturnDocuments != rdv {
							t.Fatalf("return_documents mismatch: got %t want %t", reqBody.Parameters.ReturnDocuments, rdv)
						}
						if len(res.Results) != wantTopN {
							t.Fatalf("results length mismatch: got %d want %d", len(res.Results), wantTopN)
						}
						if rdv {
							if res.Results[0].Document == "" {
								t.Fatalf("document empty when return_documents=true")
							}
						} else {
							if res.Results[0].Document != "" {
								t.Fatalf("document should be empty when return_documents=false")
							}
						}
						if res.Usage == nil || res.Usage.TotalTokens != 2 {
							t.Fatalf("usage mismatch")
						}
					})
				}
			}
		}
	})

	t.Run("baai_matrix", func(t *testing.T) {
		baseURLs := []string{
			"https://model-square.app.baizhi.cloud/v1",
			"https://model-square.app.baizhi.cloud/v1/rerank",
		}
		for _, bu := range baseURLs {
			for _, n := range []int{0, 3, 9} {
				for _, rd := range []bool{false, true} {
					name := fmt.Sprintf("baai/%d/%t", n, rd)
					t.Run(name, func(t *testing.T) {
						nn := n
						rdv := rd
						mt.build = func(r *http.Request, b []byte) ([]byte, int) {
							var req baaiReq
							_ = json.Unmarshal(b, &req)
							docs := req.Documents
							dlen := len(docs)
							resp := baaiResp{Model: "bge-m3"}
							resp.Results = make([]struct {
								Index          int     `json:"index"`
								RelevanceScore float64 `json:"relavance_score"`
								Document       *struct {
									Text string `json:"text"`
								} `json:"document,omitempty"`
							}, dlen)
							for i := 0; i < dlen; i++ {
								resp.Results[i].Index = i
								resp.Results[i].RelevanceScore = float64(i%7) * 0.1
								if req.ReturnDocuments {
									resp.Results[i].Document = &struct {
										Text string `json:"text"`
									}{Text: docs[i]}
								}
							}
							out, _ := json.Marshal(resp)
							return out, http.StatusOK
						}
						rk, err := mk.GetReranker(ctx, &domain.ModelMetadata{
							Provider:  consts.ModelProviderBaiZhiCloud,
							ModelName: "bge-m3",
							BaseURL:   bu,
							APIKey:    "test",
						})
						if err != nil {
							t.Fatalf("GetReranker: %v", err)
						}
						docs := []string{"a", "b", "c", "d"}
						res, err := rk.Rerank(ctx, domain.RerankRequest{N: &nn, Documents: docs, Query: "q", ReturnDocuments: rdv})
						if err != nil {
							t.Fatalf("Rerank: %v", err)
						}
						wantEndpoint := bu
						if !strings.HasSuffix(wantEndpoint, "/rerank") {
							wantEndpoint = wantEndpoint + "/rerank"
						}
						if mtPtr.lastReq.URL.String() != wantEndpoint {
							t.Fatalf("endpoint mismatch: got %s want %s", mtPtr.lastReq.URL.String(), wantEndpoint)
						}
						want := nn
						if want <= 0 {
							want = 1
						}
						if want > len(docs) {
							want = len(docs)
						}
						if len(res.Results) != want {
							t.Fatalf("results length mismatch: got %d want %d", len(res.Results), want)
						}
						if rdv {
							if res.Results[0].Document == "" {
								t.Fatalf("document empty when return_documents=true")
							}
						} else {
							if res.Results[0].Document != "" {
								t.Fatalf("document should be empty when return_documents=false")
							}
						}
					})
				}
			}
		}
	})
}

func logResultsHead(t *testing.T, v []domain.Result) {
	n := 5
	if len(v) < n {
		n = len(v)
	}
	s := make([]string, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, fmt.Sprintf("%d:%.4f", v[i].Index, v[i].RelevanceScore))
	}
	t.Logf("head: %s", strings.Join(s, " "))
}
