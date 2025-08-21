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
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ConvertMrLabelsMeta)
}

var ConvertMrLabelsMeta = plugin.SubTaskMeta{
	Name:             "Convert MR Labels",
	EntryPoint:       ConvertMrLabels,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_mr_labels into  domain layer table pull_request_labels",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	Dependencies:     []*plugin.SubTaskMeta{&ConvertApiMergeRequestsMeta},
}

func ConvertMrLabels(subtaskCtx plugin.SubTaskContext) errors.Error {
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(subtaskCtx, RAW_MERGE_REQUEST_TABLE)
	db := subtaskCtx.GetDal()

	mrIdGen := didgen.NewDomainIdGenerator(&models.GitlabMergeRequest{})

	converter, err := api.NewStatefulDataConverter(&api.StatefulDataConverterArgs[models.GitlabMrLabel]{
		SubtaskCommonArgs: subtaskCommonArgs,
		Input: func(stateManager *api.SubtaskStateManager) (dal.Rows, errors.Error) {
			clauses := []dal.Clause{
				dal.Select("c.*"),
				dal.From("_tool_gitlab_mr_labels c"),
				dal.Join(`LEFT JOIN _tool_gitlab_merge_requests mr ON mr.gitlab_id = c.mr_id AND c.connection_id = mr.connection_id`),
				dal.Where(`mr.project_id = ?  and mr.connection_id = ?`, data.Options.ProjectId, data.Options.ConnectionId),
			}
			if stateManager.IsIncremental() {
				since := stateManager.GetSince()
				if since != nil {
					clauses = append(clauses, dal.Where("c.updated_at >= ? ", since))
				}
			}
			return db.Cursor(clauses...)
		},
		Convert: func(mrLabel *models.GitlabMrLabel) ([]interface{}, errors.Error) {
			domainIssueLabel := &code.PullRequestLabel{
				PullRequestId: mrIdGen.Generate(data.Options.ConnectionId, mrLabel.MrId),
				LabelName:     mrLabel.LabelName,
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
