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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/testmo/models"
)

var ExtractProjectsMeta = plugin.SubTaskMeta{
	Name:             "extractProjects",
	EntryPoint:       ExtractProjects,
	EnabledByDefault: true,
	Description:      "Extract raw projects data into tool layer table testmo_projects",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

func ExtractProjects(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*TestmoTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TestmoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_PROJECT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			var apiProject struct {
				Id                           uint64 `json:"id"`
				Name                         string `json:"name"`
				IsCompleted                  bool   `json:"is_completed"`
				MilestoneCount               uint64 `json:"milestone_count"`
				MilestoneActiveCount         uint64 `json:"milestone_active_count"`
				MilestoneCompletedCount      uint64 `json:"milestone_completed_count"`
				RunCount                     uint64 `json:"run_count"`
				RunActiveCount               uint64 `json:"run_active_count"`
				RunClosedCount               uint64 `json:"run_closed_count"`
				AutomationSourceCount        uint64 `json:"automation_source_count"`
				AutomationSourceActiveCount  uint64 `json:"automation_source_active_count"`
				AutomationSourceRetiredCount uint64 `json:"automation_source_retired_count"`
				AutomationRunCount           uint64 `json:"automation_run_count"`
				AutomationRunActiveCount     uint64 `json:"automation_run_active_count"`
				AutomationRunCompletedCount  uint64 `json:"automation_run_completed_count"`
			}

			err := json.Unmarshal(row.Data, &apiProject)
			if err != nil {
				return nil, errors.Default.Wrap(err, "error unmarshaling project")
			}

			project := &models.TestmoProject{
				Scope: common.Scope{
					ConnectionId: data.Options.ConnectionId,
				},
				Id:                           apiProject.Id,
				Name:                         apiProject.Name,
				IsCompleted:                  apiProject.IsCompleted,
				MilestoneCount:               apiProject.MilestoneCount,
				MilestoneActiveCount:         apiProject.MilestoneActiveCount,
				MilestoneCompletedCount:      apiProject.MilestoneCompletedCount,
				RunCount:                     apiProject.RunCount,
				RunActiveCount:               apiProject.RunActiveCount,
				RunClosedCount:               apiProject.RunClosedCount,
				AutomationSourceCount:        apiProject.AutomationSourceCount,
				AutomationSourceActiveCount:  apiProject.AutomationSourceActiveCount,
				AutomationSourceRetiredCount: apiProject.AutomationSourceRetiredCount,
				AutomationRunCount:           apiProject.AutomationRunCount,
				AutomationRunActiveCount:     apiProject.AutomationRunActiveCount,
				AutomationRunCompletedCount:  apiProject.AutomationRunCompletedCount,
			}

			return []interface{}{project}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
