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
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	githubModels "github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertIssueCommentsMeta = core.SubTaskMeta{
	Name:             "convertIssueComments",
	EntryPoint:       ConvertIssueComments,
	EnabledByDefault: true,
	Description:      "ConvertIssueComments data from Github api",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ConvertIssueComments(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Options.GithubId

	cursor, err := db.Cursor(
		dal.From(&githubModels.GithubIssueComment{}),
		dal.Join("left join _tool_github_issues "+
			"on _tool_github_issues.github_id = _tool_github_issue_comments.issue_id"),
		dal.Where("repo_id = ? and _tool_github_issues.connection_id = ?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubIssue{})
	accountIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubAccount{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(githubModels.GithubIssueComment{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_COMMENTS_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			githubIssueComment := inputRow.(*githubModels.GithubIssueComment)
			domainIssueComment := &ticket.IssueComment{
				DomainEntity: domainlayer.DomainEntity{
					Id: issueIdGen.Generate(data.Options.ConnectionId, githubIssueComment.GithubId),
				},
				IssueId:     issueIdGen.Generate(data.Options.ConnectionId, githubIssueComment.IssueId),
				Body:        githubIssueComment.Body,
				AccountId:   accountIdGen.Generate(data.Options.ConnectionId, githubIssueComment.AuthorUserId),
				CreatedDate: githubIssueComment.GithubCreatedAt,
			}
			return []interface{}{
				domainIssueComment,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
