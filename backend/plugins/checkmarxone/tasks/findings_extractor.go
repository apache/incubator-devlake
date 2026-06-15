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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	checkmarxoneModels "github.com/apache/incubator-devlake/plugins/checkmarxone/models"
)

var ExtractFindingsMeta = plugin.SubTaskMeta{
	Name:             "extractFindings",
	EntryPoint:       ExtractFindings,
	EnabledByDefault: true,
	Description:      "Extract findings data",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_SECURITY},
}

var ExtractFindingsMeta = plugin.SubTaskMeta{
	Name:             "extractFindings",
	EntryPoint:       ExtractFindings,
	EnabledByDefault: true,
	Description:      "Extract findings data",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_SECURITY},
}

func ExtractFindings(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*CheckmarxoneTaskData)
	logger := taskCtx.GetLogger()

	extractor, err := plugin.NewDataConverter(plugin.DataConverterArgs{
		RawDataSubTaskArgs: plugin.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Table:   RAW_FINDINGS_TABLE,
		},
		InputRowType: func() interface{} {
			return make(map[string]interface{})
		},
		Input: nil,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			rawData := inputRow.(map[string]interface{})

			finding := checkmarxoneModels.CheckmarxoneFinding{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			}

			if id, ok := rawData["id"].(string); ok {
				finding.FindingId = id
			}
			if name, ok := rawData["name"].(string); ok {
				finding.Name = name
			}
			if severity, ok := rawData["severity"].(string); ok {
				finding.Severity = severity
			}
			if status, ok := rawData["status"].(string); ok {
				finding.Status = status
			}
			if desc, ok := rawData["description"].(string); ok {
				finding.Description = desc
			}
			if state, ok := rawData["state"].(string); ok {
				finding.State = state
			}
			if fType, ok := rawData["type"].(string); ok {
				finding.Type = fType
			}
			if count, ok := rawData["count"].(float64); ok {
				finding.Count = int(count)
			}

			return []interface{}{&finding}, nil
		},
	})
	if err != nil {
		logger.Error(err, "failed to create converter")
		return err
	}

	return extractor.Execute()
}
