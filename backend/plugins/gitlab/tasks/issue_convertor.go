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
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"reflect"
	"strconv"
	"strings"
)

var ConvertIssuesMeta = plugin.SubTaskMeta{
	Name:             "convertIssues",
	EntryPoint:       ConvertIssues,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_issues into  domain layer table issues",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertIssues(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUE_TABLE)
	projectId := data.Options.ProjectId

	clauses := []dal.Clause{
		dal.Select("issues.*"),
		dal.From("_tool_gitlab_issues issues"),
		dal.Where("project_id = ? and connection_id = ?", projectId, data.Options.ConnectionId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&models.GitlabIssue{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.GitlabAccount{})
	boardIdGen := didgen.NewDomainIdGenerator(&models.GitlabProject{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabIssue{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			issue := inputRow.(*models.GitlabIssue)
			domainIssue := &ticket.Issue{
				DomainEntity:            domainlayer.DomainEntity{Id: issueIdGen.Generate(data.Options.ConnectionId, issue.GitlabId)},
				IssueKey:                strconv.Itoa(issue.Number),
				Title:                   issue.Title,
				Description:             issue.Body,
				Priority:                issue.Priority,
				Type:                    issue.StdType,
				OriginalType:            issue.Type,
				LeadTimeMinutes:         int64(issue.LeadTimeMinutes),
				Url:                     issue.Url,
				CreatedDate:             &issue.GitlabCreatedAt,
				UpdatedDate:             &issue.GitlabUpdatedAt,
				ResolutionDate:          issue.ClosedAt,
				Severity:                issue.Severity,
				Component:               issue.Component,
				OriginalStatus:          issue.Status,
				OriginalEstimateMinutes: issue.TimeEstimate,
				TimeSpentMinutes:        issue.TotalTimeSpent,
				CreatorId:               accountIdGen.Generate(data.Options.ConnectionId, issue.CreatorId),
				CreatorName:             issue.CreatorName,
				AssigneeId:              accountIdGen.Generate(data.Options.ConnectionId, issue.AssigneeId),
				AssigneeName:            issue.AssigneeName,
			}

			if strings.ToUpper(issue.State) == "OPENED" {
				domainIssue.Status = ticket.TODO
			} else {
				domainIssue.Status = ticket.DONE
			}

			boardIssue := &ticket.BoardIssue{
				BoardId: boardIdGen.Generate(data.Options.ConnectionId, projectId),
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
