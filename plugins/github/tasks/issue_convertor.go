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
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"reflect"
	"strconv"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	githubModels "github.com/apache/incubator-devlake/plugins/github/models"
)

var ConvertIssuesMeta = core.SubTaskMeta{
	Name:             "convertIssues",
	EntryPoint:       ConvertIssues,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_issues into  domain layer table issues",
}

func ConvertIssues(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	issue := &githubModels.GithubIssue{}
	cursor, err := db.Cursor(
		dal.From(issue),
		dal.Where("repo_id = ? and connection_id=?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubIssue{})
	userIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubUser{})
	boardIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubRepo{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_ISSUE_TABLE,
		},
		InputRowType: reflect.TypeOf(githubModels.GithubIssue{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			issue := inputRow.(*githubModels.GithubIssue)
			domainIssue := &ticket.Issue{
				DomainEntity:    domainlayer.DomainEntity{Id: issueIdGen.Generate(data.Options.ConnectionId, issue.GithubId)},
				IssueKey:        strconv.Itoa(issue.Number),
				Title:           issue.Title,
				Description:     issue.Body,
				Priority:        issue.Priority,
				Type:            issue.Type,
				AssigneeId:      userIdGen.Generate(issue.AssigneeId),
				AssigneeName:    issue.AssigneeName,
				CreatorId:       userIdGen.Generate(issue.AuthorId),
				CreatorName:     issue.AuthorName,
				LeadTimeMinutes: issue.LeadTimeMinutes,
				Url:             issue.Url,
				CreatedDate:     &issue.GithubCreatedAt,
				UpdatedDate:     &issue.GithubUpdatedAt,
				ResolutionDate:  issue.ClosedAt,
				Severity:        issue.Severity,
				Component:       issue.Component,
			}
			if issue.State == "closed" {
				domainIssue.Status = ticket.DONE
			} else {
				domainIssue.Status = ticket.TODO
			}
			boardIssue := &ticket.BoardIssue{
				BoardId: boardIdGen.Generate(data.Options.ConnectionId, repoId),
				IssueId: domainIssue.Id,
			}
			return []interface{}{
				domainIssue,
				boardIssue,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
