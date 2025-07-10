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

func TestQDevUserDataAllMetrics(t *testing.T) {
	// Create a test user data object with all metrics
	userData := &QDevUserData{
		ConnectionId: 1,
		UserId: "test-user-id",
		Date: time.Now(),
		DisplayName: "Test User",
		
		// Set values for existing metrics
		CodeReview_FindingsCount: 10,
		CodeReview_SucceededEventCount: 11,
		InlineChat_AcceptanceEventCount: 12,
		InlineChat_AcceptedLineAdditions: 13,
		InlineChat_AcceptedLineDeletions: 14,
		InlineChat_DismissalEventCount: 15,
		InlineChat_DismissedLineAdditions: 16,
		InlineChat_DismissedLineDeletions: 17,
		InlineChat_RejectedLineAdditions: 18,
		InlineChat_RejectedLineDeletions: 19,
		InlineChat_RejectionEventCount: 20,
		InlineChat_TotalEventCount: 21,
		Inline_AICodeLines: 22,
		Inline_AcceptanceCount: 23,
		Inline_SuggestionsCount: 24,
		
		// Set values for new metrics
		Chat_AICodeLines: 25,
		Chat_MessagesInteracted: 26,
		Chat_MessagesSent: 27,
		CodeFix_AcceptanceEventCount: 28,
		CodeFix_AcceptedLines: 29,
		CodeFix_GeneratedLines: 30,
		CodeFix_GenerationEventCount: 31,
		CodeReview_FailedEventCount: 32,
		Dev_AcceptanceEventCount: 33,
		Dev_AcceptedLines: 34,
		Dev_GeneratedLines: 35,
		Dev_GenerationEventCount: 36,
		DocGeneration_AcceptedFileUpdates: 37,
		DocGeneration_AcceptedFilesCreations: 38,
		DocGeneration_AcceptedLineAdditions: 39,
		DocGeneration_AcceptedLineUpdates: 40,
		DocGeneration_EventCount: 41,
		DocGeneration_RejectedFileCreations: 42,
		DocGeneration_RejectedFileUpdates: 43,
		DocGeneration_RejectedLineAdditions: 44,
		DocGeneration_RejectedLineUpdates: 45,
		TestGeneration_AcceptedLines: 46,
		TestGeneration_AcceptedTests: 47,
		TestGeneration_EventCount: 48,
		TestGeneration_GeneratedLines: 49,
		TestGeneration_GeneratedTests: 50,
		Transformation_EventCount: 51,
		Transformation_LinesGenerated: 52,
		Transformation_LinesIngested: 53,
	}
	
	// Verify that all metrics are accessible
	// Existing metrics
	assert.Equal(t, 10, userData.CodeReview_FindingsCount)
	assert.Equal(t, 11, userData.CodeReview_SucceededEventCount)
	assert.Equal(t, 12, userData.InlineChat_AcceptanceEventCount)
	assert.Equal(t, 13, userData.InlineChat_AcceptedLineAdditions)
	assert.Equal(t, 14, userData.InlineChat_AcceptedLineDeletions)
	assert.Equal(t, 15, userData.InlineChat_DismissalEventCount)
	assert.Equal(t, 16, userData.InlineChat_DismissedLineAdditions)
	assert.Equal(t, 17, userData.InlineChat_DismissedLineDeletions)
	assert.Equal(t, 18, userData.InlineChat_RejectedLineAdditions)
	assert.Equal(t, 19, userData.InlineChat_RejectedLineDeletions)
	assert.Equal(t, 20, userData.InlineChat_RejectionEventCount)
	assert.Equal(t, 21, userData.InlineChat_TotalEventCount)
	assert.Equal(t, 22, userData.Inline_AICodeLines)
	assert.Equal(t, 23, userData.Inline_AcceptanceCount)
	assert.Equal(t, 24, userData.Inline_SuggestionsCount)
	
	// New metrics
	assert.Equal(t, 25, userData.Chat_AICodeLines)
	assert.Equal(t, 26, userData.Chat_MessagesInteracted)
	assert.Equal(t, 27, userData.Chat_MessagesSent)
	assert.Equal(t, 28, userData.CodeFix_AcceptanceEventCount)
	assert.Equal(t, 29, userData.CodeFix_AcceptedLines)
	assert.Equal(t, 30, userData.CodeFix_GeneratedLines)
	assert.Equal(t, 31, userData.CodeFix_GenerationEventCount)
	assert.Equal(t, 32, userData.CodeReview_FailedEventCount)
	assert.Equal(t, 33, userData.Dev_AcceptanceEventCount)
	assert.Equal(t, 34, userData.Dev_AcceptedLines)
	assert.Equal(t, 35, userData.Dev_GeneratedLines)
	assert.Equal(t, 36, userData.Dev_GenerationEventCount)
	assert.Equal(t, 37, userData.DocGeneration_AcceptedFileUpdates)
	assert.Equal(t, 38, userData.DocGeneration_AcceptedFilesCreations)
	assert.Equal(t, 39, userData.DocGeneration_AcceptedLineAdditions)
	assert.Equal(t, 40, userData.DocGeneration_AcceptedLineUpdates)
	assert.Equal(t, 41, userData.DocGeneration_EventCount)
	assert.Equal(t, 42, userData.DocGeneration_RejectedFileCreations)
	assert.Equal(t, 43, userData.DocGeneration_RejectedFileUpdates)
	assert.Equal(t, 44, userData.DocGeneration_RejectedLineAdditions)
	assert.Equal(t, 45, userData.DocGeneration_RejectedLineUpdates)
	assert.Equal(t, 46, userData.TestGeneration_AcceptedLines)
	assert.Equal(t, 47, userData.TestGeneration_AcceptedTests)
	assert.Equal(t, 48, userData.TestGeneration_EventCount)
	assert.Equal(t, 49, userData.TestGeneration_GeneratedLines)
	assert.Equal(t, 50, userData.TestGeneration_GeneratedTests)
	assert.Equal(t, 51, userData.Transformation_EventCount)
	assert.Equal(t, 52, userData.Transformation_LinesGenerated)
	assert.Equal(t, 53, userData.Transformation_LinesIngested)
}

func TestQDevUserDataTableName(t *testing.T) {
	userData := &QDevUserData{}
	assert.Equal(t, "_tool_q_dev_user_data", userData.TableName())
}
