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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
)

var ExtractTeamMetricsMeta = plugin.SubTaskMeta{
	Name:             "extractTeamMetrics",
	EntryPoint:       ExtractTeamMetrics,
	EnabledByDefault: true,
	Description:      "Extract raw team-level Copilot metrics into tool layer team metrics tables",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
	DependencyTables: []string{rawCopilotTeamMetricsTable},
	Dependencies:     []*plugin.SubTaskMeta{&CollectTeamMetricsMeta},
	ProductTables: []string{
		models.GhCopilotTeamDailyMetrics{}.TableName(),
		models.GhCopilotTeamCompletions{}.TableName(),
		models.GhCopilotTeamIdeChat{}.TableName(),
		models.GhCopilotTeamDotcomChat{}.TableName(),
		models.GhCopilotTeamDotcomPrs{}.TableName(),
	},
}

type teamMetricDay struct {
	Date                      string                        `json:"date"`
	TotalActiveUsers          int                           `json:"total_active_users"`
	TotalEngagedUsers         int                           `json:"total_engaged_users"`
	CopilotIdeCodeCompletions *teamIdeCodeCompletions       `json:"copilot_ide_code_completions"`
	CopilotIdeChat            *teamIdeChatMetrics           `json:"copilot_ide_chat"`
	CopilotDotcomChat         *teamDotcomChatMetrics        `json:"copilot_dotcom_chat"`
	CopilotDotcomPullRequests *teamDotcomPullRequestMetrics `json:"copilot_dotcom_pull_requests"`
}

type teamIdeCodeCompletions struct {
	TotalEngagedUsers int                    `json:"total_engaged_users"`
	Languages         []teamLangEngagedUsers `json:"languages"`
	Editors           []teamCompletionEditor `json:"editors"`
}

type teamLangEngagedUsers struct {
	Name              string `json:"name"`
	TotalEngagedUsers int    `json:"total_engaged_users"`
}

type teamCompletionEditor struct {
	Name              string                `json:"name"`
	TotalEngagedUsers int                   `json:"total_engaged_users"`
	Models            []teamCompletionModel `json:"models"`
}

type teamCompletionModel struct {
	Name                    string                   `json:"name"`
	IsCustomModel           bool                     `json:"is_custom_model"`
	CustomModelTrainingDate *string                  `json:"custom_model_training_date"`
	TotalEngagedUsers       int                      `json:"total_engaged_users"`
	Languages               []teamCompletionLanguage `json:"languages"`
}

type teamCompletionLanguage struct {
	Name                    string `json:"name"`
	TotalEngagedUsers       int    `json:"total_engaged_users"`
	TotalCodeSuggestions    int    `json:"total_code_suggestions"`
	TotalCodeAcceptances    int    `json:"total_code_acceptances"`
	TotalCodeLinesSuggested int    `json:"total_code_lines_suggested"`
	TotalCodeLinesAccepted  int    `json:"total_code_lines_accepted"`
}

type teamIdeChatMetrics struct {
	TotalEngagedUsers int                 `json:"total_engaged_users"`
	Editors           []teamIdeChatEditor `json:"editors"`
}

type teamIdeChatEditor struct {
	Name              string             `json:"name"`
	TotalEngagedUsers int                `json:"total_engaged_users"`
	Models            []teamIdeChatModel `json:"models"`
}

type teamIdeChatModel struct {
	Name                     string  `json:"name"`
	IsCustomModel            bool    `json:"is_custom_model"`
	CustomModelTrainingDate  *string `json:"custom_model_training_date"`
	TotalEngagedUsers        int     `json:"total_engaged_users"`
	TotalChats               int     `json:"total_chats"`
	TotalChatInsertionEvents int     `json:"total_chat_insertion_events"`
	TotalChatCopyEvents      int     `json:"total_chat_copy_events"`
}

type teamDotcomChatMetrics struct {
	TotalEngagedUsers int                   `json:"total_engaged_users"`
	Models            []teamDotcomChatModel `json:"models"`
}

