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
	"strconv"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ConvertIssuesMeta)
}

var ConvertIssuesMeta = plugin.SubTaskMeta{
	Name:             "Convert Issues",
	EntryPoint:       ConvertIssues,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_issues into  domain layer table issues",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	Dependencies:     []*plugin.SubTaskMeta{&ConvertApiMrCommitsMeta},
}

func ConvertIssues(subtaskCtx plugin.SubTaskContext) errors.Error {
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(subtaskCtx, RAW_ISSUE_TABLE)

	db := subtaskCtx.GetDal()
	projectId := data.Options.ProjectId
	issueIdGen := didgen.NewDomainIdGenerator(&models.GitlabIssue{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.GitlabAccount{})
	boardIdGen := didgen.NewDomainIdGenerator(&models.GitlabProject{})

	converter, err := api.NewStatefulDataConverter(&api.StatefulDataConverterArgs[models.GitlabIssue]{
		SubtaskCommonArgs: subtaskCommonArgs,
		Input: func(stateManager *api.SubtaskStateManager) (dal.Rows, errors.Error) {
			clauses := []dal.Clause{
				dal.From(&models.GitlabIssue{}),
				dal.Where("connection_id = ? AND project_id = ?", data.Options.ConnectionId, data.Options.ProjectId),
			}
			if stateManager.IsIncremental() {
				since := stateManager.GetSince()
				if since != nil {
					clauses = append(clauses, dal.Where("updated_at >= ? ", since))
				}
			}
			return db.Cursor(clauses...)
		},
		BeforeConvert: func(issue *models.GitlabIssue, stateManager *api.SubtaskStateManager) errors.Error {
			issueId := issueIdGen.Generate(data.Options.ConnectionId, issue.GitlabId)
			if err := db.Delete(&ticket.IssueLabel{}, dal.Where("issue_id = ?", issueId)); err != nil {
				return err
			}
			if err := db.Delete(&ticket.IssueAssignee{}, dal.Where("issue_id = ?", issueId)); err != nil {
				return err
			}
			return nil
		},
		Convert: func(issue *models.GitlabIssue) ([]interface{}, errors.Error) {
			domainIssue := &ticket.Issue{
				DomainEntity:            domainlayer.DomainEntity{Id: issueIdGen.Generate(data.Options.ConnectionId, issue.GitlabId)},
				IssueKey:                strconv.Itoa(issue.Number),
				Title:                   issue.Title,
				Description:             issue.Body,
				Priority:                issue.Priority,
				OriginalType:            issue.Type,
				LeadTimeMinutes:         issue.LeadTimeMinutes,
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
			if strings.ToUpper(issue.Type) == ticket.INCIDENT {
				domainIssue.Type = ticket.INCIDENT
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
