package tasks

import (
	"github.com/merico-dev/lake/plugins/github/models"
	"time"

	"github.com/merico-dev/lake/plugins/core"
)

type GithubOptions struct {
	Tasks       []string `json:"tasks,omitempty"`
	Since       string
	Owner       string
	Repo        string
	ParamString string
}

type GithubTaskData struct {
	Options   *GithubOptions
	ApiClient *core.ApiClient
	Since     *time.Time
	Repo      *models.GithubRepo
}
