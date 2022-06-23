/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_JOB_TABLE = "jenkins_api_jobs"

var CollectApiJobsMeta = core.SubTaskMeta{
	Name:             "collectApiJobs",
	EntryPoint:       CollectApiJobs,
	EnabledByDefault: true,
	Description:      "Collect jobs data from jenkins api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func CollectApiJobs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JenkinsTaskData)
	incremental := false
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
			},
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
		// jenkins api is special, 1. If the concurrency is larger than 1, then it will report 500.
		Concurrency: 1,

		UrlTemplate: "api/json",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			treeValue := fmt.Sprintf(
				"jobs[name,class,color,base]{%d,%d}",
				reqData.Pager.Skip, reqData.Pager.Skip+reqData.Pager.Size)
			query.Set("tree", treeValue)
			return query, nil
		},

		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Jobs []json.RawMessage `json:"jobs"`
			}
			err := helper.UnmarshalResponse(res, &data)
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
