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
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

const RAW_JOB_TABLE = "gitlab_api_job"

type SimpleGitlabApiJob struct {
	GitlabId  int
	CreatedAt helper.Iso8601Time `json:"created_at"`
}

var CollectApiJobsMeta = plugin.SubTaskMeta{
	Name:             "collectApiJobs",
	EntryPoint:       CollectApiJobs,
	EnabledByDefault: true,
	Description:      "Collect job data from gitlab api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func CollectApiJobs(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_JOB_TABLE)
	db := taskCtx.GetDal()
	collector, err := helper.NewStatefulApiCollectorForFinalizableEntity(helper.FinalizableApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		TimeAfter:          data.TimeAfter, // set to nil to disable timeFilter
		CollectNewRecordsByList: helper.FinalizableApiCollectorListArgs{
			PageSize:    100,
			Concurrency: 10,
			FinalizableApiCollectorCommonArgs: helper.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "projects/{{ .Params.ProjectId }}/jobs",
				Query: func(reqData *helper.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					query := url.Values{}
					query.Set("page", strconv.Itoa(reqData.Pager.Page))
					query.Set("per_page", strconv.Itoa(reqData.Pager.Size))
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
		CollectUnfinishedDetails: helper.FinalizableApiCollectorDetailArgs{
			BuildInputIterator: func() (helper.Iterator, errors.Error) {
				// select pull id from database
				cursor, err := db.Cursor(
					dal.Select("gitlab_id"),
					dal.From(&models.GitlabJob{}),
					dal.Where(
						"project_id = ? AND connection_id = ? AND finished_at is null",
						data.Options.ProjectId, data.Options.ConnectionId,
					),
				)
				if err != nil {
					return nil, err
				}
				return helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleGitlabApiJob{}))
			},
			FinalizableApiCollectorCommonArgs: helper.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "projects/{{ .Params.ProjectId }}/jobs/{{ .Input.GitlabId }}",
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					body, err := io.ReadAll(res.Body)
					if err != nil {
						return nil, errors.Convert(err)
					}
					res.Body.Close()
					return []json.RawMessage{body}, nil
				},
				AfterResponse: ignoreHTTPStatus403, // ignore 403 for CI/CD disable
			},
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
