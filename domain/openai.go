package domain

type OpenAIData struct {
	ID string `json:"id"`
}

type OpenAIResp struct {
	Object string        `json:"object"`
	Data   []*OpenAIData `json:"data"`
}

// ParseModels 实现ModelResponseParser接口
func (o *OpenAIResp) ParseModels() []ModelListItem {
	var models []ModelListItem
	for _, item := range o.Data {
		models = append(models, ModelListItem{Model: item.ID})
	}
	return models
}
