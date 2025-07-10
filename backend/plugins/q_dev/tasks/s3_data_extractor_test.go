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
	"github.com/stretchr/testify/mock"

	"github.com/apache/incubator-devlake/plugins/q_dev/models"
)

// Mock Identity Client for testing
type MockIdentityClient struct {
	mock.Mock
}

func (m *MockIdentityClient) ResolveUserDisplayName(userId string) (string, error) {
	args := m.Called(userId)
	return args.String(0), args.Error(1)
}

// Ensure MockIdentityClient implements UserDisplayNameResolver
var _ UserDisplayNameResolver = (*MockIdentityClient)(nil)

func TestCreateUserDataWithDisplayName_Success(t *testing.T) {
	headers := []string{"UserId", "Date", "CodeReview_FindingsCount", "Inline_AcceptanceCount"}
	record := []string{"user-123", "2025-06-23", "5", "10"}
	fileMeta := &models.QDevS3FileMeta{
		ConnectionId: 1,
	}
	
	mockIdentityClient := &MockIdentityClient{}
	mockIdentityClient.On("ResolveUserDisplayName", "user-123").Return("John Doe", nil)
	
	userData, err := createUserDataWithDisplayName(headers, record, fileMeta, mockIdentityClient)
	
	assert.NoError(t, err)
	assert.NotNil(t, userData)
	assert.Equal(t, "user-123", userData.UserId)
	assert.Equal(t, "John Doe", userData.DisplayName)
	assert.Equal(t, uint64(1), userData.ConnectionId)
	assert.Equal(t, 5, userData.CodeReview_FindingsCount)
	assert.Equal(t, 10, userData.Inline_AcceptanceCount)
	
	mockIdentityClient.AssertExpectations(t)
}

func TestCreateUserDataWithDisplayName_FallbackToUUID(t *testing.T) {
	headers := []string{"UserId", "Date"}
	record := []string{"user-456", "2025-06-23"}
	fileMeta := &models.QDevS3FileMeta{
		ConnectionId: 1,
	}
	
	mockIdentityClient := &MockIdentityClient{}
	mockIdentityClient.On("ResolveUserDisplayName", "user-456").Return("user-456", assert.AnError)
	
	userData, err := createUserDataWithDisplayName(headers, record, fileMeta, mockIdentityClient)
	
	assert.NoError(t, err)
	assert.NotNil(t, userData)
	assert.Equal(t, "user-456", userData.UserId)
	assert.Equal(t, "user-456", userData.DisplayName) // Should fallback to UUID
	
	mockIdentityClient.AssertExpectations(t)
}

func TestCreateUserDataWithDisplayName_NoIdentityClient(t *testing.T) {
	headers := []string{"UserId", "Date"}
	record := []string{"user-789", "2025-06-23"}
	fileMeta := &models.QDevS3FileMeta{
		ConnectionId: 1,
	}
	
	userData, err := createUserDataWithDisplayName(headers, record, fileMeta, nil)
	
	assert.NoError(t, err)
	assert.NotNil(t, userData)
	assert.Equal(t, "user-789", userData.UserId)
	assert.Equal(t, "user-789", userData.DisplayName) // Should use UUID when no client
}

func TestCreateUserDataWithDisplayName_EmptyDisplayName(t *testing.T) {
	headers := []string{"UserId", "Date"}
	record := []string{"user-empty", "2025-06-23"}
	fileMeta := &models.QDevS3FileMeta{
		ConnectionId: 1,
	}
	
	mockIdentityClient := &MockIdentityClient{}
	mockIdentityClient.On("ResolveUserDisplayName", "user-empty").Return("", nil)
	
	userData, err := createUserDataWithDisplayName(headers, record, fileMeta, mockIdentityClient)
	
	assert.NoError(t, err)
	assert.NotNil(t, userData)
	assert.Equal(t, "user-empty", userData.UserId)
	assert.Equal(t, "user-empty", userData.DisplayName) // Should fallback when empty
	
	mockIdentityClient.AssertExpectations(t)
}

