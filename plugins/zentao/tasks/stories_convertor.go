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
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"reflect"
)

var _ core.SubTaskEntryPoint = ConvertExecutions

var ConvertStoriesMeta = core.SubTaskMeta{
	Name:             "convertStories",
	EntryPoint:       ConvertStories,
	EnabledByDefault: true,
	Description:      "convert Zentao stories",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ConvertStories(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	db := taskCtx.GetDal()
	boardIdGen := didgen.NewDomainIdGenerator(&models.ZentaoStories{})
	cursor, err := db.Cursor(
		dal.From(&models.ZentaoStories{}),
		dal.Where(`_tool_zentao_stories.execution_id = ? and 
			_tool_zentao_stories.connection_id = ?`, data.Options.ExecutionId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	convertor, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.ZentaoStories{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ZentaoApiParams{
				ProductId:   data.Options.ProductId,
				ExecutionId: data.Options.ExecutionId,
				ProjectId:   data.Options.ProjectId,
			},
			Table: RAW_STORIES_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			toolStories := inputRow.(*models.ZentaoStories)

			domainBoard := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: boardIdGen.Generate(toolStories.ConnectionId, toolStories.ID),
				},
				Url:      "",
				IconURL:  "",
				IssueKey: "",
				Title:    toolStories.Title,
				//Description:             toolStories.Spec,
				EpicKey:                 "",
				Type:                    toolStories.Type,
				Status:                  toolStories.Status,
				OriginalStatus:          "",
				StoryPoint:              0,
				ResolutionDate:          nil,
				CreatedDate:             toolStories.OpenedDate.ToNullableTime(),
				UpdatedDate:             toolStories.LastEditedDate.ToNullableTime(),
				LeadTimeMinutes:         0,
				ParentIssueId:           "",
				Priority:                "",
				OriginalEstimateMinutes: 0,
				TimeSpentMinutes:        0,
				TimeRemainingMinutes:    0,
				CreatorId:               "",
				CreatorName:             toolStories.OpenedBy.OpenedByRealname,
				AssigneeId:              "",
				AssigneeName:            toolStories.AssignedTo.Realname,
				Severity:                "",
				Component:               "",
				DeploymentId:            "",
			}

			results := make([]interface{}, 0)
			results = append(results, domainBoard)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return convertor.Execute()
}
