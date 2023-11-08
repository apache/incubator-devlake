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
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&ConvertPullRequestCommitsMeta)
}

var ConvertPullRequestCommitsMeta = plugin.SubTaskMeta{
	Name:             "convertPullRequestCommits",
	EntryPoint:       ConvertPullRequestCommits,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_pull_request_commits into  domain layer table pull_request_commits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS, plugin.DOMAIN_TYPE_CODE_REVIEW},
	DependencyTables: []string{
		models.GithubPrCommit{}.TableName(),    // cursor
		models.GithubPullRequest{}.TableName(), // cursor
		models.GithubPullRequest{}.TableName(), // id generator
		RAW_PR_COMMIT_TABLE},
	ProductTables: []string{code.PullRequestCommit{}.TableName()},
}

func ConvertPullRequestCommits(taskCtx plugin.SubTaskContext) (err errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Options.GithubId

	pullIdGen := didgen.NewDomainIdGenerator(&models.GithubPullRequest{})

	cursor, err := db.Cursor(
		dal.From(&models.GithubPrCommit{}),
		dal.Join(`left join _tool_github_pull_requests on _tool_github_pull_requests.github_id = _tool_github_pull_request_commits.pull_request_id`),
		dal.Where("_tool_github_pull_requests.repo_id = ? and _tool_github_pull_requests.connection_id = ?", repoId, data.Options.ConnectionId),
		dal.Orderby("pull_request_id ASC"),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.GithubPrCommit{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_PR_COMMIT_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			githubPullRequestCommit := inputRow.(*models.GithubPrCommit)
			domainPrCommit := &code.PullRequestCommit{
				CommitSha:          githubPullRequestCommit.CommitSha,
				PullRequestId:      pullIdGen.Generate(data.Options.ConnectionId, githubPullRequestCommit.PullRequestId),
				CommitAuthorName:   githubPullRequestCommit.CommitAuthorName,
				CommitAuthorEmail:  githubPullRequestCommit.CommitAuthorEmail,
				CommitAuthoredDate: githubPullRequestCommit.CommitAuthoredDate,
			}
			return []interface{}{
				domainPrCommit,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
