package domain

import (

	"github.com/chaitin/ModelKit/backend/consts"
	"github.com/rokku-c/go-openai"
)

type CompletionRequest struct {
	openai.CompletionRequest
	Metadata map[string]string `json:"metadata"`
}

type ModelListResp struct {
	Object string       `json:"object"`
	Data   []*ModelData `json:"data"`
}

type ModelData struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Created   int64  `json:"created"`
	OwnedBy   string `json:"owned_by"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	BaseModel string `json:"base_model"`
	APIBase   string `json:"api_base"`
	IsActive  bool   `json:"is_active"`
}



type ConfigReq struct {
	Type    consts.ConfigType `json:"type" query:"type"`
	Key     string            `json:"-"`
	BaseURL string            `json:"-"`
}

type ConfigResp struct {
	Type    consts.ConfigType `json:"type"`
	Content string            `json:"content"`
}
type OpenAIResp struct {
	Object string        `json:"object"`
	Data   []*OpenAIData `json:"data"`
}

type OpenAIData struct {
	ID string `json:"id"`
}
