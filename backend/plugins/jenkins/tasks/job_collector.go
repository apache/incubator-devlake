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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

var CollectApiJobsMeta = plugin.SubTaskMeta{
	Name:             "collectApiJobs",
	EntryPoint:       CollectApiJobs,
	EnabledByDefault: true,
	Description:      "Collect jobs data from multibranch projects using jenkins api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func CollectApiJobs(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JenkinsTaskData)
	logger := taskCtx.GetLogger()

	if data.Options.Class != WORKFLOW_MULTI_BRANCH_PROJECT {
		logger.Debug("class must be %s, got %s", WORKFLOW_MULTI_BRANCH_PROJECT, data.Options.Class)
		return nil
	}

	collector, err := helper.NewStatefulApiCollectorForFinalizableEntity(helper.FinalizableApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
				FullName:     data.Options.JobFullName,
			},
			Ctx:   taskCtx,
			Table: RAW_JOB_TABLE,
		},
		ApiClient: data.ApiClient,
		CollectNewRecordsByList: helper.FinalizableApiCollectorListArgs{
			FinalizableApiCollectorCommonArgs: helper.FinalizableApiCollectorCommonArgs{
				UrlTemplate: fmt.Sprintf("%sjob/%s/api/json", data.Options.JobPath, data.Options.JobName),
				Query: func(reqData *helper.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					query := url.Values{}
					treeValue := "jobs[fullName,name,class,url,color,description]"
					query.Set("tree", treeValue)
					return query, nil
				},
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					var data struct {
						Jobs []json.RawMessage `json:"jobs"`
					}
					err := helper.UnmarshalResponse(res, &data)
					if err != nil {
						return nil, err
					}

					jobs := make([]json.RawMessage, 0, len(data.Jobs))
					for _, job := range data.Jobs {
						var jobObj map[string]interface{}
						err := json.Unmarshal(job, &jobObj)
						if err != nil {
							return nil, errors.Convert(err)
						}

						logger.Debug("%v", jobObj)
						if jobObj["color"] != "notbuilt" && jobObj["color"] != "nobuilt_anime" {
							jobs = append(jobs, job)
						}
					}

					return jobs, nil
				},
			},
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
