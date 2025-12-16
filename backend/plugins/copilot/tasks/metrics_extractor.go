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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/copilot/models"
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
	Editors []copilotEditorCompletions `json:"editors"`
}

type copilotEditorCompletions struct {
	Models []copilotModelCompletions `json:"models"`
}

type copilotModelCompletions struct {
	Languages []copilotLanguageMetrics `json:"languages"`
}

type copilotLanguageMetrics struct {
	TotalCodeSuggestions    int `json:"total_code_suggestions"`
	TotalCodeAcceptances    int `json:"total_code_acceptances"`
	TotalCodeLinesSuggested int `json:"total_code_lines_suggested"`
	TotalCodeLinesAccepted  int `json:"total_code_lines_accepted"`
}

type copilotIdeChat struct {
	TotalEngagedUsers int                 `json:"total_engaged_users"`
	Editors           []copilotEditorChat `json:"editors"`
}

type copilotEditorChat struct {
	Models []copilotModelChat `json:"models"`
}

type copilotModelChat struct {
	TotalChats               int `json:"total_chats"`
	TotalChatCopyEvents      int `json:"total_chat_copy_events"`
	TotalChatInsertionEvents int `json:"total_chat_insertion_events"`
}

type copilotDotcomChat struct {
	TotalEngagedUsers int                  `json:"total_engaged_users"`
	Models            []copilotDotcomModel `json:"models"`
}

type copilotDotcomModel struct {
	TotalChats int `json:"total_chats"`
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

			completionSuggestions := 0
			completionAcceptances := 0
			completionLinesSuggested := 0
			completionLinesAccepted := 0
			for _, editor := range day.IdeCodeCompletions.Editors {
				for _, model := range editor.Models {
					for _, lang := range model.Languages {
						completionSuggestions += lang.TotalCodeSuggestions
						completionAcceptances += lang.TotalCodeAcceptances
						completionLinesSuggested += lang.TotalCodeLinesSuggested
						completionLinesAccepted += lang.TotalCodeLinesAccepted
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
			}
			return []any{metric}, nil
		},
	})
	if err != nil {
		return err
	}
	if err := metricsExtractor.Execute(); err != nil {
		return err
	}

	seatsExtractor, err := helper.NewStatefulApiExtractor(&helper.StatefulApiExtractorArgs[copilotSeatResponse]{
		SubtaskCommonArgs: &helper.SubtaskCommonArgs{
			SubTaskContext: taskCtx,
			Table:          rawCopilotSeatsTable,
			Params:         params,
			SubtaskConfig:  params,
		},
		Extract: func(seat *copilotSeatResponse, row *helper.RawData) ([]any, errors.Error) {
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
				t, parseErr := time.Parse(time.RFC3339, *v)
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

			return []any{toolSeat}, nil
		},
	})
	if err != nil {
		return err
	}

	return seatsExtractor.Execute()
}
