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

func TestQDevUserData_WithDisplayName(t *testing.T) {
	userData := QDevUserData{
		ConnectionId: 1,
		UserId:       "uuid-123",
		DisplayName:  "John Doe",
		Date:         time.Now(),
		CodeReview_FindingsCount: 5,
		Inline_AcceptanceCount:   10,
	}

	assert.Equal(t, "John Doe", userData.DisplayName)
	assert.Equal(t, "uuid-123", userData.UserId)
	assert.Equal(t, uint64(1), userData.ConnectionId)
	assert.Equal(t, 5, userData.CodeReview_FindingsCount)
	assert.Equal(t, 10, userData.Inline_AcceptanceCount)
}

func TestQDevUserData_WithFallbackDisplayName(t *testing.T) {
	userData := QDevUserData{
		ConnectionId: 1,
		UserId:       "uuid-456",
		DisplayName:  "uuid-456", // Fallback case when display name resolution fails
		Date:         time.Now(),
	}

	assert.Equal(t, "uuid-456", userData.DisplayName)
	assert.Equal(t, userData.UserId, userData.DisplayName) // Should match when fallback
}

func TestQDevUserData_EmptyDisplayName(t *testing.T) {
	userData := QDevUserData{
		ConnectionId: 1,
		UserId:       "uuid-789",
		DisplayName:  "", // Empty display name
		Date:         time.Now(),
	}

	assert.Equal(t, "", userData.DisplayName)
	assert.Equal(t, "uuid-789", userData.UserId)
	assert.NotEqual(t, userData.UserId, userData.DisplayName)
}

func TestQDevUserData_TableName(t *testing.T) {
	userData := QDevUserData{}
	assert.Equal(t, "_tool_q_dev_user_data", userData.TableName())
}

func TestQDevUserData_AllFields(t *testing.T) {
	now := time.Now()
	userData := QDevUserData{
		ConnectionId:                      1,
		UserId:                            "test-user",
		DisplayName:                       "Test User",
		Date:                              now,
		CodeReview_FindingsCount:          1,
		CodeReview_SucceededEventCount:    2,
		InlineChat_AcceptanceEventCount:   3,
		InlineChat_AcceptedLineAdditions:  4,
		InlineChat_AcceptedLineDeletions:  5,
		InlineChat_DismissalEventCount:    6,
		InlineChat_DismissedLineAdditions: 7,
		InlineChat_DismissedLineDeletions: 8,
		InlineChat_RejectedLineAdditions:  9,
		InlineChat_RejectedLineDeletions:  10,
		InlineChat_RejectionEventCount:    11,
		InlineChat_TotalEventCount:        12,
		Inline_AICodeLines:                13,
		Inline_AcceptanceCount:            14,
		Inline_SuggestionsCount:           15,
	}

	// Verify all fields are properly set
	assert.Equal(t, uint64(1), userData.ConnectionId)
	assert.Equal(t, "test-user", userData.UserId)
	assert.Equal(t, "Test User", userData.DisplayName)
	assert.Equal(t, now, userData.Date)
	assert.Equal(t, 1, userData.CodeReview_FindingsCount)
	assert.Equal(t, 2, userData.CodeReview_SucceededEventCount)
	assert.Equal(t, 3, userData.InlineChat_AcceptanceEventCount)
	assert.Equal(t, 4, userData.InlineChat_AcceptedLineAdditions)
	assert.Equal(t, 5, userData.InlineChat_AcceptedLineDeletions)
	assert.Equal(t, 6, userData.InlineChat_DismissalEventCount)
	assert.Equal(t, 7, userData.InlineChat_DismissedLineAdditions)
	assert.Equal(t, 8, userData.InlineChat_DismissedLineDeletions)
	assert.Equal(t, 9, userData.InlineChat_RejectedLineAdditions)
	assert.Equal(t, 10, userData.InlineChat_RejectedLineDeletions)
	assert.Equal(t, 11, userData.InlineChat_RejectionEventCount)
	assert.Equal(t, 12, userData.InlineChat_TotalEventCount)
	assert.Equal(t, 13, userData.Inline_AICodeLines)
	assert.Equal(t, 14, userData.Inline_AcceptanceCount)
	assert.Equal(t, 15, userData.Inline_SuggestionsCount)
}
