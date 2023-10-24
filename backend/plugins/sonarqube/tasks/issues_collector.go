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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_ISSUES_TABLE = "sonarqube_api_issues"

var _ plugin.SubTaskEntryPoint = CollectIssues

type SonarqubeIssueIteratorNode struct {
	Severity string
	Status   string
	Type     string
}

func CollectIssues(taskCtx plugin.SubTaskContext) (err errors.Error) {
	logger := taskCtx.GetLogger()
	logger.Info("collect issues")

	iterator := helper.NewQueueIterator()
	severities := []string{"BLOCKER", "CRITICAL", "MAJOR", "MINOR", "INFO"}
	statuses := []string{"OPEN", "CONFIRMED", "REOPENED", "RESOLVED", "CLOSED"}
	types := []string{"BUG", "VULNERABILITY", "CODE_SMELL"}
	for _, severity := range severities {
		for _, status := range statuses {
			for _, typ := range types {
				iterator.Push(
					&SonarqubeIssueIteratorNode{
						Severity: severity,
						Status:   status,
						Type:     typ,
					},
				)
			}
		}
	}
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUES_TABLE)
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		UrlTemplate:        "issues/search",
		Input:              iterator,
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("componentKeys", fmt.Sprintf("%v", data.Options.ProjectKey))
			query.Set("severities", reqData.Input.(*SonarqubeIssueIteratorNode).Severity)
			query.Set("statuses", reqData.Input.(*SonarqubeIssueIteratorNode).Status)
			query.Set("types", reqData.Input.(*SonarqubeIssueIteratorNode).Type)
			query.Set("p", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("ps", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Encode()
			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var resData struct {
				Data []json.RawMessage `json:"issues"`
			}
			err = helper.UnmarshalResponse(res, &resData)
			if err != nil {
				return nil, err
			}

			// check if sonar report updated during collecting
			var issue struct {
				UpdateDate *common.Iso8601Time `json:"updateDate"`
			}
			for _, v := range resData.Data {
				err = errors.Convert(json.Unmarshal(v, &issue))
				if err != nil {
					return nil, err
				}
				if issue.UpdateDate.ToTime().After(data.TaskStartTime) {
					return nil, errors.Default.New(fmt.Sprintf(`Your data is affected by the latest analysis\n
						Please recollect this project: %s`, data.Options.ProjectKey))
				}
			}

			return resData.Data, nil
		},
	})
	if err != nil {
		return err
	}
	return collector.Execute()

}

var CollectIssuesMeta = plugin.SubTaskMeta{
	Name:             "CollectIssues",
	EntryPoint:       CollectIssues,
	EnabledByDefault: true,
	Description:      "Collect issues data from Sonarqube api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}
