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
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func init() {
	RegisterSubtaskMeta(&CollectApiJobsMeta)
}

const RAW_JOB_TABLE = "gitlab_api_job"

type SimpleGitlabApiJob struct {
	GitlabId  int
	CreatedAt common.Iso8601Time `json:"created_at"`
}

var CollectApiJobsMeta = plugin.SubTaskMeta{
	Name:             "collectApiJobs",
	EntryPoint:       CollectApiJobs,
	EnabledByDefault: true,
	Description:      "Collect job data from gitlab api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	Dependencies:     []*plugin.SubTaskMeta{&ExtractApiPipelineDetailsMeta},
}

func CollectApiJobs(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_JOB_TABLE)
	collector, err := helper.NewStatefulApiCollectorForFinalizableEntity(helper.FinalizableApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		CollectNewRecordsByList: helper.FinalizableApiCollectorListArgs{
			PageSize:    100,
			Concurrency: 10,
			FinalizableApiCollectorCommonArgs: helper.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "projects/{{ .Params.ProjectId }}/jobs",
				Query: func(reqData *helper.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					query := url.Values{}
					query.Set("page", strconv.Itoa(reqData.Pager.Page))
					query.Set("per_page", strconv.Itoa(reqData.Pager.Size))
					query.Set("scope[]", "failed")
					query.Add("scope[]", "success")
					return query, nil
				},
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					var items []json.RawMessage
					err := helper.UnmarshalResponse(res, &items)
					if err != nil {
						return nil, err
					}
					return items, nil
				},
				AfterResponse: ignoreHTTPStatus403, // ignore 403 for CI/CD disable
			},
			GetCreated: func(item json.RawMessage) (time.Time, errors.Error) {
				pr := &SimpleGitlabApiJob{}
				err := json.Unmarshal(item, pr)
				if err != nil {
					return time.Time{}, errors.BadInput.Wrap(err, "failed to unmarshal gitlab job")
				}
				return pr.CreatedAt.ToTime(), nil
			},
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
