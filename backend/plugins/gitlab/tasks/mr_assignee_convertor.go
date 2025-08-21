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
	RegisterSubtaskMeta(&ConvertMrAssigneesMeta)
}

var ConvertMrAssigneesMeta = plugin.SubTaskMeta{
	Name:             "Convert MR Assignees",
	EntryPoint:       ConvertMrAssignees,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_assignees into  domain layer table pull_request_assignees",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	Dependencies:     []*plugin.SubTaskMeta{&ConvertApiMergeRequestsMeta},
}

func ConvertMrAssignees(subtaskCtx plugin.SubTaskContext) errors.Error {
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(subtaskCtx, RAW_MERGE_REQUEST_TABLE)
	db := subtaskCtx.GetDal()
	projectId := data.Options.ProjectId
	mrIdGen := didgen.NewDomainIdGenerator(&models.GitlabMergeRequest{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.GitlabAccount{})

	converter, err := api.NewStatefulDataConverter(&api.StatefulDataConverterArgs[models.GitlabAssignee]{
		SubtaskCommonArgs: subtaskCommonArgs,
		Input: func(stateManager *api.SubtaskStateManager) (dal.Rows, errors.Error) {
			clauses := []dal.Clause{
				dal.From(&models.GitlabAssignee{}),
				dal.Where(`project_id = ? and connection_id = ?`, projectId, data.Options.ConnectionId),
				dal.Orderby("merge_request_id ASC"),
			}
			if stateManager.IsIncremental() {
				since := stateManager.GetSince()
				if since != nil {
					clauses = append(clauses, dal.Where("updated_at >= ? ", since))
				}
			}
			return db.Cursor(clauses...)
		},
		Convert: func(mrAssignee *models.GitlabAssignee) ([]interface{}, errors.Error) {
			domainPrAssigne := &code.PullRequestAssignee{
				PullRequestId: mrIdGen.Generate(data.Options.ConnectionId, mrAssignee.MergeRequestId),
				AssigneeId:    accountIdGen.Generate(data.Options.ConnectionId, mrAssignee.AssigneeId),
				Name:          mrAssignee.Name,
				UserName:      mrAssignee.Username,
			}
			return []interface{}{
				domainPrAssigne,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
