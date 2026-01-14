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
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
)

type copilotMetricsDay struct {
	Date              string `json:"date"`
	TotalActiveUsers  int    `json:"total_active_users"`
	TotalEngagedUsers int    `json:"total_engaged_users"`

	IdeCodeCompletions copilotIdeCodeCompletions `json:"copilot_ide_code_completions"`
	IdeChat            copilotIdeChat            `json:"copilot_ide_chat"`
	DotcomChat         copilotDotcomChat         `json:"copilot_dotcom_chat"`
}

type copilotIdeCodeCompletions struct {
	TotalEngagedUsers int                        `json:"total_engaged_users"`
	Editors           []copilotEditorCompletions `json:"editors"`
}

type copilotEditorCompletions struct {
	Name              string                    `json:"name"`
	TotalEngagedUsers int                       `json:"total_engaged_users"`
	Models            []copilotModelCompletions `json:"models"`
}

type copilotModelCompletions struct {
	Name              string                   `json:"name"`
	IsCustomModel     bool                     `json:"is_custom_model"`
	TotalEngagedUsers int                      `json:"total_engaged_users"`
	Languages         []copilotLanguageMetrics `json:"languages"`
}

type copilotLanguageMetrics struct {
	Name                    string `json:"name"`
	TotalEngagedUsers       int    `json:"total_engaged_users"`
	TotalCodeSuggestions    int    `json:"total_code_suggestions"`
	TotalCodeAcceptances    int    `json:"total_code_acceptances"`
	TotalCodeLinesSuggested int    `json:"total_code_lines_suggested"`
	TotalCodeLinesAccepted  int    `json:"total_code_lines_accepted"`
}

type copilotIdeChat struct {
	TotalEngagedUsers int                 `json:"total_engaged_users"`
	Editors           []copilotEditorChat `json:"editors"`
}

type copilotEditorChat struct {
	Name              string             `json:"name"`
	TotalEngagedUsers int                `json:"total_engaged_users"`
	Models            []copilotModelChat `json:"models"`
}

type copilotModelChat struct {
	Name                     string `json:"name"`
	IsCustomModel            bool   `json:"is_custom_model"`
	TotalEngagedUsers        int    `json:"total_engaged_users"`
	TotalChats               int    `json:"total_chats"`
	TotalChatCopyEvents      int    `json:"total_chat_copy_events"`
	TotalChatInsertionEvents int    `json:"total_chat_insertion_events"`
}

type copilotDotcomChat struct {
	TotalEngagedUsers int                  `json:"total_engaged_users"`
	Models            []copilotDotcomModel `json:"models"`
}

type copilotDotcomModel struct {
	Name              string `json:"name"`
	IsCustomModel     bool   `json:"is_custom_model"`
	TotalEngagedUsers int    `json:"total_engaged_users"`
	TotalChats        int    `json:"total_chats"`
}

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

