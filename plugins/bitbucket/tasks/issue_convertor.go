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
	"strconv"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	bitbucketModels "github.com/apache/incubator-devlake/plugins/bitbucket/models"
)

var ConvertIssuesMeta = core.SubTaskMeta{
	Name:             "convertIssues",
	EntryPoint:       ConvertIssues,
	EnabledByDefault: true,
	Description:      "Convert tool layer table bitbucket_issues into  domain layer table issues",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ConvertIssues(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*BitbucketTaskData)
	repoId := data.Repo.BitbucketId

	issue := &bitbucketModels.BitbucketIssue{}
	cursor, err := db.Cursor(
		dal.From(issue),
		dal.Where("repo_id = ? and connection_id=?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&bitbucketModels.BitbucketIssue{})
	boardIdGen := didgen.NewDomainIdGenerator(&bitbucketModels.BitbucketRepo{})
	accountIdGen := didgen.NewDomainIdGenerator(&bitbucketModels.BitbucketAccount{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: BitbucketApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_ISSUE_TABLE,
		},
		InputRowType: reflect.TypeOf(bitbucketModels.BitbucketIssue{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			issue := inputRow.(*bitbucketModels.BitbucketIssue)
			domainIssue := &ticket.Issue{
				DomainEntity:    domainlayer.DomainEntity{Id: issueIdGen.Generate(data.Options.ConnectionId, issue.BitbucketId)},
				IssueKey:        strconv.Itoa(issue.Number),
				Title:           issue.Title,
				Description:     issue.Body,
				Priority:        issue.Priority,
				Type:            issue.Type,
				OriginalStatus:  issue.State,
				LeadTimeMinutes: issue.LeadTimeMinutes,
				Url:             issue.Url,
				CreatedDate:     &issue.BitbucketCreatedAt,
				UpdatedDate:     &issue.BitbucketUpdatedAt,
				ResolutionDate:  issue.ClosedAt,
				Severity:        issue.Severity,
				Component:       issue.Component,
			}
			if issue.State == "closed" {
				domainIssue.Status = ticket.DONE
			} else {
				domainIssue.Status = ticket.TODO
			}
			if issue.AssigneeName != "" {
				domainIssue.AssigneeName = issue.AssigneeName
				domainIssue.AssigneeId = accountIdGen.Generate(data.Options.ConnectionId, issue.AssigneeId)
			}
			if issue.AuthorName != "" {
				domainIssue.CreatorName = issue.AuthorName
				domainIssue.CreatorId = accountIdGen.Generate(data.Options.ConnectionId, issue.AuthorId)
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
