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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserMetricsAggregationWithDisplayName_SingleUser(t *testing.T) {
	aggregation := &UserMetricsAggregationWithDisplayName{
		ConnectionId: 1,
		UserId:       "user-123",
		DisplayName:  "John Doe",
		FirstDate:    time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC),
		LastDate:     time.Date(2025, 6, 23, 0, 0, 0, 0, time.UTC),
		DataCount:    4,
		TotalCodeReview_FindingsCount:        20,
		TotalInlineChat_AcceptanceEventCount: 40,
		TotalInline_AcceptanceCount:          60,
		TotalInline_SuggestionsCount:         80,
	}

	metrics := aggregation.ToUserMetrics()

	assert.Equal(t, uint64(1), metrics.ConnectionId)
	assert.Equal(t, "user-123", metrics.UserId)
	assert.Equal(t, "John Doe", metrics.DisplayName)
	assert.Equal(t, 4, metrics.TotalDays) // 6/20 to 6/23 = 4 days
	assert.Equal(t, 20, metrics.TotalCodeReview_FindingsCount)
	assert.Equal(t, 40, metrics.TotalInlineChat_AcceptanceEventCount)
	assert.Equal(t, 60, metrics.TotalInline_AcceptanceCount)
	assert.Equal(t, 80, metrics.TotalInline_SuggestionsCount)

	// Test averages
	assert.Equal(t, 5.0, metrics.AvgCodeReview_FindingsCount)        // 20/4
	assert.Equal(t, 10.0, metrics.AvgInlineChat_AcceptanceEventCount) // 40/4
	assert.Equal(t, 15.0, metrics.AvgInline_AcceptanceCount)         // 60/4
	assert.Equal(t, 20.0, metrics.AvgInline_SuggestionsCount)        // 80/4
}

func TestUserMetricsAggregationWithDisplayName_FallbackDisplayName(t *testing.T) {
	aggregation := &UserMetricsAggregationWithDisplayName{
		ConnectionId: 1,
		UserId:       "user-456",
		DisplayName:  "user-456", // Fallback case
		FirstDate:    time.Date(2025, 6, 23, 0, 0, 0, 0, time.UTC),
		LastDate:     time.Date(2025, 6, 23, 0, 0, 0, 0, time.UTC),
		DataCount:    1,
	}

	metrics := aggregation.ToUserMetrics()

	assert.Equal(t, "user-456", metrics.UserId)
	assert.Equal(t, "user-456", metrics.DisplayName)
	assert.Equal(t, metrics.UserId, metrics.DisplayName) // Should match when fallback
	assert.Equal(t, 1, metrics.TotalDays) // Same day = 1 day
}

func TestUserMetricsAggregationWithDisplayName_AcceptanceRateCalculation(t *testing.T) {
	aggregation := &UserMetricsAggregationWithDisplayName{
		ConnectionId: 1,
		UserId:       "user-789",
		DisplayName:  "Jane Smith",
		FirstDate:    time.Date(2025, 6, 23, 0, 0, 0, 0, time.UTC),
		LastDate:     time.Date(2025, 6, 23, 0, 0, 0, 0, time.UTC),
		DataCount:    1,
		TotalInlineChat_AcceptanceEventCount: 80,  // Accepted
		TotalInlineChat_DismissalEventCount:  15,  // Dismissed
		TotalInlineChat_RejectionEventCount:  5,   // Rejected
		// Total events = 80 + 15 + 5 = 100
		// Acceptance rate = 80/100 = 0.8
	}

	metrics := aggregation.ToUserMetrics()

	assert.Equal(t, "Jane Smith", metrics.DisplayName)
	assert.Equal(t, 80, metrics.TotalInlineChat_AcceptanceEventCount)
	assert.Equal(t, 15, metrics.TotalInlineChat_DismissalEventCount)
	assert.Equal(t, 5, metrics.TotalInlineChat_RejectionEventCount)
	assert.Equal(t, 0.8, metrics.AcceptanceRate)
}

func TestUserMetricsAggregationWithDisplayName_ZeroAcceptanceRate(t *testing.T) {
	aggregation := &UserMetricsAggregationWithDisplayName{
		ConnectionId: 1,
		UserId:       "user-zero",
		DisplayName:  "Zero User",
		FirstDate:    time.Date(2025, 6, 23, 0, 0, 0, 0, time.UTC),
		LastDate:     time.Date(2025, 6, 23, 0, 0, 0, 0, time.UTC),
		DataCount:    1,
		// No events = acceptance rate should be 0
	}

	metrics := aggregation.ToUserMetrics()

	assert.Equal(t, "Zero User", metrics.DisplayName)
	assert.Equal(t, 0.0, metrics.AcceptanceRate)
}

