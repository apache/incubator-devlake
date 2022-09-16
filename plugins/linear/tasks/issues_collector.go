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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/linear/models"
	"github.com/merico-dev/graphql"
)

const RAW_ISSUES_TABLE = "linear_issues"

type GraphqlQueryIssueWrapper struct {
	IssueList struct {
		TotalCount graphql.Int
		Issues     []GraphqlQueryIssue `graphql:"nodes"`
		PageInfo   *helper.GraphqlQueryPageInfo
	} `graphql:"issues(first: $pageSize, after: $skipCursor)"`
}

type GraphqlQueryIssue struct {
	Id         string
	Identifier string
	Title      string
	Number     int
	BranchName string
}

var _ core.SubTaskEntryPoint = CollectIssues

var CollectIssuesMeta = core.SubTaskMeta{
	Name:             "CollectIssues",
	EntryPoint:       CollectIssues,
	EnabledByDefault: true,
	Description:      "Collect Issues data from Linear api",
}

func CollectIssues(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*LinearTaskData)

	collector, err := helper.NewGraphqlCollector(helper.GraphqlCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:    taskCtx,
			Params: LinearApiParams{},
			Table:  RAW_ISSUES_TABLE,
		},
		GraphqlClient: data.GraphqlClient,
		PageSize:      100,

		BuildQuery: func(reqData *helper.GraphqlRequestData) (interface{}, map[string]interface{}, error) {
			query := &GraphqlQueryIssueWrapper{}
			variables := map[string]interface{}{
				"pageSize":   graphql.Int(reqData.Pager.Size),
				"skipCursor": (*graphql.String)(reqData.Pager.SkipCursor),
			}
			return query, variables, nil
		},
		GetPageInfo: func(iQuery interface{}, args *helper.GraphqlCollectorArgs) (*helper.GraphqlQueryPageInfo, error) {
			query := iQuery.(*GraphqlQueryIssueWrapper)
			return query.IssueList.PageInfo, nil
		},
		ResponseParser: func(iQuery interface{}, variables map[string]interface{}) ([]interface{}, error) {
			query := iQuery.(*GraphqlQueryIssueWrapper)
			issues := query.IssueList.Issues

			results := make([]interface{}, 0, 1)
			for _, issue := range issues {
				githubIssue, err := convertLinearIssue(issue, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				results = append(results, githubIssue)
			}
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

func convertLinearIssue(issue GraphqlQueryIssue, connectionId uint64) (*models.LinearIssue, error) {
	linearIssue := &models.LinearIssue{
		ConnectionId: connectionId,
		LinearId:     issue.Id,
		Title:        issue.Title,
		Number:       issue.Number,
		BranchName:   issue.BranchName,
	}
	return linearIssue, nil
}
