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
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_STAGE_TABLE = "jenkins_api_stages"

var CollectApiStagesMeta = core.SubTaskMeta{
	Name:             "collectApiStages",
	EntryPoint:       CollectApiStages,
	EnabledByDefault: true,
	Description:      "Collect stages data from jenkins api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

type SimpleBuild struct {
	JobName     string
	Number      string
	DisplayName string
}

func CollectApiStages(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*JenkinsTaskData)
	clauses := []dal.Clause{
		dal.Select("tjb.job_name,tjb.number, tjb.display_name"),
		dal.From("_tool_jenkins_builds tjb"),
		dal.Where(`tjb.connection_id = ? and tjb.class = ?`,
			data.Options.ConnectionId, "WorkflowRun"),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleBuild{}))
	if err != nil {
		return err
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
			},
			Ctx:   taskCtx,
			Table: RAW_STAGE_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Input:       iterator,
		UrlTemplate: "job/{{ .Input.JobName }}/{{ .Input.Number }}/wfapi/describe",
		/*
			(Optional) Return query string for request, or you can plug them into UrlTemplate directly
		*/
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Stages []json.RawMessage `json:"stages"`
			}
			err := helper.UnmarshalResponse(res, &data)
			if err != nil {
				return nil, err
			}
			return data.Stages, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
