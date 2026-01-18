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

package token

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/impls/logruslog"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRoundTripper struct {
	mock.Mock
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestNeedsRefresh(t *testing.T) {
	tp := &TokenProvider{
		conn: &models.GithubConnection{
			GithubConn: models.GithubConn{
				RefreshToken: "refresh_token",
			},
		},
	}

	// Not expired, outside buffer
	tp.conn.TokenExpiresAt = time.Now().Add(10 * time.Minute)
	assert.False(t, tp.needsRefresh())

	// Inside buffer
	tp.conn.TokenExpiresAt = time.Now().Add(1 * time.Minute)
	assert.True(t, tp.needsRefresh())

	// Expired
	tp.conn.TokenExpiresAt = time.Now().Add(-1 * time.Minute)
	assert.True(t, tp.needsRefresh())

	// No refresh token
	tp.conn.RefreshToken = ""
	assert.False(t, tp.needsRefresh())
}

func TestTokenProviderConcurrency(t *testing.T) {
	mockRT := new(MockRoundTripper)
	client := &http.Client{Transport: mockRT}

	conn := &models.GithubConnection{
		GithubConn: models.GithubConn{
			RefreshToken:   "refresh_token",
			TokenExpiresAt: time.Now().Add(-1 * time.Minute), // Expired
			GithubAppKey: models.GithubAppKey{
				AppKey: api.AppKey{
					AppId:     "123",
					SecretKey: "secret",
				},
			},
		},
	}

	logger, _ := logruslog.NewDefaultLogger(logrus.New())
	tp := NewTokenProvider(conn, nil, client, logger)

	// Mock response for refresh
	respBody := `{"access_token":"new_token","refresh_token":"new_refresh_token","expires_in":3600,"refresh_token_expires_in":3600}`
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(respBody)),
	}

	// Expect exactly one call
	mockRT.On("RoundTrip", mock.Anything).Return(resp, nil).Once()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token, err := tp.GetToken()
			assert.NoError(t, err)
			assert.Equal(t, "new_token", token)
		}()
	}
	wg.Wait()

	mockRT.AssertExpectations(t)
}

func TestConfigurableBuffer(t *testing.T) {
	os.Setenv("GITHUB_TOKEN_REFRESH_BUFFER_MINUTES", "10")
	defer os.Unsetenv("GITHUB_TOKEN_REFRESH_BUFFER_MINUTES")

	tp := &TokenProvider{
		conn: &models.GithubConnection{
			GithubConn: models.GithubConn{
				RefreshToken: "refresh_token",
			},
		},
	}

	// 9 minutes remaining (inside 10m buffer)
	tp.conn.TokenExpiresAt = time.Now().Add(9 * time.Minute)
	assert.True(t, tp.needsRefresh())

	// 11 minutes remaining (outside 10m buffer)
	tp.conn.TokenExpiresAt = time.Now().Add(11 * time.Minute)
	assert.False(t, tp.needsRefresh())
}

func TestPersistenceFailure(t *testing.T) {
	mockRT := new(MockRoundTripper)
	client := &http.Client{Transport: mockRT}
	mockDal := new(mockdal.Dal)

	conn := &models.GithubConnection{
		GithubConn: models.GithubConn{
			RefreshToken: "refresh_token",
			GithubAccessToken: models.GithubAccessToken{
				AccessToken: api.AccessToken{
					Token: "old_token",
				},
			},
			GithubAppKey: models.GithubAppKey{
				AppKey: api.AppKey{
					AppId:     "123",
					SecretKey: "secret",
				},
			},
		},
	}

	logger, _ := logruslog.NewDefaultLogger(logrus.New())
	tp := NewTokenProvider(conn, mockDal, client, logger)

	// Mock response for refresh
	respBody := `{"access_token":"new_token","refresh_token":"new_refresh_token","expires_in":3600,"refresh_token_expires_in":3600}`
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(respBody)),
	}
	mockRT.On("RoundTrip", mock.Anything).Return(resp, nil).Once()

	// Mock DAL failure
	mockDal.On("UpdateColumns", mock.Anything, mock.Anything, mock.AnythingOfType("[]dal.Clause")).Return(errors.Default.New("db error"))
	err := tp.ForceRefresh("old_token")
	assert.NoError(t, err) // Should not return error even if persistence fails

	mockRT.AssertExpectations(t)
	mockDal.AssertExpectations(t)
}
