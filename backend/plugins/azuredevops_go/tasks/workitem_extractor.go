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
	"encoding/json"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiWorkItemsMeta)
}

var ExtractApiWorkItemsMeta = plugin.SubTaskMeta{
	Name:             "extractApiWorkItems",
	EntryPoint:       ExtractApiWorkItems,
	EnabledByDefault: false,
	Description:      "Extract raw work items data into tool layer table",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	DependencyTables: []string{RawWorkitemsTable},
	ProductTables: []string{
		models.AzuredevopsWorkItem{}.TableName()},
}

func ExtractApiWorkItems(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawWorkitemsTable)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			results := make([]interface{}, 0, 2)

			apiWorkItem := &models.AzuredevopsApiWorkItem{}
			err := errors.Convert(json.Unmarshal(row.Data, apiWorkItem))
			if err != nil {
				return nil, err
			}

			// create project/commits relationship
			repoWorkItem := &models.AzuredevopsWorkItem{
				ConnectionId: data.Options.ConnectionId,
				WorkItemID:   apiWorkItem.Id,
				Title:        apiWorkItem.Fields.SystemTitle,
				Type:         apiWorkItem.Fields.SystemWorkItemType,
				State:        apiWorkItem.Fields.SystemState,
				CreatedDate:  *common.Iso8601TimeToTime(apiWorkItem.Fields.SystemCreatedDate),
				ChangedDate:  *common.Iso8601TimeToTime(apiWorkItem.Fields.SystemChangedDate),
				CreatorId:    apiWorkItem.Fields.SystemCreatedBy.Id,
				CreatorName:  apiWorkItem.Fields.SystemCreatedBy.DisplayName,
				AssigneeName: apiWorkItem.Fields.SystemAssignedTo,
			}

			results = append(results, repoWorkItem)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
