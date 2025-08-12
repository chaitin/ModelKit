package domain

type OpenAIData struct {
	ID string `json:"id"`
}

type OpenAIResp struct {
	Object string        `json:"object"`
	Data   []*OpenAIData `json:"data"`
}