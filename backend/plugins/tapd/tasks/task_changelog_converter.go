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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"reflect"
	"time"
)

type TaskChangelogItemResult struct {
	ConnectionId      uint64     `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Id                uint64     `gorm:"primaryKey;type:BIGINT  NOT NULL" json:"id"`
	WorkspaceId       uint64     `json:"workspace_id"`
	WorkitemTypeId    uint64     `json:"workitem_type_id"`
	Creator           string     `json:"creator"`
	Created           *time.Time `json:"created"`
	ChangeSummary     string     `json:"change_summary"`
	Comment           string     `json:"comment"`
	EntityType        string     `json:"entity_type"`
	ChangeType        string     `json:"change_type"`
	ChangeTypeText    string     `json:"change_type_text"`
	TaskId            uint64     `json:"task_id"`
	ChangelogId       uint64     `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Field             string     `json:"field" gorm:"primaryKey;type:varchar(255)"`
	ValueBeforeParsed string     `json:"value_before"`
	ValueAfterParsed  string     `json:"value_after"`
	IterationIdFrom   uint64
	IterationIdTo     uint64
	common.NoPKModel
}

func ConvertTaskChangelog(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TASK_CHANGELOG_TABLE)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	customStatusMap := getStatusMapping(data)
	getTaskStdStatus := func(statusKey string) string {
		if statusKey == "done" {
			return ticket.DONE
		} else if statusKey == "progressing" {
			return ticket.IN_PROGRESS
		} else {
			return ticket.TODO
		}
	}
	logger.Info("convert changelog :%d", data.Options.WorkspaceId)
	clIdGen := didgen.NewDomainIdGenerator(&models.TapdTaskChangelog{})
	issueIdGen := didgen.NewDomainIdGenerator(&models.TapdTask{})
	clIterIdGen := didgen.NewDomainIdGenerator(&models.TapdIteration{})
	clauses := []dal.Clause{
		dal.Select("tc.created, tc.id, tc.workspace_id, tc.task_id, tc.creator, _tool_tapd_task_changelog_items.*"),
		dal.From(&models.TapdTaskChangelogItem{}),
		dal.Join("left join _tool_tapd_task_changelogs tc on tc.id = _tool_tapd_task_changelog_items.changelog_id "),
		dal.Where("tc.connection_id = ? AND tc.workspace_id = ?", data.Options.ConnectionId, data.Options.WorkspaceId),
		dal.Orderby("created DESC"),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(TaskChangelogItemResult{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			cl := inputRow.(*TaskChangelogItemResult)
			domainCl := &ticket.IssueChangelogs{
				DomainEntity: domainlayer.DomainEntity{
					Id: fmt.Sprintf("%s:%s", clIdGen.Generate(data.Options.ConnectionId, cl.Id), cl.Field),
				},
				IssueId:           issueIdGen.Generate(data.Options.ConnectionId, cl.TaskId),
				AuthorId:          getAccountIdGen().Generate(data.Options.ConnectionId, cl.Creator),
				AuthorName:        cl.Creator,
				FieldId:           cl.Field,
				FieldName:         cl.Field,
				OriginalFromValue: cl.ValueBeforeParsed,
				OriginalToValue:   cl.ValueAfterParsed,
				CreatedDate:       *cl.Created,
			}
			if domainCl.FieldName == "status" {
				if len(customStatusMap) != 0 {
					domainCl.FromValue = customStatusMap[domainCl.OriginalFromValue]
					domainCl.ToValue = customStatusMap[domainCl.OriginalToValue]
				} else {
					domainCl.FromValue = getTaskStdStatus(domainCl.OriginalFromValue)
					domainCl.ToValue = getTaskStdStatus(domainCl.OriginalToValue)
				}
			}
			if domainCl.FieldName == "iteration_id" {
				domainCl.FieldName = "Sprint"
				if cl.IterationIdFrom == 0 {
					domainCl.OriginalFromValue = ""
					domainCl.FromValue = ""
				} else {
					domainCl.OriginalFromValue = clIterIdGen.Generate(cl.ConnectionId, cl.IterationIdFrom)
					domainCl.FromValue = cl.ValueBeforeParsed
				}
				if cl.IterationIdTo == 0 {
					domainCl.OriginalToValue = ""
					domainCl.ToValue = ""
				} else {
					domainCl.OriginalToValue = clIterIdGen.Generate(cl.ConnectionId, cl.IterationIdTo)
					domainCl.ToValue = cl.ValueAfterParsed
				}
			}
			if domainCl.FieldId == "owner" {
				domainCl.FieldName = "assignee"
				domainCl.OriginalFromValue = generateDomainAccountIdForUsers(cl.ValueBeforeParsed, cl.ConnectionId)
				domainCl.OriginalToValue = generateDomainAccountIdForUsers(cl.ValueAfterParsed, cl.ConnectionId)
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

var ConvertTaskChangelogMeta = plugin.SubTaskMeta{
	Name:             "convertTaskChangelog",
	EntryPoint:       ConvertTaskChangelog,
	EnabledByDefault: true,
	Description:      "convert Tapd task changelog",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}
