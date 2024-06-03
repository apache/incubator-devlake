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
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
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
	Dependencies:     []*plugin.SubTaskMeta{&ExtractApiMergeRequestDetailsMeta},
}

func ConvertMrReviewers(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_TABLE)
	projectId := data.Options.ProjectId
	clauses := []dal.Clause{
		dal.Select("_tool_gitlab_reviewers.*"),
		dal.From(&models.GitlabReviewer{}),
		dal.Join(`left join _tool_gitlab_merge_requests mr on
			mr.gitlab_id = _tool_gitlab_reviewers.merge_request_id`),
		dal.Where(`mr.project_id = ?
			and mr.connection_id = ?`,
			projectId, data.Options.ConnectionId),
		dal.Orderby("_tool_gitlab_reviewers.merge_request_id ASC"),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	mrIdGen := didgen.NewDomainIdGenerator(&models.GitlabMergeRequest{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabReviewer{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			mrReviewer := inputRow.(*models.GitlabReviewer)
			domainPrReviewer := &code.PullRequestReviewer{
				PullRequestId: mrIdGen.Generate(data.Options.ConnectionId, mrReviewer.MergeRequestId),
				ReviewerId:    mrReviewer.ReviewerId,
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
