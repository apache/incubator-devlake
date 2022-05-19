package apiv2models

import (
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

type Project struct {
	Self string `json:"self"`
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

func (p Project) ToToolLayer(connectionId uint64) *models.JiraProject {
	return &models.JiraProject{
		ConnectionId: connectionId,
		Id:           p.ID,
		Key:          p.Key,
		Name:         p.Name,
	}
}
