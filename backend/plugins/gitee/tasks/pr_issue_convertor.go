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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitee/models"
	"github.com/spf13/cast"
	"reflect"
	"strconv"
)

var ConvertPullRequestIssuesMeta = plugin.SubTaskMeta{
	Name:             "convertPullRequestIssues",
	EntryPoint:       ConvertPullRequestIssues,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitee_pull_request_issues into  domain layer table pull_request_issues",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func ConvertPullRequestIssues(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_TABLE)
	repoId := data.Repo.GiteeId

	cursor, err := db.Cursor(
		dal.From(&models.GiteePullRequestIssue{}),
		dal.Join(`left join _tool_gitee_pull_requests on _tool_gitee_pull_requests.gitee_id = _tool_gitee_pull_request_issues.pull_request_id`),
		dal.Where("_tool_gitee_pull_requests.repo_id = ? and _tool_gitee_pull_requests.connection_id = ?", repoId, data.Options.ConnectionId),
		dal.Orderby("pull_request_id ASC"),
	)

	if err != nil {
		return err
	}
	defer cursor.Close()
	prIdGen := didgen.NewDomainIdGenerator(&models.GiteePullRequest{})
	issueIdGen := didgen.NewDomainIdGenerator(&models.GiteeIssue{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.GiteePullRequestIssue{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			giteePrIssue := inputRow.(*models.GiteePullRequestIssue)
			issueNum, _ := strconv.Atoi(giteePrIssue.IssueNumber)
			pullRequestIssue := &crossdomain.PullRequestIssue{
				PullRequestId:  prIdGen.Generate(data.Options.ConnectionId, giteePrIssue.PullRequestId),
				IssueId:        issueIdGen.Generate(data.Options.ConnectionId, giteePrIssue.IssueId),
				IssueKey:       cast.ToString(issueNum),
				PullRequestKey: giteePrIssue.PullRequestNumber,
			}
			return []interface{}{
				pullRequestIssue,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
