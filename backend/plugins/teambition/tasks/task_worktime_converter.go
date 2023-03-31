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
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/teambition/models"
	"reflect"
)

var ConvertTaskWorktimeMeta = plugin.SubTaskMeta{
	Name:             "convertTaskWorktime",
	EntryPoint:       ConvertTaskWorktime,
	EnabledByDefault: true,
	Description:      "convert teambition task worktime",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertTaskWorktime(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TASK_WORKTIME_TABLE)
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	logger.Info("convert project:%v task worktime", data.Options.ProjectId)
	clauses := []dal.Clause{
		dal.From(&models.TeambitionTaskWorktime{}),
		dal.Where("connection_id = ? AND project_id = ?", data.Options.ConnectionId, data.Options.ProjectId),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.TeambitionTaskWorktime{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			userTool := inputRow.(*models.TeambitionTaskWorktime)
			issueWorklog := &ticket.IssueWorklog{
				DomainEntity: domainlayer.DomainEntity{
					Id: getTaskWorktimeIdGen().Generate(data.Options.ConnectionId, userTool.WorktimeId),
				},
				IssueId:          getTaskIdGen().Generate(userTool.ConnectionId, userTool.TaskId),
				AuthorId:         getAccountIdGen().Generate(userTool.ConnectionId, userTool.UserId),
				LoggedDate:       userTool.CreatedAt.ToNullableTime(),
				Comment:          userTool.Description,
				TimeSpentMinutes: int(userTool.Worktime / (60 * 1000)),
				StartedDate:      userTool.Date.ToNullableTime(),
			}
			return []interface{}{
				issueWorklog,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
