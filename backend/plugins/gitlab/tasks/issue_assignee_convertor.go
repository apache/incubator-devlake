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
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ConvertIssueAssigneeMeta)
}

var ConvertIssueAssigneeMeta = plugin.SubTaskMeta{
	Name:             "convert Issue Assignees",
	EntryPoint:       ConvertIssueAssignee,
	EnabledByDefault: true,
	Description:      "Convert tool layer table _tool_gitlab_issue_assignees into  domain layer table issue_assignees",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	DependencyTables: []string{models.GitlabIssueAssignee{}.TableName()},
	Dependencies:     []*plugin.SubTaskMeta{&ConvertIssuesMeta},
}

func ConvertIssueAssignee(subtaskCtx plugin.SubTaskContext) errors.Error {
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(subtaskCtx, RAW_ISSUE_TABLE)
	db := subtaskCtx.GetDal()

	issueIdGen := didgen.NewDomainIdGenerator(&models.GitlabIssue{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.GitlabAccount{})

	converter, err := api.NewStatefulDataConverter(&api.StatefulDataConverterArgs[models.GitlabIssueAssignee]{
		SubtaskCommonArgs: subtaskCommonArgs,
		Input: func(stateManager *api.SubtaskStateManager) (dal.Rows, errors.Error) {
			clauses := []dal.Clause{
				dal.From(&models.GitlabIssueAssignee{}),
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
		Convert: func(issueAssignee *models.GitlabIssueAssignee) ([]interface{}, errors.Error) {
			domainIssueAssignee := &ticket.IssueAssignee{
				IssueId:      issueIdGen.Generate(data.Options.ConnectionId, issueAssignee.GitlabId),
				AssigneeId:   accountIdGen.Generate(data.Options.ConnectionId, issueAssignee.AssigneeId),
				AssigneeName: issueAssignee.AssigneeName,
			}
			return []interface{}{domainIssueAssignee}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
