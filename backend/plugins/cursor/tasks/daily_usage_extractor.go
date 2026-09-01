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
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/cursor/models"
)

type dailyUsageRecord struct {
	Day                      string `json:"day"`
	UserId                   string `json:"userId"`
	Email                    string `json:"email"`
	IsActive                 bool   `json:"isActive"`
	Completions              int    `json:"completions"`
	PremiumRequests          int    `json:"premiumRequests"`
	AgentRequests            int    `json:"agentRequests"`
	ChatRequests             int    `json:"chatRequests"`
	ComposerRequests         int    `json:"composerRequests"`
	TotalTabsAccepted        int    `json:"totalTabsAccepted"`
	TotalTabsShown           int    `json:"totalTabsShown"`
	UsageBasedReqs           int    `json:"usageBasedReqs"`
	SubscriptionIncludedReqs int    `json:"subscriptionIncludedReqs"`
	MostUsedModel            string `json:"mostUsedModel"`
	ClientVersion            string `json:"clientVersion"`
	TotalLinesAdded          int    `json:"totalLinesAdded"`
	TotalLinesDeleted        int    `json:"totalLinesDeleted"`
	AcceptedLinesAdded       int    `json:"acceptedLinesAdded"`
	AcceptedLinesDeleted     int    `json:"acceptedLinesDeleted"`
	TotalApplies             int    `json:"totalApplies"`
	TotalAccepts             int    `json:"totalAccepts"`
	TotalRejects             int    `json:"totalRejects"`
}

// ExtractDailyUsage parses raw daily usage records into tool-layer tables.
func ExtractDailyUsage(taskCtx plugin.SubTaskContext) errors.Error {
	data, ok := taskCtx.TaskContext().GetData().(*CursorTaskData)
	if !ok {
		return errors.Default.New("task data is not CursorTaskData")
	}
	logger := taskCtx.GetLogger()

	extractor, err := newCursorStatefulExtractor(&cursorStatefulExtractorArgs[dailyUsageRecord]{
		SubtaskCommonArgs: cursorSubtaskCommonArgs(taskCtx, data, rawDailyUsageTable),
		ConnectionId:      data.Options.ConnectionId,
		ScopeId:           data.Options.ScopeId,
		ToolTable:         models.CursorDailyUsage{}.TableName(),
		Extract: func(record *dailyUsageRecord, row *helper.RawData) ([]any, errors.Error) {
			userId := strings.TrimSpace(record.UserId)
			if userId == "" {
				userId = strings.TrimSpace(record.Email)
			}
			if userId == "" {
				rawID := uint64(0)
				if row != nil {
					rawID = row.ID
				}
				logger.Warn(nil, "skipping daily usage raw row id=%d: missing userId and email", rawID)
				return nil, nil
			}

			usageDate, parseErr := time.Parse("2006-01-02", strings.TrimSpace(record.Day))
			if parseErr != nil {
				rawID := uint64(0)
				if row != nil {
					rawID = row.ID
				}
				logger.Warn(nil, "invalid day in daily usage raw row id=%d: %v", rawID, parseErr)
				return nil, errors.BadInput.Wrap(parseErr, "invalid day format in daily usage record")
			}

			usage := &models.CursorDailyUsage{
				ConnectionId:             data.Options.ConnectionId,
				ScopeId:                  data.Options.ScopeId,
				UserId:                   userId,
				UsageDate:                usageDate,
				Email:                    strings.TrimSpace(record.Email),
				IsActive:                 record.IsActive,
				Completions:              record.Completions,
				PremiumRequests:          record.PremiumRequests,
				AgentRequests:            record.AgentRequests,
				ChatRequests:             record.ChatRequests,
				ComposerRequests:         record.ComposerRequests,
				TabsAccepted:             record.TotalTabsAccepted,
				TabsShown:                record.TotalTabsShown,
				UsageBasedReqs:           record.UsageBasedReqs,
				SubscriptionIncludedReqs: record.SubscriptionIncludedReqs,
				MostUsedModel:            strings.TrimSpace(record.MostUsedModel),
				ClientVersion:            strings.TrimSpace(record.ClientVersion),
				AcceptedLinesAdded:       record.AcceptedLinesAdded,
				AcceptedLinesDeleted:     record.AcceptedLinesDeleted,
				TotalLinesAdded:          record.TotalLinesAdded,
				TotalLinesDeleted:        record.TotalLinesDeleted,
				TotalApplies:             record.TotalApplies,
				TotalAccepts:             record.TotalAccepts,
				TotalRejects:             record.TotalRejects,
				LinesAdded:               record.AcceptedLinesAdded,
				LinesDeleted:             record.AcceptedLinesDeleted,
			}
			return []any{usage}, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}
