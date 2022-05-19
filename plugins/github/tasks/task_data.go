package tasks

import (
	"time"

	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type GithubOptions struct {
	Tasks []string `json:"tasks,omitempty"`
	Since string
	Owner string
	Repo  string
	models.Config
}

type GithubTaskData struct {
	Options   *GithubOptions
	ApiClient *helper.ApiAsyncClient
	Since     *time.Time
	Repo      *models.GithubRepo
}
