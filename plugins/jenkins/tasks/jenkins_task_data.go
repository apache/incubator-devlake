package tasks

import (
	"github.com/merico-dev/lake/plugins/core"
	"time"
)

type JenkinsOptions struct {
	Host     string
	Username string
	Password string
	Since    string
	Tasks    []string `json:"tasks,omitempty"`
}

type JenkinsTaskData struct {
	Options   *JenkinsOptions
	ApiClient *core.ApiClient
	Since     *time.Time
}
