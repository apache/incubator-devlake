package tasks

import (
	"time"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
)

type JiraOptions struct {
	SourceId uint64   `json:"sourceId"`
	BoardId  uint64   `json:"boardId"`
	Tasks    []string `json:"tasks,omitempty"`
	Since    string
}

type JiraTaskData struct {
	Options   *JiraOptions
	ApiClient *core.ApiClient
	Source    *models.JiraSource
	Since     *time.Time
}
