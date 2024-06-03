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
	RegisterSubtaskMeta(&ConvertMrAssigneesMeta)
}

var ConvertMrAssigneesMeta = plugin.SubTaskMeta{
	Name:             "Convert MR Assignees",
	EntryPoint:       ConvertMrAssignees,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_assignees into  domain layer table pull_request_assignees",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	Dependencies:     []*plugin.SubTaskMeta{&ExtractApiMergeRequestDetailsMeta},
}

func ConvertMrAssignees(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_TABLE)
	projectId := data.Options.ProjectId
	clauses := []dal.Clause{
		dal.Select("_tool_gitlab_assignees.*"),
		dal.From(&models.GitlabAssignee{}),
		dal.Join(`left join _tool_gitlab_merge_requests mr on
			mr.gitlab_id = _tool_gitlab_assignees.merge_request_id`),
		dal.Where(`mr.project_id = ?
			and mr.connection_id = ?`,
			projectId, data.Options.ConnectionId),
		dal.Orderby("_tool_gitlab_assignees.merge_request_id ASC"),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	mrIdGen := didgen.NewDomainIdGenerator(&models.GitlabMergeRequest{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabAssignee{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			mrAssignee := inputRow.(*models.GitlabAssignee)
			domainPrAssigne := &code.PullRequestAssignee{
				PullRequestId: mrIdGen.Generate(data.Options.ConnectionId, mrAssignee.MergeRequestId),
				AssigneeId:    mrAssignee.AssigneeId,
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