func TestCreateUserDataWithDisplayName_AllExistingMetrics(t *testing.T) {
	headers := []string{
		"UserId", "Date", "CodeReview_FindingsCount", "CodeReview_SucceededEventCount",
		"InlineChat_AcceptanceEventCount", "InlineChat_AcceptedLineAdditions",
		"InlineChat_AcceptedLineDeletions", "InlineChat_DismissalEventCount",
		"InlineChat_DismissedLineAdditions", "InlineChat_DismissedLineDeletions",
		"InlineChat_RejectedLineAdditions", "InlineChat_RejectedLineDeletions",
		"InlineChat_RejectionEventCount", "InlineChat_TotalEventCount",
		"Inline_AICodeLines", "Inline_AcceptanceCount", "Inline_SuggestionsCount",
	}
	record := []string{
		"test-user", "2025-06-23", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15",
	}
	fileMeta := &models.QDevS3FileMeta{
		ConnectionId: 123,
	}
	
	mockIdentityClient := &MockIdentityClient{}
	mockIdentityClient.On("ResolveUserDisplayName", "test-user").Return("Test User", nil)
	
	userData, err := createUserDataWithDisplayName(headers, record, fileMeta, mockIdentityClient)
	
	assert.NoError(t, err)
	assert.NotNil(t, userData)
	
	// Verify basic fields
	assert.Equal(t, "test-user", userData.UserId)
	assert.Equal(t, "Test User", userData.DisplayName)
	assert.Equal(t, uint64(123), userData.ConnectionId)
	
	// Verify date parsing
	expectedDate, _ := time.Parse("2006-01-02", "2025-06-23")
	assert.Equal(t, expectedDate, userData.Date)
	
	// Verify all existing metric fields
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
	
	mockIdentityClient.AssertExpectations(t)
}

