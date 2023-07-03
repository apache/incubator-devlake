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
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&CollectApiEventsMeta)
}

const RAW_EVENTS_TABLE = "github_api_events"

type SimpleGithubApiEvents struct {
	GithubId  int64
	CreatedAt helper.Iso8601Time `json:"created_at"`
}

var CollectApiEventsMeta = plugin.SubTaskMeta{
	Name:             "collectApiEvents",
	EntryPoint:       CollectApiEvents,
	EnabledByDefault: true,
	Description:      "Collect Events data from Github api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	DependencyTables: []string{},
	ProductTables:    []string{RAW_EVENTS_TABLE},
}

func CollectApiEvents(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)
	db := taskCtx.GetDal()

	collector, err := helper.NewStatefulApiCollectorForFinalizableEntity(helper.FinalizableApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_EVENTS_TABLE,
		},
		ApiClient: data.ApiClient,
		TimeAfter: data.TimeAfter, // set to nil to disable timeFilter
		CollectNewRecordsByList: helper.FinalizableApiCollectorListArgs{
			PageSize:    100,
			Concurrency: 10,
			//GetTotalPages: GetTotalPagesFromResponse,
			FinalizableApiCollectorCommonArgs: helper.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "repos/{{ .Params.Name }}/issues/events",
				Query: func(reqData *helper.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					query := url.Values{}
					query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
					query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
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
				AfterResponse: ignoreHTTPStatus422,
			},
			GetCreated: func(item json.RawMessage) (time.Time, errors.Error) {
				e := &SimpleGithubApiEvents{}
				err := json.Unmarshal(item, e)
				if err != nil {
					return time.Time{}, errors.BadInput.Wrap(err, "failed to unmarshal github events")
				}
				return e.CreatedAt.ToTime(), nil
			},
		},
		CollectUnfinishedDetails: helper.FinalizableApiCollectorDetailArgs{
			BuildInputIterator: func() (helper.Iterator, errors.Error) {
				cursor, err := db.Cursor(
					dal.Select("github_id"),
					dal.From(&models.GithubIssueEvent{}),
					dal.Where(
						"github_id = ? AND connection_id = ? and type NOT IN ('closed', 'merged')",
						data.Options.GithubId, data.Options.ConnectionId,
					),
				)
				if err != nil {
					return nil, err
				}
				return helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleGithubApiEvents{}))
			},
			FinalizableApiCollectorCommonArgs: helper.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "repos/{{ .Params.Name }}/issues/events/{{ .Input.GithubId }}",
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					body, err := io.ReadAll(res.Body)
					if err != nil {
						return nil, errors.Convert(err)
					}
					res.Body.Close()
					return []json.RawMessage{body}, nil
				},
				AfterResponse: ignoreHTTPStatus422,
			},
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}
