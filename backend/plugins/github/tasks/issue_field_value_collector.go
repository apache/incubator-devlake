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
	"reflect"

	"github.com/apache/devlake/core/dal"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&CollectApiIssueFieldValuesMeta)
}

const RAW_ISSUE_FIELD_VALUE_TABLE = "github_api_issue_field_values"

// SimpleIssue carries the two identifiers the field-value endpoint needs: the number to
// build the URL, and the id to key the extracted rows on.
type SimpleIssue struct {
	Number   int
	GithubId int
}

var CollectApiIssueFieldValuesMeta = plugin.SubTaskMeta{
	Name:             "Collect Issue Field Values",
	EntryPoint:       CollectApiIssueFieldValues,
	EnabledByDefault: false,
	Description:      "Collect issue field values from the GitHub API, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	DependencyTables: []string{models.GithubIssue{}.TableName()},
	ProductTables:    []string{RAW_ISSUE_FIELD_VALUE_TABLE},
}

func CollectApiIssueFieldValues(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)

	apiCollector, err := helper.NewStatefulApiCollector(helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: GithubApiParams{
			ConnectionId: data.Options.ConnectionId,
			Name:         data.Options.Name,
		},
		Table: RAW_ISSUE_FIELD_VALUE_TABLE,
	})
	if err != nil {
		return err
	}

	clauses := []dal.Clause{
		dal.Select("number, github_id"),
		dal.From(models.GithubIssue{}.TableName()),
		dal.Where("repo_id = ? and connection_id = ?", data.Options.GithubId, data.Options.ConnectionId),
	}
	if apiCollector.IsIncremental() && apiCollector.GetSince() != nil {
		clauses = append(clauses, dal.Where("github_updated_at > ?", apiCollector.GetSince()))
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}

	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleIssue{}))
	if err != nil {
		return err
	}

	err = apiCollector.InitCollector(helper.ApiCollectorArgs{
		ApiClient: data.ApiClient,
		PageSize:  100,
		Input:     iterator,

		UrlTemplate: "repos/{{ .Params.Name }}/issues/{{ .Input.Number }}/issue-field-values",

		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			// An organization that has never configured issue fields, or a token without
			// visibility of them, gets 404 rather than an empty list. Treat that as "this
			// issue has no field values" so one such repository cannot fail the whole task.
			if res.StatusCode == http.StatusNotFound {
				return nil, nil
			}
			var items []json.RawMessage
			err := helper.UnmarshalResponse(res, &items)
			if err != nil {
				return nil, err
			}
			return items, nil
		},
	})
	if err != nil {
		return err
	}
	return apiCollector.Execute()
}
