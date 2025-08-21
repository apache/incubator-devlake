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
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"reflect"
)

func init() {
	RegisterSubtaskMeta(&ConvertApiPullRequestsMeta)
}

var ConvertApiPullRequestsMeta = plugin.SubTaskMeta{
	Name:             "convertApiPullRequests",
	EntryPoint:       ConvertApiPullRequests,
	EnabledByDefault: true,
	Description:      "Add domain layer PullRequest according to Azure DevOps Pull Requests",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	DependencyTables: []string{
		models.AzuredevopsPullRequest{}.TableName(),
	},
}

func ConvertApiPullRequests(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawPullRequestTable)
	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.From(&models.AzuredevopsPullRequest{}),
		dal.Where("repository_id=? and connection_id = ?", data.Options.RepositoryId, data.Options.ConnectionId),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	domainPrIdGenerator := didgen.NewDomainIdGenerator(&models.AzuredevopsPullRequest{})
	domainRepoIdGenerator := didgen.NewDomainIdGenerator(&models.AzuredevopsRepo{})
	domainUserIdGen := didgen.NewDomainIdGenerator(&models.AzuredevopsUser{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.AzuredevopsPullRequest{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			azuredevopsPr := inputRow.(*models.AzuredevopsPullRequest)

			domainPr := &code.PullRequest{
				DomainEntity: domainlayer.DomainEntity{
					Id: domainPrIdGenerator.Generate(data.Options.ConnectionId, azuredevopsPr.AzuredevopsId),
				},
				HeadRepoId:     domainRepoIdGenerator.Generate(data.Options.ConnectionId, azuredevopsPr.RepositoryId),
				BaseRepoId:     domainRepoIdGenerator.Generate(data.Options.ConnectionId, azuredevopsPr.RepositoryId),
				OriginalStatus: azuredevopsPr.Status,
				PullRequestKey: azuredevopsPr.AzuredevopsId,
				Title:          azuredevopsPr.Title,
				Description:    azuredevopsPr.Description,
				Type:           azuredevopsPr.Type,
				Url:            azuredevopsPr.Url,
				AuthorName:     azuredevopsPr.CreatedByName,
				AuthorId:       domainUserIdGen.Generate(data.Options.ConnectionId, azuredevopsPr.CreatedById),
				CreatedDate:    *azuredevopsPr.CreationDate,
				MergedDate:     azuredevopsPr.ClosedDate,
				ClosedDate:     azuredevopsPr.ClosedDate,
				MergeCommitSha: azuredevopsPr.MergeCommitSha,
				HeadRef:        azuredevopsPr.SourceRefName,
				BaseRef:        azuredevopsPr.TargetRefName,
				Component:      "", // not supported
				BaseCommitSha:  azuredevopsPr.TargetCommitSha,
				HeadCommitSha:  azuredevopsPr.SourceCommitSha,
			}
			switch azuredevopsPr.Status {
			case "opened":
				domainPr.Status = code.OPEN
			case "merged":
				domainPr.Status = code.MERGED
			case "closed", "locked":
				domainPr.Status = code.CLOSED
			default:
				domainPr.Status = azuredevopsPr.Status
			}

			return []interface{}{
				domainPr,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
