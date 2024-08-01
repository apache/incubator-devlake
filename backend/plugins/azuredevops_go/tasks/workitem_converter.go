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
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
)

func init() {
	RegisterSubtaskMeta(&ConvertWortItemsMeta)
}

var ConvertWortItemsMeta = plugin.SubTaskMeta{
	Name:             "convertApiWorkItems",
	EntryPoint:       ConvertApiWorkItems,
	EnabledByDefault: true,
	Description:      "Update domain layer ticket according to Azure DevOps Work Item",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	DependencyTables: []string{
		models.AzuredevopsWorkItem{}.TableName(),
	},
}

func ConvertApiWorkItems(taskCtx plugin.SubTaskContext) errors.Error {
	_, data := CreateRawDataSubTaskArgs(taskCtx, RawWorkitemsTable)
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()

	var existingWorkItem models.AzuredevopsWorkItem
	var existingWorkItems []models.AzuredevopsWorkItem
	err := db.All(&existingWorkItems, dal.Where("connection_id = ?", data.Options.ConnectionId))
	if err != nil {
		return err
	}

	logger.Debug("Total number of work items: #%s", len(existingWorkItems))

	for _, existingWorkItem = range existingWorkItems {
		finalIssue := &ticket.Issue{}

		finalIssue.Id = existingWorkItem.WorkItemID
		finalIssue.Component = existingWorkItem.Area
		finalIssue.Title = existingWorkItem.Title
		finalIssue.Type = existingWorkItem.Type
		finalIssue.Status = existingWorkItem.State
		finalIssue.CreatedDate = existingWorkItem.CreatedDate
		finalIssue.UpdatedDate = existingWorkItem.ChangedDate
		finalIssue.CreatorName = existingWorkItem.CreatorName
		finalIssue.CreatorId = existingWorkItem.CreatorId
		finalIssue.AssigneeName = existingWorkItem.AssigneeName
		finalIssue.Status = existingWorkItem.State
		finalIssue.Url = existingWorkItem.Url
		finalIssue.StoryPoint = &existingWorkItem.StoryPoint
		finalIssue.Severity = existingWorkItem.Severity
		finalIssue.Priority = existingWorkItem.Priority

		if existingWorkItem.ResolvedDate != nil {
			finalIssue.ResolutionDate = existingWorkItem.ResolvedDate
			temp := uint(existingWorkItem.ResolvedDate.Sub(*existingWorkItem.CreatedDate).Minutes())
			finalIssue.LeadTimeMinutes = &temp
		}

		err := db.CreateOrUpdate(finalIssue, dal.Where("connection_id = ? AND id = ?", data.Options.ConnectionId, finalIssue.Id))
		if err != nil {
			return err
		}
	}

	return nil
}
