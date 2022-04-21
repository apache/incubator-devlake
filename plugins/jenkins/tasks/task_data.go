package tasks

import (
	"time"

	"github.com/merico-dev/lake/plugins/helper"
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
	ApiClient *helper.ApiAsyncClient
	Since     *time.Time
}