func ExtractCopilotOrgMetrics(taskCtx plugin.SubTaskContext) errors.Error {
	data, ok := taskCtx.TaskContext().GetData().(*CopilotTaskData)
	if !ok {
		return errors.Default.New("task data is not CopilotTaskData")
	}
	connection := data.Connection
	connection.Normalize()

	params := copilotRawParams{
		ConnectionId: data.Options.ConnectionId,
		ScopeId:      data.Options.ScopeId,
		Organization: connection.Organization,
		Endpoint:     connection.Endpoint,
	}

	// Extract seat assignments first so we can derive seat counts for org metrics.
	// NOTE: Keep this extractor stateless to avoid SubtaskStateManager collisions inside this subtask.
	// The state key does not include raw table name, so multiple stateful extractors would race/skip.
	seatsExtractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Table:   rawCopilotSeatsTable,
			Options: params,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			seat := &copilotSeatResponse{}
			if err := errors.Convert(json.Unmarshal(row.Data, seat)); err != nil {
				return nil, err
			}

			createdAt, parseErr := time.Parse(time.RFC3339, seat.CreatedAt)
			if parseErr != nil {
				return nil, errors.BadInput.Wrap(parseErr, "invalid seat created_at")
			}
			updatedAt, parseErr := time.Parse(time.RFC3339, seat.UpdatedAt)
			if parseErr != nil {
				return nil, errors.BadInput.Wrap(parseErr, "invalid seat updated_at")
			}

			parseOptional := func(v *string) (*time.Time, errors.Error) {
				if v == nil || *v == "" {
					return nil, nil
				}
				// GitHub may return RFC3339 timestamps or date-only strings (YYYY-MM-DD) for some fields.
				if t, parseErr := time.Parse(time.RFC3339, *v); parseErr == nil {
					return &t, nil
				}
				t, parseErr := time.Parse("2006-01-02", *v)
				if parseErr != nil {
					return nil, errors.BadInput.Wrap(parseErr, "invalid timestamp")
				}
				return &t, nil
			}

			lastAuth, err := parseOptional(seat.LastAuthenticatedAt)
			if err != nil {
				return nil, err
			}
			lastAct, err := parseOptional(seat.LastActivityAt)
			if err != nil {
				return nil, err
			}
			pendingCancel, err := parseOptional(seat.PendingCancellationDate)
			if err != nil {
				return nil, err
			}

			toolSeat := &models.CopilotSeat{
				ConnectionId:            data.Options.ConnectionId,
				Organization:            connection.Organization,
				UserLogin:               seat.Assignee.Login,
				UserId:                  seat.Assignee.Id,
				PlanType:                seat.PlanType,
				CreatedAt:               createdAt,
				LastActivityAt:          lastAct,
				LastActivityEditor:      seat.LastActivityEditor,
				LastAuthenticatedAt:     lastAuth,
				PendingCancellationDate: pendingCancel,
				UpdatedAt:               updatedAt,
			}

			return []interface{}{toolSeat}, nil
		},
	})
	if err != nil {
		return err
	}
	if err := seatsExtractor.Execute(); err != nil {
		return err
	}

	// Derive seat counts from extracted assignments.
	db := taskCtx.GetDal()
	seatTotal, err := db.Count(dal.From(&models.CopilotSeat{}), dal.Where(
		"connection_id = ? AND organization = ?",
		data.Options.ConnectionId,
		connection.Organization,
	))
	if err != nil {
		return errors.Default.Wrap(err, "failed to count copilot seats")
	}
	seatActive, err := db.Count(dal.From(&models.CopilotSeat{}), dal.Where(
		"connection_id = ? AND organization = ? AND last_activity_at IS NOT NULL",
		data.Options.ConnectionId,
		connection.Organization,
	))
	if err != nil {
		return errors.Default.Wrap(err, "failed to count active copilot seats")
	}

	// Keep existing org metrics in sync even when the stateful metrics extractor
	// has nothing new to process (e.g., incremental runs with no new raw metrics).
	if db.HasTable(&models.CopilotOrgMetrics{}) {
		err = db.UpdateColumns(
			&models.CopilotOrgMetrics{},
			[]dal.DalSet{{ColumnName: "seat_total", Value: seatTotal}, {ColumnName: "seat_active_count", Value: seatActive}},
			dal.Where("connection_id = ? AND scope_id = ?", data.Options.ConnectionId, data.Options.ScopeId),
		)
		if err != nil {
			return errors.Default.Wrap(err, "failed to update copilot org metrics seat counts")
		}
	}

	metricsExtractor, err := helper.NewStatefulApiExtractor(&helper.StatefulApiExtractorArgs[copilotMetricsDay]{
		SubtaskCommonArgs: &helper.SubtaskCommonArgs{
			SubTaskContext: taskCtx,
			Table:          rawCopilotMetricsTable,
			Params:         params,
			SubtaskConfig:  params,
		},
		Extract: func(day *copilotMetricsDay, row *helper.RawData) ([]any, errors.Error) {
			date, err := time.Parse("2006-01-02", day.Date)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "invalid metrics date")
			}

			normalizeDim := func(v, fallback string) string {
				v = strings.TrimSpace(v)
				if v == "" {
					return fallback
				}
				return v
			}

			completionSuggestions := 0
			completionAcceptances := 0
			completionLinesSuggested := 0
			completionLinesAccepted := 0
			languageMetrics := make([]any, 0, 64)
			for _, editor := range day.IdeCodeCompletions.Editors {
				editorName := normalizeDim(editor.Name, "unknown")
				for _, model := range editor.Models {
					for _, lang := range model.Languages {
						completionSuggestions += lang.TotalCodeSuggestions
						completionAcceptances += lang.TotalCodeAcceptances
						completionLinesSuggested += lang.TotalCodeLinesSuggested
						completionLinesAccepted += lang.TotalCodeLinesAccepted

						toolLang := &models.CopilotLanguageMetrics{
							ConnectionId:   data.Options.ConnectionId,
							ScopeId:        data.Options.ScopeId,
							Date:           date,
							Editor:         editorName,
							Language:       normalizeDim(lang.Name, "unknown"),
							EngagedUsers:   lang.TotalEngagedUsers,
							Suggestions:    lang.TotalCodeSuggestions,
							Acceptances:    lang.TotalCodeAcceptances,
							LinesSuggested: lang.TotalCodeLinesSuggested,
							LinesAccepted:  lang.TotalCodeLinesAccepted,
						}
						languageMetrics = append(languageMetrics, toolLang)
					}
				}
			}

			ideChats := 0
			ideChatCopyEvents := 0
			ideChatInsertionEvents := 0
			for _, editor := range day.IdeChat.Editors {
				for _, model := range editor.Models {
					ideChats += model.TotalChats
					ideChatCopyEvents += model.TotalChatCopyEvents
					ideChatInsertionEvents += model.TotalChatInsertionEvents
				}
			}

			dotcomChats := 0
			for _, model := range day.DotcomChat.Models {
				dotcomChats += model.TotalChats
			}

			metric := &models.CopilotOrgMetrics{
				ConnectionId:             data.Options.ConnectionId,
				ScopeId:                  data.Options.ScopeId,
				Date:                     date,
				TotalActiveUsers:         day.TotalActiveUsers,
				TotalEngagedUsers:        day.TotalEngagedUsers,
				CompletionSuggestions:    completionSuggestions,
				CompletionAcceptances:    completionAcceptances,
				CompletionLinesSuggested: completionLinesSuggested,
				CompletionLinesAccepted:  completionLinesAccepted,
				IdeChats:                 ideChats,
				IdeChatCopyEvents:        ideChatCopyEvents,
				IdeChatInsertionEvents:   ideChatInsertionEvents,
				IdeChatEngagedUsers:      day.IdeChat.TotalEngagedUsers,
				DotcomChats:              dotcomChats,
				DotcomChatEngagedUsers:   day.DotcomChat.TotalEngagedUsers,
				SeatTotal:                int(seatTotal),
				SeatActiveCount:          int(seatActive),
			}
			results := make([]any, 0, 1+len(languageMetrics))
			results = append(results, metric)
			results = append(results, languageMetrics...)
			return results, nil
		},
	})
	if err != nil {
		return err
	}
	if err := metricsExtractor.Execute(); err != nil {
		return err
	}

	return nil
}
