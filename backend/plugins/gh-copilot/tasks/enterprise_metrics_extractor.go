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
	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
)

// --- Enterprise report JSON structures ---

type enterpriseReport struct {
	ReportStartDay string                  `json:"report_start_day"`
	ReportEndDay   string                  `json:"report_end_day"`
	EnterpriseId   string                  `json:"enterprise_id"`
	DayTotals      []enterpriseDayTotal    `json:"day_totals"`
}

type enterpriseDayTotal struct {
	Day                           string                  `json:"day"`
	EnterpriseId                  string                  `json:"enterprise_id"`
	DailyActiveUsers              int                     `json:"daily_active_users"`
	WeeklyActiveUsers             int                     `json:"weekly_active_users"`
	MonthlyActiveUsers            int                     `json:"monthly_active_users"`
	MonthlyActiveChatUsers        int                     `json:"monthly_active_chat_users"`
	MonthlyActiveAgentUsers       int                     `json:"monthly_active_agent_users"`
	UserInitiatedInteractionCount int                     `json:"user_initiated_interaction_count"`
	CodeGenerationActivityCount   int                     `json:"code_generation_activity_count"`
	CodeAcceptanceActivityCount   int                     `json:"code_acceptance_activity_count"`
	LocSuggestedToAddSum          int                     `json:"loc_suggested_to_add_sum"`
	LocSuggestedToDeleteSum       int                     `json:"loc_suggested_to_delete_sum"`
	LocAddedSum                   int                     `json:"loc_added_sum"`
	LocDeletedSum                 int                     `json:"loc_deleted_sum"`
	TotalsByIde                   []totalsByIde           `json:"totals_by_ide"`
	TotalsByFeature               []totalsByFeature       `json:"totals_by_feature"`
	TotalsByLanguageFeature       []totalsByLangFeature   `json:"totals_by_language_feature"`
	TotalsByLanguageModel         []totalsByLangModel     `json:"totals_by_language_model"`
	TotalsByModelFeature          []totalsByModelFeature  `json:"totals_by_model_feature"`
	PullRequests                  *pullRequestStats       `json:"pull_requests"`
}

type totalsByIde struct {
	Ide                           string `json:"ide"`
	UserInitiatedInteractionCount int    `json:"user_initiated_interaction_count"`
	CodeGenerationActivityCount   int    `json:"code_generation_activity_count"`
	CodeAcceptanceActivityCount   int    `json:"code_acceptance_activity_count"`
	LocSuggestedToAddSum          int    `json:"loc_suggested_to_add_sum"`
	LocSuggestedToDeleteSum       int    `json:"loc_suggested_to_delete_sum"`
	LocAddedSum                   int    `json:"loc_added_sum"`
	LocDeletedSum                 int    `json:"loc_deleted_sum"`
}

type totalsByFeature struct {
	Feature                       string `json:"feature"`
	UserInitiatedInteractionCount int    `json:"user_initiated_interaction_count"`
	CodeGenerationActivityCount   int    `json:"code_generation_activity_count"`
	CodeAcceptanceActivityCount   int    `json:"code_acceptance_activity_count"`
	LocSuggestedToAddSum          int    `json:"loc_suggested_to_add_sum"`
	LocSuggestedToDeleteSum       int    `json:"loc_suggested_to_delete_sum"`
	LocAddedSum                   int    `json:"loc_added_sum"`
	LocDeletedSum                 int    `json:"loc_deleted_sum"`
}

type totalsByLangFeature struct {
	Language                    string `json:"language"`
	Feature                     string `json:"feature"`
	CodeGenerationActivityCount int    `json:"code_generation_activity_count"`
	CodeAcceptanceActivityCount int    `json:"code_acceptance_activity_count"`
	LocSuggestedToAddSum        int    `json:"loc_suggested_to_add_sum"`
	LocSuggestedToDeleteSum     int    `json:"loc_suggested_to_delete_sum"`
	LocAddedSum                 int    `json:"loc_added_sum"`
	LocDeletedSum               int    `json:"loc_deleted_sum"`
}

type totalsByLangModel struct {
	Language                    string `json:"language"`
	Model                       string `json:"model"`
	CodeGenerationActivityCount int    `json:"code_generation_activity_count"`
	CodeAcceptanceActivityCount int    `json:"code_acceptance_activity_count"`
	LocSuggestedToAddSum        int    `json:"loc_suggested_to_add_sum"`
	LocSuggestedToDeleteSum     int    `json:"loc_suggested_to_delete_sum"`
	LocAddedSum                 int    `json:"loc_added_sum"`
	LocDeletedSum               int    `json:"loc_deleted_sum"`
}