func TestCreateUserDataWithDisplayName_AllNewMetrics(t *testing.T) {
	headers := []string{
		"UserId", "Date",
		"Chat_AICodeLines", "Chat_MessagesInteracted", "Chat_MessagesSent",
		"CodeFix_AcceptanceEventCount", "CodeFix_AcceptedLines", "CodeFix_GeneratedLines", "CodeFix_GenerationEventCount",
		"CodeReview_FailedEventCount",
		"Dev_AcceptanceEventCount", "Dev_AcceptedLines", "Dev_GeneratedLines", "Dev_GenerationEventCount",
		"DocGeneration_AcceptedFileUpdates", "DocGeneration_AcceptedFilesCreations", "DocGeneration_AcceptedLineAdditions",
		"DocGeneration_AcceptedLineUpdates", "DocGeneration_EventCount", "DocGeneration_RejectedFileCreations",
		"DocGeneration_RejectedFileUpdates", "DocGeneration_RejectedLineAdditions", "DocGeneration_RejectedLineUpdates",
		"TestGeneration_AcceptedLines", "TestGeneration_AcceptedTests", "TestGeneration_EventCount",
		"TestGeneration_GeneratedLines", "TestGeneration_GeneratedTests",
		"Transformation_EventCount", "Transformation_LinesGenerated", "Transformation_LinesIngested",
	}
	
	record := []string{
		"test-user", "2025-06-23",
		"101", "102", "103", "104", "105", "106", "107", "108", "109", "110",
		"111", "112", "113", "114", "115", "116", "117", "118", "119", "120",
		"121", "122", "123", "124", "125", "126", "127", "128", "129",
	}
	
	fileMeta := &models.QDevS3FileMeta{
		ConnectionId: 123,
	}
	
	mockIdentityClient := &MockIdentityClient{}
	mockIdentityClient.On("ResolveUserDisplayName", "test-user").Return("Test User", nil)
	
	userData, err := createUserDataWithDisplayName(headers, record, fileMeta, mockIdentityClient)
	
	assert.NoError(t, err)
	assert.NotNil(t, userData)
	
	// Verify basic fields
	assert.Equal(t, "test-user", userData.UserId)
	assert.Equal(t, "Test User", userData.DisplayName)
	
	// Verify all new metric fields
	assert.Equal(t, 101, userData.Chat_AICodeLines)
	assert.Equal(t, 102, userData.Chat_MessagesInteracted)
	assert.Equal(t, 103, userData.Chat_MessagesSent)
	assert.Equal(t, 104, userData.CodeFix_AcceptanceEventCount)
	assert.Equal(t, 105, userData.CodeFix_AcceptedLines)
	assert.Equal(t, 106, userData.CodeFix_GeneratedLines)
	assert.Equal(t, 107, userData.CodeFix_GenerationEventCount)
	assert.Equal(t, 108, userData.CodeReview_FailedEventCount)
	assert.Equal(t, 109, userData.Dev_AcceptanceEventCount)
	assert.Equal(t, 110, userData.Dev_AcceptedLines)
	assert.Equal(t, 111, userData.Dev_GeneratedLines)
	assert.Equal(t, 112, userData.Dev_GenerationEventCount)
	assert.Equal(t, 113, userData.DocGeneration_AcceptedFileUpdates)
	assert.Equal(t, 114, userData.DocGeneration_AcceptedFilesCreations)
	assert.Equal(t, 115, userData.DocGeneration_AcceptedLineAdditions)
	assert.Equal(t, 116, userData.DocGeneration_AcceptedLineUpdates)
	assert.Equal(t, 117, userData.DocGeneration_EventCount)
	assert.Equal(t, 118, userData.DocGeneration_RejectedFileCreations)
	assert.Equal(t, 119, userData.DocGeneration_RejectedFileUpdates)
	assert.Equal(t, 120, userData.DocGeneration_RejectedLineAdditions)
	assert.Equal(t, 121, userData.DocGeneration_RejectedLineUpdates)
	assert.Equal(t, 122, userData.TestGeneration_AcceptedLines)
	assert.Equal(t, 123, userData.TestGeneration_AcceptedTests)
	assert.Equal(t, 124, userData.TestGeneration_EventCount)
	assert.Equal(t, 125, userData.TestGeneration_GeneratedLines)
	assert.Equal(t, 126, userData.TestGeneration_GeneratedTests)
	assert.Equal(t, 127, userData.Transformation_EventCount)
	assert.Equal(t, 128, userData.Transformation_LinesGenerated)
	assert.Equal(t, 129, userData.Transformation_LinesIngested)
	
	mockIdentityClient.AssertExpectations(t)
}

func TestCreateUserDataWithDisplayName_MissingMetrics(t *testing.T) {
	// Only provide a few metrics in the CSV
	headers := []string{"UserId", "Date", "CodeReview_FindingsCount", "Chat_AICodeLines"}
	record := []string{"test-user", "2025-06-23", "42", "99"}
	
	fileMeta := &models.QDevS3FileMeta{
		ConnectionId: 123,
	}
	
	mockIdentityClient := &MockIdentityClient{}
	mockIdentityClient.On("ResolveUserDisplayName", "test-user").Return("Test User", nil)
	
	userData, err := createUserDataWithDisplayName(headers, record, fileMeta, mockIdentityClient)
	
	assert.NoError(t, err)
	assert.NotNil(t, userData)
	
	// Verify provided metrics are set correctly
	assert.Equal(t, 42, userData.CodeReview_FindingsCount)
	assert.Equal(t, 99, userData.Chat_AICodeLines)
	
	// Verify missing metrics are set to 0
	assert.Equal(t, 0, userData.CodeReview_SucceededEventCount)
	assert.Equal(t, 0, userData.InlineChat_AcceptanceEventCount)
	assert.Equal(t, 0, userData.Chat_MessagesInteracted)
	assert.Equal(t, 0, userData.TestGeneration_AcceptedTests)
	assert.Equal(t, 0, userData.Transformation_LinesIngested)
	
	mockIdentityClient.AssertExpectations(t)
}

