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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiBuildRecordsMeta)
}

var ExtractApiBuildRecordsMeta = plugin.SubTaskMeta{
	Name:             "extractApiTimelineRecords",
	EntryPoint:       ExtractApiTimelineTasks,
	EnabledByDefault: true,
	Description:      "Extract raw timeline record data into tool layer table",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	DependencyTables: []string{RawTimelineRecordTable},
	ProductTables: []string{
		models.AzuredevopsTimelineRecord{}.TableName(),
	},
}

func truncateString(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength]
	}
	return s
}

func ExtractApiTimelineTasks(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawTimelineRecordTable)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			results := make([]interface{}, 0, 1)

			recordApi := &models.AzuredevopsApiTimelineRecord{}
			err := errors.Convert(json.Unmarshal(row.Data, recordApi))
			if err != nil {
				return nil, err
			}

			input := &SimplePr{}
			err = errors.Convert(json.Unmarshal(row.Input, &input))
			if err != nil {
				return nil, err
			}

			record := &models.AzuredevopsTimelineRecord{
				ConnectionId: data.Options.ConnectionId,
				RecordId:     recordApi.Id,
				BuildId:      input.AzuredevopsId,
				ParentId:     recordApi.ParentId,
				Type:         "",
				Name:         truncateString(recordApi.Name, 255),
				StartTime:    recordApi.StartTime,
				FinishTime:   recordApi.FinishTime,
				State:        recordApi.State,
				Result:       recordApi.Result,
				ChangeId:     recordApi.ChangeId,
				LastModified: recordApi.LastModified,
			}

			results = append(results, record)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
