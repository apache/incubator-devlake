package apiv2models

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/jira/models"
)

type Project struct {
	Self string `json:"self"`
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

func (p Project) ToToolLayer(sourceId uint64) *models.JiraProject {
	return &models.JiraProject{
		SourceId: sourceId,
		Id:       p.ID,
		Key:      p.Key,
		Name:     p.Name,
	}
}
func (Project) FromAPI(sourceId uint64, raw json.RawMessage) (interface{}, error) {
	var vv []Project
	err := json.Unmarshal(raw, &vv)
	if err != nil {
		return nil, err
	}
	list := make([]*models.JiraProject, len(vv))
	for i, item := range vv {
		list[i] = item.ToToolLayer(sourceId)
	}
	return list, nil
}
func (Project) ExtractRawMessage(blob []byte) (json.RawMessage, error) {
	return blob, nil
}
