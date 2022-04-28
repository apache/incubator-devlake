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
	Config
}

type Config struct {
	GITHUB_PR_TYPE                string `json:"GITHUB_PR_TYPE,omitempty"`
	GITHUB_PR_COMPONENT           string `json:"GITHUB_PR_COMPONENT,omitempty"`
	GITHUB_ISSUE_SEVERITY         string `json:"GITHUB_ISSUE_SEVERITY,omitempty"`
	GITHUB_ISSUE_COMPONENT        string `json:"GITHUB_ISSUE_COMPONENT,omitempty"`
	GITHUB_ISSUE_PRIORITY         string `json:"GITHUB_ISSUE_PRIORITY,omitempty"`
	GITHUB_ISSUE_TYPE_REQUIREMENT string `json:"GITHUB_ISSUE_TYPE_REQUIREMENT,omitempty"`
	GITHUB_ISSUE_TYPE_BUG         string `json:"GITHUB_ISSUE_TYPE_BUG,omitempty"`
	GITHUB_ISSUE_TYPE_INCIDENT    string `json:"GITHUB_ISSUE_TYPE_INCIDENT,omitempty"`
}

type GithubTaskData struct {
	Options   *GithubOptions
	ApiClient *helper.ApiAsyncClient
	Since     *time.Time
	Repo      *models.GithubRepo
}
