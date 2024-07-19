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
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
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
	Dependencies:     []*plugin.SubTaskMeta{&ExtractApiIssuesMeta},
}

func ConvertIssueAssignee(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUE_TABLE)

	cursor, err := db.Cursor(
		dal.From(&models.GitlabIssueAssignee{}),
		dal.Where("project_id = ? and connection_id=?", data.Options.ProjectId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&models.GitlabIssue{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.GitlabAccount{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabIssueAssignee{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			input := inputRow.(*models.GitlabIssueAssignee)
			domainIssueAssignee := &ticket.IssueAssignee{
				IssueId:      issueIdGen.Generate(data.Options.ConnectionId, input.GitlabId),
				AssigneeId:   accountIdGen.Generate(data.Options.ConnectionId, input.AssigneeId),
				AssigneeName: input.AssigneeName,
			}
			return []interface{}{domainIssueAssignee}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
