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
	"fmt"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/q_dev/models"
	"math"
	"time"
)

var _ plugin.SubTaskEntryPoint = ConvertQDevUserMetrics

// ConvertQDevUserMetrics 按用户聚合指标 (enhanced with display name support)
func ConvertQDevUserMetrics(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*QDevTaskData)
	db := taskCtx.GetDal()

	// 清空之前聚合的数据
	clauses := []dal.Clause{
		dal.Where("connection_id = ?", data.Options.ConnectionId),
	}
	err := db.Delete(&models.QDevUserMetrics{}, clauses...)
	if err != nil {
		return errors.Default.Wrap(err, "failed to delete previous user metrics")
	}

	// 聚合数据 (updated to include display name)
	userDataMap := make(map[string]*UserMetricsAggregationWithDisplayName)

	cursor, err := db.Cursor(
		dal.From(&models.QDevUserData{}),
		dal.Where("connection_id = ?", data.Options.ConnectionId),
	)
	if err != nil {
		return errors.Default.Wrap(err, "failed to get user data cursor")
	}
	defer cursor.Close()

	taskCtx.SetProgress(0, -1)

	// 汇总每个用户的数据
	for cursor.Next() {
		userData := &models.QDevUserData{}
		err = db.Fetch(cursor, userData)
		if err != nil {
			return errors.Default.Wrap(err, "failed to fetch user data")
		}

		// 获取或创建用户聚合
		aggregation, ok := userDataMap[userData.UserId]
		if !ok {
			// Resolve display name for new user (new functionality)
			displayName := resolveDisplayNameForAggregation(userData.UserId, data.IdentityClient)
			// If user data already has display name, use it; otherwise use resolved name
			if userData.DisplayName != "" {
				displayName = userData.DisplayName
			}

			aggregation = &UserMetricsAggregationWithDisplayName{
				ConnectionId: userData.ConnectionId,
				UserId:       userData.UserId,
				DisplayName:  displayName, // New field
				FirstDate:    userData.Date,
				LastDate:     userData.Date,
				DataCount:    0,
			}
			userDataMap[userData.UserId] = aggregation
		}

		// 更新日期范围
		if userData.Date.Before(aggregation.FirstDate) {
			aggregation.FirstDate = userData.Date
		}
		if userData.Date.After(aggregation.LastDate) {
			aggregation.LastDate = userData.Date
		}

		// 累加指标
		aggregation.DataCount++
		aggregation.TotalCodeReview_FindingsCount += userData.CodeReview_FindingsCount
		aggregation.TotalCodeReview_SucceededEventCount += userData.CodeReview_SucceededEventCount
		aggregation.TotalInlineChat_AcceptanceEventCount += userData.InlineChat_AcceptanceEventCount
		aggregation.TotalInlineChat_AcceptedLineAdditions += userData.InlineChat_AcceptedLineAdditions
		aggregation.TotalInlineChat_AcceptedLineDeletions += userData.InlineChat_AcceptedLineDeletions
		aggregation.TotalInlineChat_DismissalEventCount += userData.InlineChat_DismissalEventCount
		aggregation.TotalInlineChat_DismissedLineAdditions += userData.InlineChat_DismissedLineAdditions
		aggregation.TotalInlineChat_DismissedLineDeletions += userData.InlineChat_DismissedLineDeletions
		aggregation.TotalInlineChat_RejectedLineAdditions += userData.InlineChat_RejectedLineAdditions
		aggregation.TotalInlineChat_RejectedLineDeletions += userData.InlineChat_RejectedLineDeletions
		aggregation.TotalInlineChat_RejectionEventCount += userData.InlineChat_RejectionEventCount
		aggregation.TotalInlineChat_TotalEventCount += userData.InlineChat_TotalEventCount
		aggregation.TotalInline_AICodeLines += userData.Inline_AICodeLines
		aggregation.TotalInline_AcceptanceCount += userData.Inline_AcceptanceCount
		aggregation.TotalInline_SuggestionsCount += userData.Inline_SuggestionsCount
	}

	// 计算每个用户的平均指标和总天数
	for _, aggregation := range userDataMap {
		// 创建指标记录 (updated to use new method)
		metrics := aggregation.ToUserMetrics()

		// 存储聚合指标
		err = db.Create(metrics)
		if err != nil {
			return errors.Default.Wrap(err, "failed to create user metrics")
		}

		taskCtx.IncProgress(1)
	}

	return nil
}

// UserMetricsAggregationWithDisplayName 聚合过程中用于保存用户指标的结构 (enhanced with display name)
type UserMetricsAggregationWithDisplayName struct {
	ConnectionId                           uint64
	UserId                                 string
	DisplayName                            string // New field for display name
	FirstDate                              time.Time
	LastDate                               time.Time
	DataCount                              int
	TotalCodeReview_FindingsCount          int
	TotalCodeReview_SucceededEventCount    int
	TotalInlineChat_AcceptanceEventCount   int
	TotalInlineChat_AcceptedLineAdditions  int
	TotalInlineChat_AcceptedLineDeletions  int
	TotalInlineChat_DismissalEventCount    int
	TotalInlineChat_DismissedLineAdditions int
	TotalInlineChat_DismissedLineDeletions int
	TotalInlineChat_RejectedLineAdditions  int
	TotalInlineChat_RejectedLineDeletions  int
	TotalInlineChat_RejectionEventCount    int
	TotalInlineChat_TotalEventCount        int
	TotalInline_AICodeLines                int
	TotalInline_AcceptanceCount            int
	TotalInline_SuggestionsCount           int
}