type teamDotcomChatModel struct {
	Name                    string  `json:"name"`
	IsCustomModel           bool    `json:"is_custom_model"`
	CustomModelTrainingDate *string `json:"custom_model_training_date"`
	TotalEngagedUsers       int     `json:"total_engaged_users"`
	TotalChats              int     `json:"total_chats"`
}

type teamDotcomPullRequestMetrics struct {
	TotalEngagedUsers int                      `json:"total_engaged_users"`
	Repositories      []teamDotcomPrRepository `json:"repositories"`
}

type teamDotcomPrRepository struct {
	Name              string              `json:"name"`
	TotalEngagedUsers int                 `json:"total_engaged_users"`
	Models            []teamDotcomPrModel `json:"models"`
}

type teamDotcomPrModel struct {
	Name                    string  `json:"name"`
	IsCustomModel           bool    `json:"is_custom_model"`
	CustomModelTrainingDate *string `json:"custom_model_training_date"`
	TotalEngagedUsers       int     `json:"total_engaged_users"`
	TotalPrSummariesCreated int     `json:"total_pr_summaries_created"`
}

func parseTeamMetricsDate(raw string) (time.Time, errors.Error) {
	date, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Time{}, errors.BadInput.Wrap(err, "invalid team metrics date")
	}
	return date, nil
}

func parseTeamCustomModelTrainingDate(raw *string) (*time.Time, errors.Error) {
	if raw == nil {
		return nil, nil
	}
	value := strings.TrimSpace(*raw)
	if value == "" {
		return nil, nil
	}

	date, err := time.Parse("2006-01-02", value)
	if err == nil {
		return &date, nil
	}

	date, err = time.Parse(time.RFC3339, value)
	if err == nil {
		return &date, nil
	}

	return nil, errors.BadInput.Wrap(err, "invalid custom model training date in team metrics")
}

