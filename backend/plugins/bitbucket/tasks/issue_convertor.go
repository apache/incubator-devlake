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
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"reflect"
	"strconv"
)

var ConvertIssuesMeta = plugin.SubTaskMeta{
	Name:             "convertIssues",
	EntryPoint:       ConvertIssues,
	EnabledByDefault: true,
	Description:      "Convert tool layer table bitbucket_issues into  domain layer table issues",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertIssues(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUE_COMMENTS_TABLE)
	db := taskCtx.GetDal()
	repoId := data.Options.FullName

	issue := &models.BitbucketIssue{}
	cursor, err := db.Cursor(
		dal.From(issue),
		dal.Where("repo_id = ? and connection_id=?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&models.BitbucketIssue{})
	boardIdGen := didgen.NewDomainIdGenerator(&models.BitbucketRepo{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.BitbucketAccount{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.BitbucketIssue{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			issue := inputRow.(*models.BitbucketIssue)
			domainIssue := &ticket.Issue{
				DomainEntity:    domainlayer.DomainEntity{Id: issueIdGen.Generate(data.Options.ConnectionId, issue.RepoId, issue.BitbucketId)},
				IssueKey:        strconv.Itoa(issue.Number),
				Title:           issue.Title,
				Description:     issue.Body,
				Priority:        issue.Priority,
				Type:            issue.Type,
				Status:          issue.StdState,
				OriginalStatus:  issue.State,
				LeadTimeMinutes: int64(issue.LeadTimeMinutes),
				Url:             issue.Url,
				CreatedDate:     &issue.BitbucketCreatedAt,
				UpdatedDate:     &issue.BitbucketUpdatedAt,
				ResolutionDate:  issue.ClosedAt,
				Severity:        issue.Severity,
				Component:       issue.Component,
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