func TestUserMetricsAggregationWithDisplayName_AllFields(t *testing.T) {
	firstDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	lastDate := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)
	
	aggregation := &UserMetricsAggregationWithDisplayName{
		ConnectionId: 123,
		UserId:       "test-user",
		DisplayName:  "Test User",
		FirstDate:    firstDate,
		LastDate:     lastDate,
		DataCount:    30,

		// Set all total fields
		TotalCodeReview_FindingsCount:          300,
		TotalCodeReview_SucceededEventCount:    270,
		TotalInlineChat_AcceptanceEventCount:   240,
		TotalInlineChat_AcceptedLineAdditions:  210,
		TotalInlineChat_AcceptedLineDeletions:  180,
		TotalInlineChat_DismissalEventCount:    150,
		TotalInlineChat_DismissedLineAdditions: 120,
		TotalInlineChat_DismissedLineDeletions: 90,
		TotalInlineChat_RejectedLineAdditions:  60,
		TotalInlineChat_RejectedLineDeletions:  30,
		TotalInlineChat_RejectionEventCount:    15,
		TotalInlineChat_TotalEventCount:        600,
		TotalInline_AICodeLines:                3000,
		TotalInline_AcceptanceCount:            450,
		TotalInline_SuggestionsCount:           540,
	}

	metrics := aggregation.ToUserMetrics()

	// Verify basic fields
	assert.Equal(t, uint64(123), metrics.ConnectionId)
	assert.Equal(t, "test-user", metrics.UserId)
	assert.Equal(t, "Test User", metrics.DisplayName)
	assert.Equal(t, firstDate, metrics.FirstDate)
	assert.Equal(t, lastDate, metrics.LastDate)
	assert.Equal(t, 30, metrics.TotalDays) // June 1-30 = 30 days

	// Verify all total fields
	assert.Equal(t, 300, metrics.TotalCodeReview_FindingsCount)
	assert.Equal(t, 270, metrics.TotalCodeReview_SucceededEventCount)
	assert.Equal(t, 240, metrics.TotalInlineChat_AcceptanceEventCount)
	assert.Equal(t, 210, metrics.TotalInlineChat_AcceptedLineAdditions)
	assert.Equal(t, 180, metrics.TotalInlineChat_AcceptedLineDeletions)
	assert.Equal(t, 150, metrics.TotalInlineChat_DismissalEventCount)
	assert.Equal(t, 120, metrics.TotalInlineChat_DismissedLineAdditions)
	assert.Equal(t, 90, metrics.TotalInlineChat_DismissedLineDeletions)
	assert.Equal(t, 60, metrics.TotalInlineChat_RejectedLineAdditions)
	assert.Equal(t, 30, metrics.TotalInlineChat_RejectedLineDeletions)
	assert.Equal(t, 15, metrics.TotalInlineChat_RejectionEventCount)
	assert.Equal(t, 600, metrics.TotalInlineChat_TotalEventCount)
	assert.Equal(t, 3000, metrics.TotalInline_AICodeLines)
	assert.Equal(t, 450, metrics.TotalInline_AcceptanceCount)
	assert.Equal(t, 540, metrics.TotalInline_SuggestionsCount)

	// Verify average fields (all divided by 30 days)
	assert.Equal(t, 10.0, metrics.AvgCodeReview_FindingsCount)        // 300/30
	assert.Equal(t, 9.0, metrics.AvgCodeReview_SucceededEventCount)   // 270/30
	assert.Equal(t, 8.0, metrics.AvgInlineChat_AcceptanceEventCount)  // 240/30
	assert.Equal(t, 20.0, metrics.AvgInlineChat_TotalEventCount)      // 600/30
	assert.Equal(t, 100.0, metrics.AvgInline_AICodeLines)            // 3000/30
	assert.Equal(t, 15.0, metrics.AvgInline_AcceptanceCount)         // 450/30
	assert.Equal(t, 18.0, metrics.AvgInline_SuggestionsCount)        // 540/30

	// Verify acceptance rate: 240 / (240 + 150 + 15) = 240/405 ≈ 0.593
	expectedAcceptanceRate := 240.0 / (240.0 + 150.0 + 15.0)
	assert.InDelta(t, expectedAcceptanceRate, metrics.AcceptanceRate, 0.001)
}

func TestResolveDisplayNameForAggregation_Success(t *testing.T) {
	mockIdentityClient := &MockIdentityClient{}
	mockIdentityClient.On("ResolveUserDisplayName", "user-123").Return("John Doe", nil)

	displayName := resolveDisplayNameForAggregation("user-123", mockIdentityClient)
	assert.Equal(t, "John Doe", displayName)

	mockIdentityClient.AssertExpectations(t)
}

func TestResolveDisplayNameForAggregation_NoClient(t *testing.T) {
	displayName := resolveDisplayNameForAggregation("user-456", nil)
	assert.Equal(t, "user-456", displayName)
}

func TestResolveDisplayNameForAggregation_Error(t *testing.T) {
	mockIdentityClient := &MockIdentityClient{}
	mockIdentityClient.On("ResolveUserDisplayName", "user-error").Return("user-error", assert.AnError)

	displayName := resolveDisplayNameForAggregation("user-error", mockIdentityClient)
	assert.Equal(t, "user-error", displayName) // Should fallback to UUID

	mockIdentityClient.AssertExpectations(t)
}