type pullRequestStats struct {
	TotalReviewed          int `json:"total_reviewed"`
	TotalCreated           int `json:"total_created"`
	TotalCreatedByCopilot  int `json:"total_created_by_copilot"`
	TotalReviewedByCopilot int `json:"total_reviewed_by_copilot"`
}

type totalsByModelFeature struct {
	Model                         string `json:"model"`
	Feature                       string `json:"feature"`
	UserInitiatedInteractionCount int    `json:"user_initiated_interaction_count"`
	CodeGenerationActivityCount   int    `json:"code_generation_activity_count"`
	CodeAcceptanceActivityCount   int    `json:"code_acceptance_activity_count"`
	LocSuggestedToAddSum          int    `json:"loc_suggested_to_add_sum"`
	LocSuggestedToDeleteSum       int    `json:"loc_suggested_to_delete_sum"`
	LocAddedSum                   int    `json:"loc_added_sum"`
	LocDeletedSum                 int    `json:"loc_deleted_sum"`
}

// ExtractEnterpriseMetrics parses enterprise report JSON and extracts to tool-layer tables.
func ExtractEnterpriseMetrics(taskCtx plugin.SubTaskContext) errors.Error {
	data, ok := taskCtx.TaskContext().GetData().(*GhCopilotTaskData)
	if !ok {
		return errors.Default.New("task data is not GhCopilotTaskData")
	}
	connection := data.Connection
	connection.Normalize()

	if !connection.HasEnterprise() {
		taskCtx.GetLogger().Info("No enterprise configured, skipping enterprise metrics extraction")
		return nil
	}

	params := copilotRawParams{
		ConnectionId: data.Options.ConnectionId,
		ScopeId:      data.Options.ScopeId,
		Organization: connection.Organization,
		Endpoint:     connection.Endpoint,
	}

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Table:   rawEnterpriseMetricsTable,
			Options: params,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			// The API returns a flat enterpriseDayTotal object per raw row, not a wrapper.
			var dt enterpriseDayTotal
			if err := errors.Convert(json.Unmarshal(row.Data, &dt)); err != nil {
				return nil, err
			}

			day, parseErr := time.Parse("2006-01-02", dt.Day)
			if parseErr != nil {
				return nil, errors.BadInput.Wrap(parseErr, "invalid day in enterprise report")
			}

			var results []interface{}

			// Main daily metrics
			dailyMetrics := &models.GhCopilotEnterpriseDailyMetrics{
				ConnectionId:        data.Options.ConnectionId,
				ScopeId:             data.Options.ScopeId,
				Day:                 day,
				EnterpriseId:        dt.EnterpriseId,
				DailyActiveUsers:    dt.DailyActiveUsers,
				WeeklyActiveUsers:   dt.WeeklyActiveUsers,
				MonthlyActiveUsers:  dt.MonthlyActiveUsers,
				MonthlyActiveChatUsers:  dt.MonthlyActiveChatUsers,
				MonthlyActiveAgentUsers: dt.MonthlyActiveAgentUsers,
				CopilotActivityMetrics: models.CopilotActivityMetrics{
					UserInitiatedInteractionCount: dt.UserInitiatedInteractionCount,
					CodeGenerationActivityCount:   dt.CodeGenerationActivityCount,
					CodeAcceptanceActivityCount:   dt.CodeAcceptanceActivityCount,
					LocSuggestedToAddSum:          dt.LocSuggestedToAddSum,
					LocSuggestedToDeleteSum:       dt.LocSuggestedToDeleteSum,
					LocAddedSum:                   dt.LocAddedSum,
					LocDeletedSum:                 dt.LocDeletedSum,
				},
			}
			if dt.PullRequests != nil {
				dailyMetrics.PRTotalReviewed = dt.PullRequests.TotalReviewed
				dailyMetrics.PRTotalCreated = dt.PullRequests.TotalCreated
				dailyMetrics.PRTotalCreatedByCopilot = dt.PullRequests.TotalCreatedByCopilot
				dailyMetrics.PRTotalReviewedByCopilot = dt.PullRequests.TotalReviewedByCopilot
			}
			results = append(results, dailyMetrics)

			// By IDE
			for _, ide := range dt.TotalsByIde {
				results = append(results, &models.GhCopilotMetricsByIde{
					ConnectionId: data.Options.ConnectionId,
					ScopeId:      data.Options.ScopeId,
					Day:          day,
					Ide:          ide.Ide,
					CopilotActivityMetrics: models.CopilotActivityMetrics{
						UserInitiatedInteractionCount: ide.UserInitiatedInteractionCount,
						CodeGenerationActivityCount:   ide.CodeGenerationActivityCount,
						CodeAcceptanceActivityCount:   ide.CodeAcceptanceActivityCount,
						LocSuggestedToAddSum:          ide.LocSuggestedToAddSum,
						LocSuggestedToDeleteSum:       ide.LocSuggestedToDeleteSum,
						LocAddedSum:                   ide.LocAddedSum,
						LocDeletedSum:                 ide.LocDeletedSum,
					},
				})
			}

			// By Feature
			for _, f := range dt.TotalsByFeature {
				results = append(results, &models.GhCopilotMetricsByFeature{
					ConnectionId: data.Options.ConnectionId,
					ScopeId:      data.Options.ScopeId,
					Day:          day,
					Feature:      f.Feature,
					CopilotActivityMetrics: models.CopilotActivityMetrics{
						UserInitiatedInteractionCount: f.UserInitiatedInteractionCount,
						CodeGenerationActivityCount:   f.CodeGenerationActivityCount,
						CodeAcceptanceActivityCount:   f.CodeAcceptanceActivityCount,
						LocSuggestedToAddSum:          f.LocSuggestedToAddSum,
						LocSuggestedToDeleteSum:       f.LocSuggestedToDeleteSum,
						LocAddedSum:                   f.LocAddedSum,
						LocDeletedSum:                 f.LocDeletedSum,
					},
				})
			}

			// By Language+Feature
			for _, lf := range dt.TotalsByLanguageFeature {
				results = append(results, &models.GhCopilotMetricsByLanguageFeature{
					ConnectionId: data.Options.ConnectionId,
					ScopeId:      data.Options.ScopeId,
					Day:          day,
					Language:     lf.Language,
					Feature:      lf.Feature,
					CopilotCodeMetrics: models.CopilotCodeMetrics{
						CodeGenerationActivityCount: lf.CodeGenerationActivityCount,
						CodeAcceptanceActivityCount: lf.CodeAcceptanceActivityCount,
						LocSuggestedToAddSum:        lf.LocSuggestedToAddSum,
						LocSuggestedToDeleteSum:     lf.LocSuggestedToDeleteSum,
						LocAddedSum:                 lf.LocAddedSum,
						LocDeletedSum:               lf.LocDeletedSum,
					},
				})
			}

			// By Language+Model
			for _, lm := range dt.TotalsByLanguageModel {
				results = append(results, &models.GhCopilotMetricsByLanguageModel{
					ConnectionId: data.Options.ConnectionId,
					ScopeId:      data.Options.ScopeId,
					Day:          day,
					Language:     lm.Language,
					Model:        lm.Model,
					CopilotCodeMetrics: models.CopilotCodeMetrics{
						CodeGenerationActivityCount: lm.CodeGenerationActivityCount,
						CodeAcceptanceActivityCount: lm.CodeAcceptanceActivityCount,
						LocSuggestedToAddSum:        lm.LocSuggestedToAddSum,
						LocSuggestedToDeleteSum:     lm.LocSuggestedToDeleteSum,
						LocAddedSum:                 lm.LocAddedSum,
						LocDeletedSum:               lm.LocDeletedSum,
					},
				})
			}

			// By Model+Feature
			for _, mf := range dt.TotalsByModelFeature {
				results = append(results, &models.GhCopilotMetricsByModelFeature{
					ConnectionId: data.Options.ConnectionId,
					ScopeId:      data.Options.ScopeId,
					Day:          day,
					Model:        mf.Model,
					Feature:      mf.Feature,
					CopilotActivityMetrics: models.CopilotActivityMetrics{
						UserInitiatedInteractionCount: mf.UserInitiatedInteractionCount,
						CodeGenerationActivityCount:   mf.CodeGenerationActivityCount,
						CodeAcceptanceActivityCount:   mf.CodeAcceptanceActivityCount,
						LocSuggestedToAddSum:          mf.LocSuggestedToAddSum,
						LocSuggestedToDeleteSum:       mf.LocSuggestedToDeleteSum,
						LocAddedSum:                   mf.LocAddedSum,
						LocDeletedSum:                 mf.LocDeletedSum,
					},
				})
			}

			return results, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}
