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
	"github.com/apache/incubator-devlake/errors"
	"reflect"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

func ConvertWorklog(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_WORKLOG_TABLE, false)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	logger.Info("convert board:%d", data.Options.WorkspaceId)
	worklogIdGen := didgen.NewDomainIdGenerator(&models.TapdWorklog{})
	clauses := []dal.Clause{
		dal.From(&models.TapdWorklog{}),
		dal.Where("connection_id = ? AND workspace_id = ?", data.Options.ConnectionId, data.Options.WorkspaceId),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	taskIdGen := didgen.NewDomainIdGenerator(&models.TapdTask{})
	storyIdGen := didgen.NewDomainIdGenerator(&models.TapdStory{})
	bugIdGen := didgen.NewDomainIdGenerator(&models.TapdBug{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.TapdWorklog{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			toolL := inputRow.(*models.TapdWorklog)
			domainL := &ticket.IssueWorklog{
				DomainEntity: domainlayer.DomainEntity{
					Id: worklogIdGen.Generate(data.Options.ConnectionId, toolL.Id),
				},
				AuthorId:         getAccountIdGen().Generate(data.Options.ConnectionId, toolL.Owner),
				Comment:          toolL.Memo,
				TimeSpentMinutes: int(toolL.Timespent),
				LoggedDate:       (*time.Time)(toolL.Created),
			}
			switch strings.ToUpper(toolL.EntityType) {
			case "TASK":
				domainL.IssueId = taskIdGen.Generate(toolL.EntityId)
			case "BUG":
				domainL.IssueId = bugIdGen.Generate(toolL.EntityId)
			case "STORY":
				domainL.IssueId = storyIdGen.Generate(toolL.EntityId)
			}
			return []interface{}{
				domainL,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

var ConvertWorklogMeta = core.SubTaskMeta{
	Name:             "convertWorklog",
	EntryPoint:       ConvertWorklog,
	EnabledByDefault: true,
	Description:      "convert Tapd Worklog",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}