func ExtractTeamMetrics(taskCtx plugin.SubTaskContext) errors.Error {
	data, ok := taskCtx.TaskContext().GetData().(*GhCopilotTaskData)
	if !ok {
		return errors.Default.New("task data is not GhCopilotTaskData")
	}
	connection := data.Connection
	connection.Normalize()

	org := strings.TrimSpace(connection.Organization)
	if org == "" {
		taskCtx.GetLogger().Info("No organization configured, skipping team metrics extraction")
		return nil
	}

	params := copilotRawParams{
		ConnectionId: data.Options.ConnectionId,
		ScopeId:      data.Options.ScopeId,
		Organization: org,
		Endpoint:     connection.Endpoint,
	}

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Table:   rawCopilotTeamMetricsTable,
			Options: params,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			apiDay := &teamMetricDay{}
			if e := json.Unmarshal(row.Data, apiDay); e != nil {
				return nil, errors.Convert(e)
			}

			team := &simpleCopilotTeam{}
			if e := json.Unmarshal(row.Input, team); e != nil {
				return nil, errors.Convert(e)
			}
			if team.Slug == "" {
				return nil, errors.BadInput.New("missing team slug in raw team metrics input")
			}

			date, parseErr := parseTeamMetricsDate(apiDay.Date)
			if parseErr != nil {
				return nil, parseErr
			}

			results := make([]interface{}, 0)

			daily := &models.GhCopilotTeamDailyMetrics{
				ConnectionId:      data.Options.ConnectionId,
				ScopeId:           data.Options.ScopeId,
				TeamSlug:          team.Slug,
				Date:              date,
				TotalActiveUsers:  apiDay.TotalActiveUsers,
				TotalEngagedUsers: apiDay.TotalEngagedUsers,
			}
			if apiDay.CopilotIdeCodeCompletions != nil {
				daily.CompletionsTotalEngagedUsers = apiDay.CopilotIdeCodeCompletions.TotalEngagedUsers
			}
			if apiDay.CopilotIdeChat != nil {
				daily.IdeChatTotalEngagedUsers = apiDay.CopilotIdeChat.TotalEngagedUsers
			}
			if apiDay.CopilotDotcomChat != nil {
				daily.DotcomChatTotalEngagedUsers = apiDay.CopilotDotcomChat.TotalEngagedUsers
			}
			if apiDay.CopilotDotcomPullRequests != nil {
				daily.DotcomPrTotalEngagedUsers = apiDay.CopilotDotcomPullRequests.TotalEngagedUsers
			}
			results = append(results, daily)

			if apiDay.CopilotIdeCodeCompletions != nil {
				for _, editor := range apiDay.CopilotIdeCodeCompletions.Editors {
					for _, model := range editor.Models {
						trainingDate, dateErr := parseTeamCustomModelTrainingDate(model.CustomModelTrainingDate)
						if dateErr != nil {
							return nil, dateErr
						}
						for _, language := range model.Languages {
							results = append(results, &models.GhCopilotTeamCompletions{
								ConnectionId:            data.Options.ConnectionId,
								ScopeId:                 data.Options.ScopeId,
								TeamSlug:                team.Slug,
								Date:                    date,
								Editor:                  editor.Name,
								Model:                   model.Name,
								Language:                language.Name,
								TotalEngagedUsers:       language.TotalEngagedUsers,
								TotalCodeSuggestions:    language.TotalCodeSuggestions,
								TotalCodeAcceptances:    language.TotalCodeAcceptances,
								TotalCodeLinesSuggested: language.TotalCodeLinesSuggested,
								TotalCodeLinesAccepted:  language.TotalCodeLinesAccepted,
								IsCustomModel:           model.IsCustomModel,
								CustomModelTrainingDate: trainingDate,
							})
						}
					}
				}
			}

			if apiDay.CopilotIdeChat != nil {
				for _, editor := range apiDay.CopilotIdeChat.Editors {
					for _, model := range editor.Models {
						trainingDate, dateErr := parseTeamCustomModelTrainingDate(model.CustomModelTrainingDate)
						if dateErr != nil {
							return nil, dateErr
						}
						results = append(results, &models.GhCopilotTeamIdeChat{
							ConnectionId:             data.Options.ConnectionId,
							ScopeId:                  data.Options.ScopeId,
							TeamSlug:                 team.Slug,
							Date:                     date,
							Editor:                   editor.Name,
							Model:                    model.Name,
							TotalEngagedUsers:        model.TotalEngagedUsers,
							TotalChats:               model.TotalChats,
							TotalChatInsertionEvents: model.TotalChatInsertionEvents,
							TotalChatCopyEvents:      model.TotalChatCopyEvents,
							IsCustomModel:            model.IsCustomModel,
							CustomModelTrainingDate:  trainingDate,
						})
					}
				}
			}

			if apiDay.CopilotDotcomChat != nil {
				for _, model := range apiDay.CopilotDotcomChat.Models {
					trainingDate, dateErr := parseTeamCustomModelTrainingDate(model.CustomModelTrainingDate)
					if dateErr != nil {
						return nil, dateErr
					}
					results = append(results, &models.GhCopilotTeamDotcomChat{
						ConnectionId:            data.Options.ConnectionId,
						ScopeId:                 data.Options.ScopeId,
						TeamSlug:                team.Slug,
						Date:                    date,
						Model:                   model.Name,
						TotalEngagedUsers:       model.TotalEngagedUsers,
						TotalChats:              model.TotalChats,
						IsCustomModel:           model.IsCustomModel,
						CustomModelTrainingDate: trainingDate,
					})
				}
			}

			if apiDay.CopilotDotcomPullRequests != nil {
				for _, repository := range apiDay.CopilotDotcomPullRequests.Repositories {
					for _, model := range repository.Models {
						trainingDate, dateErr := parseTeamCustomModelTrainingDate(model.CustomModelTrainingDate)
						if dateErr != nil {
							return nil, dateErr
						}
						results = append(results, &models.GhCopilotTeamDotcomPrs{
							ConnectionId:            data.Options.ConnectionId,
							ScopeId:                 data.Options.ScopeId,
							TeamSlug:                team.Slug,
							Date:                    date,
							Repository:              repository.Name,
							Model:                   model.Name,
							TotalEngagedUsers:       model.TotalEngagedUsers,
							TotalPrSummariesCreated: model.TotalPrSummariesCreated,
							IsCustomModel:           model.IsCustomModel,
							CustomModelTrainingDate: trainingDate,
						})
					}
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
