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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/testmo/models"
)

var ExtractAutomationRunsMeta = plugin.SubTaskMeta{
	Name:             "extractAutomationRuns",
	EntryPoint:       ExtractAutomationRuns,
	EnabledByDefault: true,
	Description:      "Extract raw automation runs data into tool layer table testmo_automation_runs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

func ExtractAutomationRuns(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*TestmoTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TestmoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_AUTOMATION_RUN_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			var apiRun struct {
				Id                   uint64     `json:"id"`
				ProjectId            uint64     `json:"project_id"`
				SourceId             uint64     `json:"source_id"`
				Name                 string     `json:"name"`
				Status               int32      `json:"status"`
				ConfigId             uint64     `json:"config_id"`
				MilestoneId          uint64     `json:"milestone_id"`
				Elapsed              *int64     `json:"elapsed"`
				IsCompleted          bool       `json:"is_completed"`
				UntestedCount        uint64     `json:"untested_count"`
				Status1Count         uint64     `json:"status1_count"`
				Status2Count         uint64     `json:"status2_count"`
				Status3Count         uint64     `json:"status3_count"`
				Status4Count         uint64     `json:"status4_count"`
				Status5Count         uint64     `json:"status5_count"`
				Status6Count         uint64     `json:"status6_count"`
				Status7Count         uint64     `json:"status7_count"`
				Status8Count         uint64     `json:"status8_count"`
				Status9Count         uint64     `json:"status9_count"`
				Status10Count        uint64     `json:"status10_count"`
				Status11Count        uint64     `json:"status11_count"`
				Status12Count        uint64     `json:"status12_count"`
				Status13Count        uint64     `json:"status13_count"`
				Status14Count        uint64     `json:"status14_count"`
				Status15Count        uint64     `json:"status15_count"`
				Status16Count        uint64     `json:"status16_count"`
				Status17Count        uint64     `json:"status17_count"`
				Status18Count        uint64     `json:"status18_count"`
				Status19Count        uint64     `json:"status19_count"`
				Status20Count        uint64     `json:"status20_count"`
				Status21Count        uint64     `json:"status21_count"`
				Status22Count        uint64     `json:"status22_count"`
				Status23Count        uint64     `json:"status23_count"`
				Status24Count        uint64     `json:"status24_count"`
				SuccessCount         uint64     `json:"success_count"`
				FailureCount         uint64     `json:"failure_count"`
				CompletedCount       uint64     `json:"completed_count"`
				TotalCount           uint64     `json:"total_count"`
				ThreadCount          uint64     `json:"thread_count"`
				ThreadActiveCount    uint64     `json:"thread_active_count"`
				ThreadCompletedCount uint64     `json:"thread_completed_count"`
				CreatedAt            *time.Time `json:"created_at"`
				CreatedBy            uint64     `json:"created_by"`
				UpdatedAt            *time.Time `json:"updated_at"`
				UpdatedBy            *uint64    `json:"updated_by"`
				CompletedAt          *time.Time `json:"completed_at"`
				CompletedBy          *uint64    `json:"completed_by"`
			}

			err := json.Unmarshal(row.Data, &apiRun)
			if err != nil {
				return nil, errors.Default.Wrap(err, "error unmarshaling automation run")
			}

			automationRun := &models.TestmoAutomationRun{
				ConnectionId:         data.Options.ConnectionId,
				Id:                   apiRun.Id,
				ProjectId:            apiRun.ProjectId,
				SourceId:             apiRun.SourceId,
				Name:                 apiRun.Name,
				Status:               apiRun.Status,
				ConfigId:             apiRun.ConfigId,
				MilestoneId:          apiRun.MilestoneId,
				Elapsed:              apiRun.Elapsed,
				IsCompleted:          apiRun.IsCompleted,
				UntestedCount:        apiRun.UntestedCount,
				Status1Count:         apiRun.Status1Count,
				Status2Count:         apiRun.Status2Count,
				Status3Count:         apiRun.Status3Count,
				Status4Count:         apiRun.Status4Count,
				Status5Count:         apiRun.Status5Count,
				Status6Count:         apiRun.Status6Count,
				Status7Count:         apiRun.Status7Count,
				Status8Count:         apiRun.Status8Count,
				Status9Count:         apiRun.Status9Count,
				Status10Count:        apiRun.Status10Count,
				Status11Count:        apiRun.Status11Count,
				Status12Count:        apiRun.Status12Count,
				Status13Count:        apiRun.Status13Count,
				Status14Count:        apiRun.Status14Count,
				Status15Count:        apiRun.Status15Count,
				Status16Count:        apiRun.Status16Count,
				Status17Count:        apiRun.Status17Count,
				Status18Count:        apiRun.Status18Count,
				Status19Count:        apiRun.Status19Count,
				Status20Count:        apiRun.Status20Count,
				Status21Count:        apiRun.Status21Count,
				Status22Count:        apiRun.Status22Count,
				Status23Count:        apiRun.Status23Count,
				Status24Count:        apiRun.Status24Count,
				SuccessCount:         apiRun.SuccessCount,
				FailureCount:         apiRun.FailureCount,
				CompletedCount:       apiRun.CompletedCount,
				TotalCount:           apiRun.TotalCount,
				ThreadCount:          apiRun.ThreadCount,
				ThreadActiveCount:    apiRun.ThreadActiveCount,
				ThreadCompletedCount: apiRun.ThreadCompletedCount,
				TestmoCreatedAt:      apiRun.CreatedAt,
				CreatedBy:            apiRun.CreatedBy,
				TestmoUpdatedAt:      apiRun.UpdatedAt,
				UpdatedBy:            apiRun.UpdatedBy,
				CompletedAt:          apiRun.CompletedAt,
				CompletedBy:          apiRun.CompletedBy,
			}

			return []interface{}{automationRun}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
