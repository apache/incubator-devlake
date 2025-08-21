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
	RegisterSubtaskMeta(&ConvertMrReviewersMeta)
}

var ConvertMrReviewersMeta = plugin.SubTaskMeta{
	Name:             "Convert MR Reviewers",
	EntryPoint:       ConvertMrReviewers,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_reviewers into domain layer table pull_request_reviewers",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	Dependencies:     []*plugin.SubTaskMeta{&ConvertApiMergeRequestsMeta},
}

func ConvertMrReviewers(subtaskCtx plugin.SubTaskContext) errors.Error {
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(subtaskCtx, RAW_MERGE_REQUEST_TABLE)
	db := subtaskCtx.GetDal()

	mrIdGen := didgen.NewDomainIdGenerator(&models.GitlabMergeRequest{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.GitlabAccount{})

	converter, err := api.NewStatefulDataConverter(&api.StatefulDataConverterArgs[models.GitlabReviewer]{
		SubtaskCommonArgs: subtaskCommonArgs,
		Input: func(stateManager *api.SubtaskStateManager) (dal.Rows, errors.Error) {
			clauses := []dal.Clause{
				dal.Select("c.*"),
				dal.From("_tool_gitlab_reviewers c"),
				dal.Join(`LEFT JOIN _tool_gitlab_merge_requests mr ON mr.gitlab_id = c.merge_request_id AND c.connection_id = mr.connection_id`),
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
		Convert: func(mrReviewer *models.GitlabReviewer) ([]interface{}, errors.Error) {
			domainPrReviewer := &code.PullRequestReviewer{
				PullRequestId: mrIdGen.Generate(data.Options.ConnectionId, mrReviewer.MergeRequestId),
				ReviewerId:    accountIdGen.Generate(data.Options.ConnectionId, mrReviewer.ReviewerId),
				Name:          mrReviewer.Name,
				UserName:      mrReviewer.Username,
			}
			return []interface{}{
				domainPrReviewer,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
