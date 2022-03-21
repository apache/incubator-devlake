package tasks

import (
	"time"

	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
)

type GithubOptions struct {
	Tasks []string `json:"tasks,omitempty"`
	Since string
	Owner string
	Repo  string
}

type GithubTaskData struct {
	Options   *GithubOptions
	ApiClient *helper.ApiAsyncClient
	Since     *time.Time
	Repo      *models.GithubRepo
}
