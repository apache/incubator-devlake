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
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"reflect"
)

var _ plugin.SubTaskEntryPoint = ConvertExecutions

var ConvertExecutionMeta = plugin.SubTaskMeta{
	Name:             "convertExecutions",
	EntryPoint:       ConvertExecutions,
	EnabledByDefault: true,
	Description:      "convert Zentao executions",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertExecutions(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	db := taskCtx.GetDal()
	executionIdGen := didgen.NewDomainIdGenerator(&models.ZentaoExecution{})
	projectIdGen := didgen.NewDomainIdGenerator(&models.ZentaoProject{})
	cursor, err := db.Cursor(
		dal.From(&models.ZentaoExecution{}),
		dal.Where(`project_id = ? and connection_id = ?`, data.Options.ProjectId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	convertor, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.ZentaoExecution{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ZentaoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProductId:    data.Options.ProductId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_EXECUTION_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			toolExecution := inputRow.(*models.ZentaoExecution)

			domainStatus := ``
			switch toolExecution.Status {
			case `wait`:
				domainStatus = `FUTURE`
			case `doing`:
				domainStatus = `ACTIVE`
			case `suspended`:
				domainStatus = `SUSPENDED`
			case `closed`:
			case `done`:
				domainStatus = `CLOSED`
			}

			sprint := &ticket.Sprint{
				DomainEntity: domainlayer.DomainEntity{
					Id: executionIdGen.Generate(toolExecution.ConnectionId, toolExecution.Id),
				},
				Name:            toolExecution.Name,
				Url:             toolExecution.Path,
				Status:          domainStatus,
				StartedDate:     toolExecution.RealBegan.ToNullableTime(),
				EndedDate:       toolExecution.RealEnd.ToNullableTime(),
				CompletedDate:   toolExecution.ClosedDate.ToNullableTime(),
				OriginalBoardID: projectIdGen.Generate(toolExecution.ConnectionId, toolExecution.Id),
			}
			boardSprint := &ticket.BoardSprint{
				BoardId:  projectIdGen.Generate(toolExecution.ConnectionId, toolExecution.Id),
				SprintId: sprint.Id,
			}
			results := make([]interface{}, 0)
			results = append(results, sprint, boardSprint)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return convertor.Execute()
}
