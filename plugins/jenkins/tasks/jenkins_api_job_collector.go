package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"net/http"
)

const RAW_JOB_TABLE = "jenkins_api_jobs"

var _ core.SubTaskEntryPoint = CollectApiJobs

func CollectApiJobs(taskCtx core.SubTaskContext) error {
	//db := taskCtx.GetDb()
	data := taskCtx.GetData().(*JenkinsTaskData)

	//since := data.Since
	incremental := false
	// user didn't specify a time range to sync, try load from database

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			/*
				Table store raw data
			*/
			Table: RAW_JOB_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: incremental,

		UrlTemplate: "api/json",

		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Jobs []json.RawMessage `json:"jobs"`
			}
			err := core.UnmarshalResponse(res, &data)
			if err != nil {
				return nil, err
			}
			return data.Jobs, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
