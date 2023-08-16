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

	"github.com/apache/incubator-devlake/plugins/jira/models"
)

var ConvertIssueRelationshipsMeta = plugin.SubTaskMeta{
	Name:             "convertIssueRelationships",
	EntryPoint:       ConvertIssueRelationships,
	EnabledByDefault: true,
	Description:      "Convert tool layer table jira_issue_relationships into  domain layer table issue_relationships",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertIssueRelationships(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*JiraTaskData)

	cursor, err := db.Cursor(
		dal.Select("jir.*"),
		dal.From("_tool_jira_issue_relationships jir"),
		dal.Join(`LEFT JOIN _tool_jira_board_issues jbi
              ON jir.connection_id = jbi.connection_id AND jir.issue_id = jbi.issue_id`),
		dal.Where("jir.connection_id = ? AND jbi.board_id = ?", data.Options.ConnectionId, data.Options.BoardId),
		dal.Orderby("issue_id ASC"),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&models.JiraIssue{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_ISSUE_TABLE,
		},
		InputRowType: reflect.TypeOf(models.JiraIssueRelationship{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			issueRelationship := inputRow.(*models.JiraIssueRelationship)
			domainIssueRelationship := &ticket.IssueRelationship{
				SourceIssueId: issueIdGen.Generate(issueRelationship.ConnectionId, issueRelationship.IssueId),
			}
			if issueRelationship.InwardIssueId != 0 {
				domainIssueRelationship.TargetIssueId = issueIdGen.Generate(issueRelationship.ConnectionId, issueRelationship.InwardIssueId)
				domainIssueRelationship.OriginalType = issueRelationship.Inward
			} else {
				domainIssueRelationship.TargetIssueId = issueIdGen.Generate(issueRelationship.ConnectionId, issueRelationship.OutwardIssueId)
				domainIssueRelationship.OriginalType = issueRelationship.Outward
			}

			return []interface{}{
				domainIssueRelationship,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