// ToUserMetrics converts aggregation data to QDevUserMetrics model
func (aggregation *UserMetricsAggregationWithDisplayName) ToUserMetrics() *models.QDevUserMetrics {
	metrics := &models.QDevUserMetrics{
		ConnectionId: aggregation.ConnectionId,
		UserId:       aggregation.UserId,
		DisplayName:  aggregation.DisplayName, // New field
		FirstDate:    aggregation.FirstDate,
		LastDate:     aggregation.LastDate,
	}

	// 计算总天数
	metrics.TotalDays = int(math.Round(aggregation.LastDate.Sub(aggregation.FirstDate).Hours()/24)) + 1

	// 设置总计指标
	metrics.TotalCodeReview_FindingsCount = aggregation.TotalCodeReview_FindingsCount
	metrics.TotalCodeReview_SucceededEventCount = aggregation.TotalCodeReview_SucceededEventCount
	metrics.TotalInlineChat_AcceptanceEventCount = aggregation.TotalInlineChat_AcceptanceEventCount
	metrics.TotalInlineChat_AcceptedLineAdditions = aggregation.TotalInlineChat_AcceptedLineAdditions
	metrics.TotalInlineChat_AcceptedLineDeletions = aggregation.TotalInlineChat_AcceptedLineDeletions
	metrics.TotalInlineChat_DismissalEventCount = aggregation.TotalInlineChat_DismissalEventCount
	metrics.TotalInlineChat_DismissedLineAdditions = aggregation.TotalInlineChat_DismissedLineAdditions
	metrics.TotalInlineChat_DismissedLineDeletions = aggregation.TotalInlineChat_DismissedLineDeletions
	metrics.TotalInlineChat_RejectedLineAdditions = aggregation.TotalInlineChat_RejectedLineAdditions
	metrics.TotalInlineChat_RejectedLineDeletions = aggregation.TotalInlineChat_RejectedLineDeletions
	metrics.TotalInlineChat_RejectionEventCount = aggregation.TotalInlineChat_RejectionEventCount
	metrics.TotalInlineChat_TotalEventCount = aggregation.TotalInlineChat_TotalEventCount
	metrics.TotalInline_AICodeLines = aggregation.TotalInline_AICodeLines
	metrics.TotalInline_AcceptanceCount = aggregation.TotalInline_AcceptanceCount
	metrics.TotalInline_SuggestionsCount = aggregation.TotalInline_SuggestionsCount

	// 计算平均值指标
	if metrics.TotalDays > 0 {
		metrics.AvgCodeReview_FindingsCount = float64(aggregation.TotalCodeReview_FindingsCount) / float64(metrics.TotalDays)
		metrics.AvgCodeReview_SucceededEventCount = float64(aggregation.TotalCodeReview_SucceededEventCount) / float64(metrics.TotalDays)
		metrics.AvgInlineChat_AcceptanceEventCount = float64(aggregation.TotalInlineChat_AcceptanceEventCount) / float64(metrics.TotalDays)
		metrics.AvgInlineChat_TotalEventCount = float64(aggregation.TotalInlineChat_TotalEventCount) / float64(metrics.TotalDays)
		metrics.AvgInline_AICodeLines = float64(aggregation.TotalInline_AICodeLines) / float64(metrics.TotalDays)
		metrics.AvgInline_AcceptanceCount = float64(aggregation.TotalInline_AcceptanceCount) / float64(metrics.TotalDays)
		metrics.AvgInline_SuggestionsCount = float64(aggregation.TotalInline_SuggestionsCount) / float64(metrics.TotalDays)
	}

	// 计算接受率
	totalEvents := aggregation.TotalInlineChat_AcceptanceEventCount +
		aggregation.TotalInlineChat_DismissalEventCount +
		aggregation.TotalInlineChat_RejectionEventCount

	if totalEvents > 0 {
		metrics.AcceptanceRate = float64(aggregation.TotalInlineChat_AcceptanceEventCount) / float64(totalEvents)
	}

	return metrics
}

// resolveDisplayNameForAggregation resolves display name for user metrics aggregation
func resolveDisplayNameForAggregation(userId string, identityClient UserDisplayNameResolver) string {
	// If no identity client available, use userId as fallback
	if identityClient == nil {
		return userId
	}

	// Try to resolve display name
	displayName, err := identityClient.ResolveUserDisplayName(userId)
	if err != nil {
		// Log error but continue with userId as fallback
		fmt.Printf("Failed to resolve display name for user %s during aggregation: %v\n", userId, err)
		return userId
	}

	// If display name is empty, use userId as fallback
	if displayName == "" {
		return userId
	}

	return displayName
}

var ConvertQDevUserMetricsMeta = plugin.SubTaskMeta{
	Name:             "convertQDevUserMetrics",
	EntryPoint:       ConvertQDevUserMetrics,
	EnabledByDefault: true,
	Description:      "Convert user data to metrics by each user",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}
