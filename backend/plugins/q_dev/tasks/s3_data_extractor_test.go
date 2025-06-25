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

func TestCreateUserDataWithDisplayName_AllFields(t *testing.T) {
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
	
	// Verify all metric fields
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
