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
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"reflect"
)

func init() {
	RegisterSubtaskMeta(&ConvertApiPrCommitsMeta)
}

var ConvertApiPrCommitsMeta = plugin.SubTaskMeta{
	Name:             "convertApiPullRequestsCommits",
	EntryPoint:       ConvertApiPullRequestsCommits,
	EnabledByDefault: true,
	Description:      "Add domain layer PullRequestCommit according to Azure DevOps Pull Request Commit",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	DependencyTables: []string{
		models.AzuredevopsPrCommit{}.TableName(),
	},
}

func ConvertApiPullRequestsCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawPrCommitTable)
	db := taskCtx.GetDal()

	clauses := []dal.Clause{
		dal.From(&models.AzuredevopsPrCommit{}),
		dal.Join(`left join _tool_azuredevops_go_pull_requests
			on _tool_azuredevops_go_pull_requests.azuredevops_id =
			_tool_azuredevops_go_pull_request_commits.pull_request_id`),
		dal.Where(`_tool_azuredevops_go_pull_requests.repository_id = ?
			and _tool_azuredevops_go_pull_requests.connection_id = ?`,
			data.Options.RepositoryId, data.Options.ConnectionId),
		dal.Orderby("pull_request_id ASC"),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}

	domainIdGenerator := didgen.NewDomainIdGenerator(&models.AzuredevopsPullRequest{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.AzuredevopsPrCommit{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			azuredevopsPrCommit := inputRow.(*models.AzuredevopsPrCommit)
			domainPrcommit := &code.PullRequestCommit{
				CommitSha:          azuredevopsPrCommit.CommitSha,
				PullRequestId:      domainIdGenerator.Generate(data.Options.ConnectionId, azuredevopsPrCommit.PullRequestId),
				CommitAuthorName:   azuredevopsPrCommit.AuthorName,
				CommitAuthorEmail:  azuredevopsPrCommit.AuthorEmail,
				CommitAuthoredDate: *azuredevopsPrCommit.AuthorDate,
			}
			return []interface{}{
				domainPrcommit,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
