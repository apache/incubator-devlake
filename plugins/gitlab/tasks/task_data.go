package tasks

import (
	"time"

	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/helper"
)

type GitlabOptions struct {
	ConnectionId uint64   `json:"connectionId"`
	ProjectId    int      `json:"projectId"`
	Tasks        []string `json:"tasks,omitempty"`
	//Since    string
}

type GitlabTaskData struct {
	Options       *GitlabOptions
	ApiClient     *helper.ApiAsyncClient
	ProjectCommit *models.GitlabProjectCommit
	Since         *time.Time
}
