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

package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQDevUserMetrics_WithDisplayName(t *testing.T) {
	userMetrics := QDevUserMetrics{
		ConnectionId: 1,
		UserId:       "uuid-123",
		DisplayName:  "John Doe",
		FirstDate:    time.Now().AddDate(0, 0, -30),
		LastDate:     time.Now(),
		TotalDays:    30,
		TotalCodeReview_FindingsCount: 50,
		AcceptanceRate:                0.85,
	}

	assert.Equal(t, "John Doe", userMetrics.DisplayName)
	assert.Equal(t, "uuid-123", userMetrics.UserId)
	assert.Equal(t, uint64(1), userMetrics.ConnectionId)
	assert.Equal(t, 50, userMetrics.TotalCodeReview_FindingsCount)
	assert.Equal(t, 0.85, userMetrics.AcceptanceRate)
}

func TestQDevUserMetrics_WithFallbackDisplayName(t *testing.T) {
	userMetrics := QDevUserMetrics{
		ConnectionId: 1,
		UserId:       "uuid-456",
		DisplayName:  "uuid-456", // Fallback case when display name resolution fails
		TotalDays:    15,
	}

	assert.Equal(t, "uuid-456", userMetrics.DisplayName)
	assert.Equal(t, userMetrics.UserId, userMetrics.DisplayName) // Should match when fallback
	assert.Equal(t, 15, userMetrics.TotalDays)
}

func TestQDevUserMetrics_EmptyDisplayName(t *testing.T) {
	userMetrics := QDevUserMetrics{
		ConnectionId: 1,
		UserId:       "uuid-789",
		DisplayName:  "", // Empty display name
		TotalDays:    5,
	}

	assert.Equal(t, "", userMetrics.DisplayName)
	assert.Equal(t, "uuid-789", userMetrics.UserId)
	assert.NotEqual(t, userMetrics.UserId, userMetrics.DisplayName)
}

func TestQDevUserMetrics_TableName(t *testing.T) {
	userMetrics := QDevUserMetrics{}
	assert.Equal(t, "_tool_q_dev_user_metrics", userMetrics.TableName())
}

func TestQDevUserMetrics_AllFields(t *testing.T) {
	firstDate := time.Now().AddDate(0, 0, -30)
	lastDate := time.Now()
	
	userMetrics := QDevUserMetrics{
		ConnectionId: 1,
		UserId:       "test-user",
		DisplayName:  "Test User",
		FirstDate:    firstDate,
		LastDate:     lastDate,
		TotalDays:    30,

		// 聚合指标
		TotalCodeReview_FindingsCount:          100,
		TotalCodeReview_SucceededEventCount:    90,
		TotalInlineChat_AcceptanceEventCount:   80,
		TotalInlineChat_AcceptedLineAdditions:  70,
		TotalInlineChat_AcceptedLineDeletions:  60,
		TotalInlineChat_DismissalEventCount:    50,
		TotalInlineChat_DismissedLineAdditions: 40,
		TotalInlineChat_DismissedLineDeletions: 30,
		TotalInlineChat_RejectedLineAdditions:  20,
		TotalInlineChat_RejectedLineDeletions:  10,
		TotalInlineChat_RejectionEventCount:    5,
		TotalInlineChat_TotalEventCount:        200,
		TotalInline_AICodeLines:                1000,
		TotalInline_AcceptanceCount:            150,
		TotalInline_SuggestionsCount:           180,

		// 平均指标
		AvgCodeReview_FindingsCount:        3.33,
		AvgCodeReview_SucceededEventCount:  3.0,
		AvgInlineChat_AcceptanceEventCount: 2.67,
		AvgInlineChat_TotalEventCount:      6.67,
		AvgInline_AICodeLines:              33.33,
		AvgInline_AcceptanceCount:          5.0,
		AvgInline_SuggestionsCount:         6.0,

		// 接受率指标
		AcceptanceRate: 0.83,
	}

	// Verify all fields are properly set
	assert.Equal(t, uint64(1), userMetrics.ConnectionId)
	assert.Equal(t, "test-user", userMetrics.UserId)
	assert.Equal(t, "Test User", userMetrics.DisplayName)
	assert.Equal(t, firstDate, userMetrics.FirstDate)
	assert.Equal(t, lastDate, userMetrics.LastDate)
	assert.Equal(t, 30, userMetrics.TotalDays)

	// Test aggregated metrics
	assert.Equal(t, 100, userMetrics.TotalCodeReview_FindingsCount)
	assert.Equal(t, 90, userMetrics.TotalCodeReview_SucceededEventCount)
	assert.Equal(t, 80, userMetrics.TotalInlineChat_AcceptanceEventCount)
	assert.Equal(t, 70, userMetrics.TotalInlineChat_AcceptedLineAdditions)
	assert.Equal(t, 60, userMetrics.TotalInlineChat_AcceptedLineDeletions)
	assert.Equal(t, 50, userMetrics.TotalInlineChat_DismissalEventCount)
	assert.Equal(t, 40, userMetrics.TotalInlineChat_DismissedLineAdditions)
	assert.Equal(t, 30, userMetrics.TotalInlineChat_DismissedLineDeletions)
	assert.Equal(t, 20, userMetrics.TotalInlineChat_RejectedLineAdditions)
	assert.Equal(t, 10, userMetrics.TotalInlineChat_RejectedLineDeletions)
	assert.Equal(t, 5, userMetrics.TotalInlineChat_RejectionEventCount)
	assert.Equal(t, 200, userMetrics.TotalInlineChat_TotalEventCount)
	assert.Equal(t, 1000, userMetrics.TotalInline_AICodeLines)
	assert.Equal(t, 150, userMetrics.TotalInline_AcceptanceCount)
	assert.Equal(t, 180, userMetrics.TotalInline_SuggestionsCount)

	// Test average metrics
	assert.Equal(t, 3.33, userMetrics.AvgCodeReview_FindingsCount)
	assert.Equal(t, 3.0, userMetrics.AvgCodeReview_SucceededEventCount)
	assert.Equal(t, 2.67, userMetrics.AvgInlineChat_AcceptanceEventCount)
	assert.Equal(t, 6.67, userMetrics.AvgInlineChat_TotalEventCount)
	assert.Equal(t, 33.33, userMetrics.AvgInline_AICodeLines)
	assert.Equal(t, 5.0, userMetrics.AvgInline_AcceptanceCount)
	assert.Equal(t, 6.0, userMetrics.AvgInline_SuggestionsCount)

	// Test acceptance rate
	assert.Equal(t, 0.83, userMetrics.AcceptanceRate)
}
