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
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/identitystore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/apache/incubator-devlake/plugins/q_dev/models"
)

// Mock IdentityStore interface for testing
type MockIdentityStoreAPI struct {
	mock.Mock
}

func (m *MockIdentityStoreAPI) DescribeUser(input *identitystore.DescribeUserInput) (*identitystore.DescribeUserOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*identitystore.DescribeUserOutput), args.Error(1)
}

func TestNewQDevIdentityClient_Success(t *testing.T) {
	connection := &models.QDevConnection{
		QDevConn: models.QDevConn{
			AccessKeyId:         "test-key",
			SecretAccessKey:     "test-secret",
			IdentityStoreId:     "d-1234567890",
			IdentityStoreRegion: "us-west-2",
		},
	}

	client, err := NewQDevIdentityClient(connection)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "d-1234567890", client.StoreId)
	assert.Equal(t, "us-west-2", client.Region)
}

func TestNewQDevIdentityClient_EmptyIdentityStoreId(t *testing.T) {
	connection := &models.QDevConnection{
		QDevConn: models.QDevConn{
			AccessKeyId:         "test-key",
			SecretAccessKey:     "test-secret",
			IdentityStoreId:     "", // Empty identity store ID
			IdentityStoreRegion: "us-west-2",
		},
	}

	client, err := NewQDevIdentityClient(connection)
	assert.NoError(t, err)
	assert.Nil(t, client) // Should return nil when no identity store configured
}

func TestNewQDevIdentityClient_EmptyIdentityStoreRegion(t *testing.T) {
	connection := &models.QDevConnection{
		QDevConn: models.QDevConn{
			AccessKeyId:         "test-key",
			SecretAccessKey:     "test-secret",
			IdentityStoreId:     "d-1234567890",
			IdentityStoreRegion: "", // Empty identity store region
		},
	}

	client, err := NewQDevIdentityClient(connection)
	assert.NoError(t, err)
	assert.Nil(t, client) // Should return nil when no region configured
}

func TestQDevIdentityClient_ResolveUserDisplayName_Success(t *testing.T) {
	mockAPI := &MockIdentityStoreAPI{}
	client := &QDevIdentityClient{
		IdentityStore: mockAPI,
		StoreId:       "d-1234567890",
		Region:        "us-west-2",
	}

	displayName := "John Doe"
	mockAPI.On("DescribeUser", mock.AnythingOfType("*identitystore.DescribeUserInput")).Return(
		&identitystore.DescribeUserOutput{
			DisplayName: &displayName,
		}, nil)

	result, err := client.ResolveUserDisplayName("user-123")
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", result)

	mockAPI.AssertExpectations(t)
}

func TestQDevIdentityClient_ResolveUserDisplayName_NoDisplayName(t *testing.T) {
	mockAPI := &MockIdentityStoreAPI{}
	client := &QDevIdentityClient{
		IdentityStore: mockAPI,
		StoreId:       "d-1234567890",
		Region:        "us-west-2",
	}

	// Return output with nil DisplayName
	mockAPI.On("DescribeUser", mock.AnythingOfType("*identitystore.DescribeUserInput")).Return(
		&identitystore.DescribeUserOutput{
			DisplayName: nil,
		}, nil)

	result, err := client.ResolveUserDisplayName("user-123")
	assert.NoError(t, err)
	assert.Equal(t, "user-123", result) // Should fallback to UUID

	mockAPI.AssertExpectations(t)
}

func TestQDevIdentityClient_ResolveUserDisplayName_EmptyDisplayName(t *testing.T) {
	mockAPI := &MockIdentityStoreAPI{}
	client := &QDevIdentityClient{
		IdentityStore: mockAPI,
		StoreId:       "d-1234567890",
		Region:        "us-west-2",
	}

	emptyName := ""
	mockAPI.On("DescribeUser", mock.AnythingOfType("*identitystore.DescribeUserInput")).Return(
		&identitystore.DescribeUserOutput{
			DisplayName: &emptyName,
		}, nil)

	result, err := client.ResolveUserDisplayName("user-123")
	assert.NoError(t, err)
	assert.Equal(t, "user-123", result) // Should fallback to UUID when empty

	mockAPI.AssertExpectations(t)
}

func TestQDevIdentityClient_ResolveUserDisplayName_APIError(t *testing.T) {
	mockAPI := &MockIdentityStoreAPI{}
	client := &QDevIdentityClient{
		IdentityStore: mockAPI,
		StoreId:       "d-1234567890",
		Region:        "us-west-2",
	}

	mockAPI.On("DescribeUser", mock.AnythingOfType("*identitystore.DescribeUserInput")).Return(
		nil, errors.New("user not found"))

	result, err := client.ResolveUserDisplayName("user-123")
	assert.Error(t, err)
	assert.Equal(t, "user-123", result) // Should fallback to UUID on error
	assert.Contains(t, err.Error(), "user not found")

	mockAPI.AssertExpectations(t)
}

func TestQDevIdentityClient_ResolveUserDisplayName_InputValidation(t *testing.T) {
	mockAPI := &MockIdentityStoreAPI{}
	client := &QDevIdentityClient{
		IdentityStore: mockAPI,
		StoreId:       "d-1234567890",
		Region:        "us-west-2",
	}

	displayName := "Jane Smith"
	mockAPI.On("DescribeUser", mock.MatchedBy(func(input *identitystore.DescribeUserInput) bool {
		// Verify the input parameters are correctly set
		return *input.IdentityStoreId == "d-1234567890" && *input.UserId == "test-user-456"
	})).Return(
		&identitystore.DescribeUserOutput{
			DisplayName: &displayName,
		}, nil)

	result, err := client.ResolveUserDisplayName("test-user-456")
	assert.NoError(t, err)
	assert.Equal(t, "Jane Smith", result)

	mockAPI.AssertExpectations(t)
}
