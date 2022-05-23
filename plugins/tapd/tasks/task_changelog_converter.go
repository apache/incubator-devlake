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
	"fmt"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"reflect"
	"time"
)

type TaskChangelogItemResult struct {
	ConnectionId      uint64    `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID                uint64    `gorm:"primaryKey;type:BIGINT  NOT NULL" json:"id"`
	WorkspaceID       uint64    `json:"workspace_id"`
	WorkitemTypeID    uint64    `json:"workitem_type_id"`
	Creator           string    `json:"creator"`
	Created           time.Time `json:"created"`
	ChangeSummary     string    `json:"change_summary"`
	Comment           string    `json:"comment"`
	EntityType        string    `json:"entity_type"`
	ChangeType        string    `json:"change_type"`
	ChangeTypeText    string    `json:"change_type_text"`
	TaskID            uint64    `json:"task_id"`
	ChangelogId       uint64    `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Field             string    `json:"field" gorm:"primaryKey;type:varchar(255)"`
	ValueBeforeParsed string    `json:"value_before"`
	ValueAfterParsed  string    `json:"value_after"`
	IterationIdFrom   uint64
	IterationIdTo     uint64
	common.NoPKModel
}

func ConvertTaskChangelog(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert changelog :%d", data.Options.WorkspaceID)
	clIdGen := didgen.NewDomainIdGenerator(&models.TapdTaskChangelog{})

	cursor, err := db.Table("_tool_tapd_task_changelog_items").
		Joins("left join _tool_tapd_task_changelogs tc on tc.id = _tool_tapd_task_changelog_items.changelog_id ").
		Where("tc.connection_id = ? AND tc.workspace_id = ?", data.Connection.ID, data.Options.WorkspaceID).
		Select("tc.*, _tool_tapd_task_changelog_items.*").
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,

				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_TASK_CHANGELOG_TABLE,
		},
		InputRowType: reflect.TypeOf(TaskChangelogItemResult{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			cl := inputRow.(*TaskChangelogItemResult)
			domainCl := &ticket.Changelog{
				DomainEntity: domainlayer.DomainEntity{
					Id: fmt.Sprintf("%s:%s", clIdGen.Generate(data.Connection.ID, cl.ID), cl.Field),
				},
				IssueId:     IssueIdGen.Generate(data.Connection.ID, cl.TaskID),
				AuthorId:    UserIdGen.Generate(data.Connection.ID, data.Options.WorkspaceID, cl.Creator),
				AuthorName:  cl.Creator,
				FieldId:     cl.Field,
				FieldName:   cl.Field,
				From:        cl.ValueBeforeParsed,
				To:          cl.ValueAfterParsed,
				CreatedDate: cl.Created,
			}

			return []interface{}{
				domainCl,
			}, nil
		},
	})
	if err != nil {
		logger.Info(err.Error())
		return err
	}

	return converter.Execute()
}

var ConvertTaskChangelogMeta = core.SubTaskMeta{
	Name:             "convertTaskChangelog",
	EntryPoint:       ConvertTaskChangelog,
	EnabledByDefault: true,
	Description:      "convert Tapd task changelog",
}