func TestCreateUserDataWithDisplayName_InvalidMetricValues(t *testing.T) {
	headers := []string{
		"UserId", "Date", "CodeReview_FindingsCount", "Chat_AICodeLines", 
		"InlineChat_AcceptanceEventCount", "TestGeneration_AcceptedTests",
	}
	record := []string{"test-user", "2025-06-23", "42", "not-a-number", "abc", ""}
	
	fileMeta := &models.QDevS3FileMeta{
		ConnectionId: 123,
	}
	
	mockIdentityClient := &MockIdentityClient{}
	mockIdentityClient.On("ResolveUserDisplayName", "test-user").Return("Test User", nil)
	
	userData, err := createUserDataWithDisplayName(headers, record, fileMeta, mockIdentityClient)
	
	assert.NoError(t, err)
	assert.NotNil(t, userData)
	
	// Verify valid metric is set correctly
	assert.Equal(t, 42, userData.CodeReview_FindingsCount)
	
	// Verify invalid metrics are set to 0
	assert.Equal(t, 0, userData.Chat_AICodeLines)
	assert.Equal(t, 0, userData.InlineChat_AcceptanceEventCount)
	assert.Equal(t, 0, userData.TestGeneration_AcceptedTests)
	
	mockIdentityClient.AssertExpectations(t)
}

func TestCreateUserDataWithDisplayName_MissingUserId(t *testing.T) {
	headers := []string{"Date", "CodeReview_FindingsCount"}
	record := []string{"2025-06-23", "5"}
	fileMeta := &models.QDevS3FileMeta{
		ConnectionId: 1,
	}
	
	userData, err := createUserDataWithDisplayName(headers, record, fileMeta, nil)
	
	assert.Error(t, err)
	assert.Nil(t, userData)
	assert.Contains(t, err.Error(), "UserId not found")
}

func TestCreateUserDataWithDisplayName_MissingDate(t *testing.T) {
	headers := []string{"UserId", "CodeReview_FindingsCount"}
	record := []string{"user-123", "5"}
	fileMeta := &models.QDevS3FileMeta{
		ConnectionId: 1,
	}
	
	userData, err := createUserDataWithDisplayName(headers, record, fileMeta, nil)
	
	assert.Error(t, err)
	assert.Nil(t, userData)
	assert.Contains(t, err.Error(), "Date not found")
}

func TestParseDate(t *testing.T) {
	testCases := []struct {
		dateStr      string
		expectedDate time.Time
		expectError  bool
	}{
		{"2025-07-10", time.Date(2025, 7, 10, 0, 0, 0, 0, time.UTC), false},
		{"2025/07/10", time.Date(2025, 7, 10, 0, 0, 0, 0, time.UTC), false},
		{"07/10/2025", time.Date(2025, 7, 10, 0, 0, 0, 0, time.UTC), false},
		{"07-10-2025", time.Date(2025, 7, 10, 0, 0, 0, 0, time.UTC), false},
		{"2025-07-10T15:04:05Z", time.Date(2025, 7, 10, 15, 4, 5, 0, time.UTC), false},
		{"invalid-date", time.Time{}, true},
	}
	
	for _, tc := range testCases {
		date, err := parseDate(tc.dateStr)
		
		if tc.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedDate, date)
		}
	}
}

func TestParseInt(t *testing.T) {
	fieldMap := map[string]string{
		"ValidInt": "42",
		"ZeroInt": "0",
		"NegativeInt": "-10",
		"InvalidInt": "not-a-number",
		"EmptyString": "",
	}
	
	assert.Equal(t, 42, parseInt(fieldMap, "ValidInt"))
	assert.Equal(t, 0, parseInt(fieldMap, "ZeroInt"))
	assert.Equal(t, -10, parseInt(fieldMap, "NegativeInt"))
	assert.Equal(t, 0, parseInt(fieldMap, "InvalidInt"))
	assert.Equal(t, 0, parseInt(fieldMap, "EmptyString"))
	assert.Equal(t, 0, parseInt(fieldMap, "NonExistentField"))
}
