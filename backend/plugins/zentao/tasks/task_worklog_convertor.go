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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var ConvertTaskWorklogsMeta = plugin.SubTaskMeta{
	Name:             "convertTaskWorklogs",
	EntryPoint:       ConvertTaskWorklogs,
	EnabledByDefault: true,
	Description:      "convert Zentao task worklogs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertTaskWorklogs(taskCtx plugin.SubTaskContext) errors.Error {
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*ZentaoTaskData)
	logger.Info(
		"convert Zentao task worklogs of %d in %d",
		data.Options.ProjectId,
		data.Options.ConnectionId,
	)
	worklogIdGen := didgen.NewDomainIdGenerator(&models.ZentaoWorklog{})
	clauses := []dal.Clause{
		dal.From(&models.ZentaoWorklog{}),
		dal.Where(
			"connection_id = ? AND project = ? AND object_type = ?",
			data.Options.ConnectionId,
			data.Options.ProjectId,
			"task",
		),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	taskIdGen := didgen.NewDomainIdGenerator(&models.ZentaoTask{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_TASK_WORKLOGS_TABLE,
		},
		InputRowType: reflect.TypeOf(models.ZentaoWorklog{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			toolL := inputRow.(*models.ZentaoWorklog)
			domainL := &ticket.IssueWorklog{
				DomainEntity: domainlayer.DomainEntity{
					Id: worklogIdGen.Generate(data.Options.ConnectionId, toolL.Id),
				},
				Comment:          toolL.Work,
				TimeSpentMinutes: int(toolL.Consumed * 60),
			}
			timeData, err := common.ConvertStringToTime(toolL.Date)
			if err != nil {
				return nil, errors.Default.Wrap(err, "failed to convert zentao task worklog date")
			}
			// zentao task only has one field as date type for worklog creation
			domainL.StartedDate = &timeData
			domainL.LoggedDate = &timeData

			domainL.IssueId = taskIdGen.Generate(data.Options.ConnectionId, toolL.ObjectId)

			// get ID of account by username
			var account models.ZentaoAccount
			err = db.First(&account, dal.Where("connection_id = ? AND account = ?",
				data.Options.ConnectionId, toolL.Account))
			if err != nil {
				// if account isn't available, giving empty string as ID
				if db.IsErrorNotFound(err) {
					logger.Warn(nil, "cannot find zentao account by account: %s", toolL.Account)
					domainL.AuthorId = ""
				} else {
					return nil, errors.Default.Wrap(err, "failed to get zentao account by account")
				}
			} else {
				accountIdGen := didgen.NewDomainIdGenerator(&models.ZentaoAccount{})
				domainL.AuthorId = accountIdGen.Generate(account.ConnectionId, account.ID)
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
