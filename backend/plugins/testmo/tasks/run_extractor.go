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
	"regexp"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/testmo/models"
)

var ExtractRunsMeta = plugin.SubTaskMeta{
	Name:             "extractRuns",
	EntryPoint:       ExtractRuns,
	EnabledByDefault: true,
	Description:      "Extract raw runs data into tool layer table testmo_runs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

func ExtractRuns(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*TestmoTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TestmoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_RUN_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			var apiRun struct {
				Id          uint64     `json:"id"`
				ProjectId   uint64     `json:"project_id"`
				Name        string     `json:"name"`
				Description string     `json:"description"`
				Url         string     `json:"url"`
				StateId     int32      `json:"state_id"`
				MilestoneId *uint64    `json:"milestone_id"`
				ConfigIds   []uint64   `json:"config_ids"`
				IsClosed    bool       `json:"is_closed"`
				IsCompleted bool       `json:"is_completed"`
				Elapsed     *int64     `json:"elapsed"`
				CreatedAt   *time.Time `json:"created_at"`
				CreatedBy   *uint64    `json:"created_by"`
				UpdatedAt   *time.Time `json:"updated_at"`
				UpdatedBy   *uint64    `json:"updated_by"`
				ClosedAt    *time.Time `json:"closed_at"`
				ClosedBy    *uint64    `json:"closed_by"`
			}

			err := json.Unmarshal(row.Data, &apiRun)
			if err != nil {
				return nil, errors.Default.Wrap(err, "error unmarshaling run")
			}

			run := &models.TestmoRun{
				ConnectionId:    data.Options.ConnectionId,
				Id:              apiRun.Id,
				ProjectId:       apiRun.ProjectId,
				Name:            apiRun.Name,
				Status:          apiRun.StateId,
				StatusName:      getStatusName(apiRun.StateId),
				Elapsed:         apiRun.Elapsed,
				Message:         apiRun.Description,
				TestmoCreatedAt: apiRun.CreatedAt,
				TestmoUpdatedAt: apiRun.UpdatedAt,
			}

			// If TestmoUpdatedAt is null, use TestmoCreatedAt as fallback
			if run.TestmoUpdatedAt == nil && run.TestmoCreatedAt != nil {
				run.TestmoUpdatedAt = run.TestmoCreatedAt
			}

			// Classify run types based on scope config patterns
			if data.Options.ScopeConfig != nil {
				run.IsAcceptanceTest = matchesPattern(run.Name, data.Options.ScopeConfig.AcceptanceTestPattern)
				run.IsSmokeTest = matchesPattern(run.Name, data.Options.ScopeConfig.SmokeTestPattern)
				run.Team = extractTeam(run.Name, data.Options.ScopeConfig.TeamPattern)
			}

			return []interface{}{run}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func matchesPattern(runName, pattern string) bool {
	if pattern == "" {
		return false
	}

	return strings.Contains(strings.ToLower(runName), strings.ToLower(pattern))
}

func extractTeam(runName, pattern string) string {
	if pattern == "" {
		return ""
	}

	if regex, err := regexp.Compile(pattern); err == nil {
		matches := regex.FindStringSubmatch(runName)
		if len(matches) > 1 {
			return matches[1] // Return first capture group
		}
	}

	return ""
}

func getStatusName(stateId int32) string {
	// Map Testmo state IDs to human-readable names
	switch stateId {
	case 1:
		return "Passed"
	case 2:
		return "Failed"
	case 3:
		return "Blocked"
	case 4:
		return "Skipped"
	case 5:
		return "Retest"
	case 6:
		return "In Progress"
	case 7:
		return "Active"
	default:
		return "Unknown"
	}
}
