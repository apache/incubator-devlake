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

//import (
//	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
//	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
//	"github.com/apache/incubator-devlake/plugins/core"
//	"github.com/apache/incubator-devlake/plugins/helper"
//	"github.com/apache/incubator-devlake/plugins/tapd/models"
//	"reflect"
//)
//
//func ConvertIssueSprintsHistory(taskCtx core.SubTaskContext) error {
//	data := taskCtx.GetData().(*TapdTaskData)
//	logger := taskCtx.GetLogger()
//	db := taskCtx.GetDb()
//	logger.Info("convert board:%d", data.Options.WorkspaceID)
//	iterIdGen := didgen.NewDomainIdGenerator(&models.TapdIteration{})
//	cursor, err := db.Model(&models.TapdIssueSprintHistory{}).Where("connection_id = ? AND workspace_id = ?", data.Connection.ID, data.Options.WorkspaceID).Rows()
//	if err != nil {
//		return err
//	}
//	defer cursor.Close()
//	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
//		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
//			Ctx: taskCtx,
//			Params: TapdApiParams{
//				ConnectionId: data.Connection.ID,
//
//				WorkspaceID: data.Options.WorkspaceID,
//			},
//			Table: "tapd_api_%",
//		},
//		InputRowType: reflect.TypeOf(models.TapdIssueSprintHistory{}),
//		Input:        cursor,
//		Convert: func(inputRow interface{}) ([]interface{}, error) {
//			toolL := inputRow.(*models.TapdIssueSprintHistory)
//			domainL := &ticket.IssueSprintsHistory{
//				IssueId:   IssueIdGen.Generate(data.Connection.ID, toolL.IssueId),
//				SprintId:  iterIdGen.Generate(data.Connection.ID, toolL.SprintId),
//				StartDate: toolL.StartDate,
//				EndDate:   &toolL.EndDate,
//			}
//			return []interface{}{
//				domainL,
//			}, nil
//		},
//	})
//	if err != nil {
//		return err
//	}
//
//	return converter.Execute()
//}
//
//var ConvertIssueSprintsHistoryMeta = core.SubTaskMeta{
//	Name:             "convertIssueSprintsHistory",
//	EntryPoint:       ConvertIssueSprintsHistory,
//	EnabledByDefault: true,
//	Description:      "convert Tapd IssueSprintsHistory",
//}
