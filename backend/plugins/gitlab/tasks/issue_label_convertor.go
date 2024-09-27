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
	RegisterSubtaskMeta(&ConvertIssueLabelsMeta)
}

var ConvertIssueLabelsMeta = plugin.SubTaskMeta{
	Name:             "Convert Issue Labels",
	EntryPoint:       ConvertIssueLabels,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_issue_labels into  domain layer table issue_labels",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	Dependencies:     []*plugin.SubTaskMeta{&ConvertIssuesMeta},
}

func ConvertIssueLabels(subtaskCtx plugin.SubTaskContext) errors.Error {
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(subtaskCtx, RAW_ISSUE_TABLE)

	db := subtaskCtx.GetDal()
	projectId := data.Options.ProjectId
	issueIdGen := didgen.NewDomainIdGenerator(&models.GitlabIssue{})

	converter, err := api.NewStatefulDataConverter(&api.StatefulDataConverterArgs[models.GitlabIssueLabel]{
		SubtaskCommonArgs: subtaskCommonArgs,
		Input: func(stateManager *api.SubtaskStateManager) (dal.Rows, errors.Error) {
			clauses := []dal.Clause{
				dal.Select("l.*"),
				dal.From("_tool_gitlab_issue_labels l"),
				dal.Join("LEFT JOIN _tool_gitlab_issues s ON s.gitlab_id = l.issue_id AND l.connection_id = s.connection_id"),
				dal.Where("s.project_id = ?  AND s.connection_id = ?", projectId, data.Options.ConnectionId),
				// dal.Orderby("s.issue_id ASC"),
			}
			if stateManager.IsIncremental() {
				since := stateManager.GetSince()
				if since != nil {
					clauses = append(clauses, dal.Where("l.updated_at >= ? ", since))
				}
			}
			return db.Cursor(clauses...)
		},
		Convert: func(issueLabel *models.GitlabIssueLabel) ([]interface{}, errors.Error) {
			domainIssueLabel := &ticket.IssueLabel{
				IssueId:   issueIdGen.Generate(data.Options.ConnectionId, issueLabel.IssueId),
				LabelName: issueLabel.LabelName,
			}
			return []interface{}{
				domainIssueLabel,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
