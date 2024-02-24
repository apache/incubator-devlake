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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"reflect"
)

func init() {
	RegisterSubtaskMeta(&ConvertCommitsMeta)
}

var ConvertCommitsMeta = plugin.SubTaskMeta{
	Name:             "convertApiCommits",
	EntryPoint:       ConvertApiCommits,
	EnabledByDefault: false,
	Description:      "Update domain layer commit according to Azure DevOps Commit",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
	DependencyTables: []string{
		models.AzuredevopsCommit{}.TableName(),
	},
}

func ConvertApiCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawCommitTable)
	db := taskCtx.GetDal()

	// select all commits belongs to the repository
	clauses := []dal.Clause{
		dal.Select("ac.*"),
		dal.From("_tool_azuredevops_go_commits ac"),
		dal.Join(`left join _tool_azuredevops_go_repo_commits arc on (
			arc.commit_sha = ac.sha
		)`),
		dal.Where("arc.repository_id = ? and arc.connection_id = ? ",
			data.Options.RepositoryId, data.Options.ConnectionId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.AzuredevopsCommit{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			azuredevopsCommit := inputRow.(*models.AzuredevopsCommit)

			commit := &code.Commit{}
			commit.Sha = azuredevopsCommit.Sha
			commit.Message = azuredevopsCommit.Message
			commit.Additions = azuredevopsCommit.Additions
			commit.Deletions = azuredevopsCommit.Deletions
			commit.AuthorId = azuredevopsCommit.AuthorEmail
			commit.AuthorName = azuredevopsCommit.AuthorName
			commit.AuthorEmail = azuredevopsCommit.AuthorEmail
			commit.AuthoredDate = *azuredevopsCommit.AuthoredDate
			commit.CommitterName = azuredevopsCommit.CommitterName
			commit.CommitterEmail = azuredevopsCommit.CommitterEmail
			commit.CommittedDate = *azuredevopsCommit.CommittedDate
			commit.CommitterId = azuredevopsCommit.CommitterEmail

			// convert repo / commits relationship
			repoCommit := &code.RepoCommit{
				RepoId:    didgen.NewDomainIdGenerator(&models.AzuredevopsRepo{}).Generate(data.Options.ConnectionId, data.Options.RepositoryId),
				CommitSha: azuredevopsCommit.Sha,
			}

			return []interface{}{
				commit,
				repoCommit,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
