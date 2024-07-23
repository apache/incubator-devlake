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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"reflect"
)

func init() {
	RegisterSubtaskMeta(&ConvertWortItemsMeta)
}

var ConvertWortItemsMeta = plugin.SubTaskMeta{
	Name:             "convertApiWorkItems",
	EntryPoint:       ConvertApiWorkItems,
	EnabledByDefault: false,
	Description:      "Update domain layer ticket according to Azure DevOps Work Item",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	DependencyTables: []string{
		models.AzuredevopsWorkItem{}.TableName(),
	},
}

func ConvertApiWorkItems(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawWorkitemsTable)
	db := taskCtx.GetDal()

	//// select all work items belongs to the repository
	//clauses := []dal.Clause{
	//	dal.Select("aw.*"),
	//	dal.From("_tool_azuredevops_go_workitem aw"),
	//	dal.Join(`left join _tool_azuredevops_go_repo_work_items arw on (
	//		arw.work_item_id = aw.id
	//	)`),
	//	dal.Where("arw.repository_id = ? and arw.connection_id = ? ",
	//		data.Options.RepositoryId, data.Options.ConnectionId),
	//}

	clauses := []dal.Clause{
		dal.Select("*"),
		dal.From(models.AzuredevopsWorkItem{}),
		dal.Where("connection_id = ?", data.Options.ConnectionId)}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	//domainIdGenerator := didgen.NewDomainIdGenerator(&models.AzuredevopsWorkItem{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.AzuredevopsWorkItem{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			azureDevOpsWorkItem := inputRow.(*models.AzuredevopsWorkItem)

			workItem := &ticket.Issue{}
			//ConnectionId uint64 `gorm:"primaryKey"`
			//	//AzuredevopsId int    `gorm:"primaryKey"`
			//	WorkItemID   int `gorm:"primaryKey"`

			workItem.Title = azureDevOpsWorkItem.Title
			workItem.Type = azureDevOpsWorkItem.Type
			workItem.Status = azureDevOpsWorkItem.State
			workItem.CreatedDate = &azureDevOpsWorkItem.CreatedDate
			workItem.UpdatedDate = &azureDevOpsWorkItem.ChangedDate
			workItem.CreatorName = azureDevOpsWorkItem.CreatorName
			workItem.CreatorId = azureDevOpsWorkItem.CreatorId
			workItem.AssigneeName = azureDevOpsWorkItem.AssigneeName
			workItem.Status = azureDevOpsWorkItem.State

			return []interface{}{
				workItem,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
