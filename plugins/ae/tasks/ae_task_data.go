package tasks

import "github.com/merico-dev/lake/plugins/helper"

type AeOptions struct {
	ProjectId int
	Tasks     []string `json:"tasks,omitempty"`
}

type AeTaskData struct {
	Options   *AeOptions
	ApiClient *helper.ApiAsyncClient
}
type AeApiParams struct {
	ProjectId int
}
