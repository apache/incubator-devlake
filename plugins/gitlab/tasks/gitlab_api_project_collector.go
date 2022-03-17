package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

const RAW_PROJECT_TABLE = "gitlab_api_project"

func CollectApiProject(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		Incremental:        false,
		UrlTemplate:        "projects/{{ .Params.ProjectId }}",
		Query:              GetQuery,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			rawMessages := make([]json.RawMessage, 1)

			defer res.Body.Close()
			resBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, fmt.Errorf("%w %s", err, res.Request.URL.String())
			}
			err = json.Unmarshal(resBody, &(rawMessages[0]))
			if err != nil {
				return nil, fmt.Errorf("%w %s %s", err, res.Request.URL.String(), string(resBody))
			}

			return rawMessages, nil
		}})

	if err != nil {
		return err
	}

	return collector.Execute()
}
