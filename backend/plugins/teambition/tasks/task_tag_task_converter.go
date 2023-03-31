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
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/teambition/models"
	"reflect"
)

var ConvertTaskTagTasksMeta = plugin.SubTaskMeta{
	Name:             "convertTaskTags",
	EntryPoint:       ConvertTaskTagTasks,
	EnabledByDefault: true,
	Description:      "convert teambition task tags",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertTaskTagTasks(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TASK_TAG_TABLE)
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	logger.Info("convert project:%v task tag tasks", data.Options.ProjectId)
	clauses := []dal.Clause{
		dal.Select("b.name as name, a.task_id as task_id, a.connection_id as connection_id, a.project_id as project_id"),
		dal.From("_tool_teambition_task_tag_tasks a"),
		dal.Join(`left join _tool_teambition_task_tags b on (
			a.connection_id = b.connection_id
			AND a.project_id = b.project_id
			AND a.task_tag_id = b.id
		)`),
		dal.Where("a.connection_id = ? AND a.project_id = ?", data.Options.ConnectionId, data.Options.ProjectId),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.TeambitionTaskTagTask{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			userTool := inputRow.(*models.TeambitionTaskTagTask)
			issue := &ticket.IssueLabel{
				IssueId:   getTaskIdGen().Generate(userTool.ConnectionId, userTool.TaskId),
				LabelName: userTool.Name,
			}
			return []interface{}{
				issue,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
