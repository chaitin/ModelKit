package domain

type GPUStackListModelResp struct {
	Items []*struct {
		Name string `json:"name"`
	} `json:"items"`
}

// ParseModels 实现ModelResponseParser接口
func (o *GPUStackListModelResp) ParseModels() []ModelListItem {
	var models []ModelListItem
	for _, item := range o.Items {
		models = append(models, ModelListItem{Model: item.Name})
	}
	return models
}
