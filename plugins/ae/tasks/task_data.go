package tasks

import "github.com/apache/incubator-devlake/plugins/helper"

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
