package apiv2models

import (
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
