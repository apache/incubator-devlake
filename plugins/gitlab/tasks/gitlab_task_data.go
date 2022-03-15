package tasks

import (
	"time"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
)

type GitlabOptions struct {
	SourceId  uint64   `json:"sourceId"`
	ProjectId int      `json:"projectId"`
	Tasks     []string `json:"tasks,omitempty"`
	//Since    string
}

type GitlabTaskData struct {
	Options       *GitlabOptions
	ApiClient     *core.ApiClient
	ProjectCommit *models.GitlabProjectCommit
	Since         *time.Time
}
