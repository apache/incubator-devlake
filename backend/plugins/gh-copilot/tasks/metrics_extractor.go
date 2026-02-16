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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
)

// Seat response structs (used by seat_extractor.go)

type copilotSeatResponse struct {
	CreatedAt               string          `json:"created_at"`
	UpdatedAt               string          `json:"updated_at"`
	PlanType                string          `json:"plan_type"`
	PendingCancellationDate *string         `json:"pending_cancellation_date"`
	LastAuthenticatedAt     *string         `json:"last_authenticated_at"`
	LastActivityAt          *string         `json:"last_activity_at"`
	LastActivityEditor      string          `json:"last_activity_editor"`
	Assignee                copilotAssignee `json:"assignee"`
}

type copilotAssignee struct {
	Login string `json:"login"`
	Id    int64  `json:"id"`
	Type  string `json:"type"`
}

// ExtractOrgMetrics parses org report data from the new report download API.
// The org report uses the same flat format as enterprise reports (day, totals_by_*).
// It extracts to GhCopilotOrgMetrics + GhCopilotLanguageMetrics tool-layer tables,
// aggregating the flat breakdown data into the org-level summary format.
func ExtractOrgMetrics(taskCtx plugin.SubTaskContext) errors.Error {
	data, ok := taskCtx.TaskContext().GetData().(*GhCopilotTaskData)
	if !ok {
		return errors.Default.New("task data is not GhCopilotTaskData")
	}
	connection := data.Connection
	connection.Normalize()

	if connection.Organization == "" {
		taskCtx.GetLogger().Info("No organization configured, skipping org metrics extraction")
		return nil
	}

	params := copilotRawParams{
		ConnectionId: data.Options.ConnectionId,
		ScopeId:      data.Options.ScopeId,
		Organization: connection.Organization,
		Endpoint:     connection.Endpoint,
	}

	// Derive seat counts from extracted seat assignments
	db := taskCtx.GetDal()
	seatTotal, err := db.Count(dal.From(&models.GhCopilotSeat{}), dal.Where(
		"connection_id = ? AND organization = ?",
		data.Options.ConnectionId,
		connection.Organization,
	))
	if err != nil {
		seatTotal = 0
	}
	seatActive, err := db.Count(dal.From(&models.GhCopilotSeat{}), dal.Where(
		"connection_id = ? AND organization = ? AND last_activity_at IS NOT NULL",
		data.Options.ConnectionId,
		connection.Organization,
	))
	if err != nil {
		seatActive = 0
	}

	// Org reports use the same flat format as enterprise reports (enterpriseDayTotal struct).
	// Each raw data row is one day's metrics with totals_by_* breakdowns.
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Table:   rawOrgMetricsTable,
			Options: params,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			var dt enterpriseDayTotal
			if err := errors.Convert(json.Unmarshal(row.Data, &dt)); err != nil {
				return nil, err
			}

			day, parseErr := time.Parse("2006-01-02", dt.Day)
			if parseErr != nil {
				return nil, errors.BadInput.Wrap(parseErr, "invalid day in org report")
			}

			// Aggregate code completions from feature breakdowns
			completionSuggestions := 0
			completionAcceptances := 0
			completionLinesSuggested := 0
			completionLinesAccepted := 0
			for _, f := range dt.TotalsByFeature {
				if f.Feature == "code_completion" {
					completionSuggestions += f.CodeGenerationActivityCount
					completionAcceptances += f.CodeAcceptanceActivityCount
					completionLinesSuggested += f.LocSuggestedToAddSum
					completionLinesAccepted += f.LocAddedSum
				}
			}

			// Aggregate chat metrics from feature breakdowns
			ideChats := 0
			for _, f := range dt.TotalsByFeature {
				if f.Feature == "chat_panel_ask_mode" || f.Feature == "chat_panel_agent_mode" ||
					f.Feature == "chat_panel_edit_mode" || f.Feature == "chat_panel_custom_mode" ||
					f.Feature == "chat_panel_unknown_mode" || f.Feature == "chat_inline" {
					ideChats += f.UserInitiatedInteractionCount
				}
			}

			var results []interface{}

			results = append(results, &models.GhCopilotOrgMetrics{
				ConnectionId:             data.Options.ConnectionId,
				ScopeId:                  data.Options.ScopeId,
				Date:                     day,
				TotalActiveUsers:         dt.DailyActiveUsers,
				TotalEngagedUsers:        0, // not available in flat format
				CompletionSuggestions:    completionSuggestions,
				CompletionAcceptances:    completionAcceptances,
				CompletionLinesSuggested: completionLinesSuggested,
				CompletionLinesAccepted:  completionLinesAccepted,
				IdeChats:                 ideChats,
				SeatTotal:                int(seatTotal),
				SeatActiveCount:          int(seatActive),
			})

			// Language metrics from totals_by_language_feature
			for _, lf := range dt.TotalsByLanguageFeature {
				if lf.Feature == "code_completion" {
					results = append(results, &models.GhCopilotLanguageMetrics{
						ConnectionId:   data.Options.ConnectionId,
						ScopeId:        data.Options.ScopeId,
						Date:           day,
						Editor:         "all",
						Language:       lf.Language,
						Suggestions:    lf.CodeGenerationActivityCount,
						Acceptances:    lf.CodeAcceptanceActivityCount,
						LinesSuggested: lf.LocSuggestedToAddSum,
						LinesAccepted:  lf.LocAddedSum,
					})
				}
			}

			return results, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}
