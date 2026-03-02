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
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/impls/logruslog"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRoundTripper401Refresh(t *testing.T) {
	mockRT := new(MockRoundTripper)
	client := &http.Client{Transport: mockRT}

	expiry := time.Now().Add(10 * time.Minute) // Not expired
	conn := &models.GithubConnection{
		GithubConn: models.GithubConn{
			RefreshToken: "refresh_token",
			GithubAccessToken: models.GithubAccessToken{
				AccessToken: api.AccessToken{
					Token: "old_token",
				},
			},
			TokenExpiresAt: &expiry,
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
	rt := NewRefreshRoundTripper(mockRT, tp)

	// Request
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)

	// 1. First call returns 401
	resp401 := &http.Response{
		StatusCode: 401,
		Body:       io.NopCloser(bytes.NewBufferString("Unauthorized")),
	}
	mockRT.On("RoundTrip", mock.MatchedBy(func(r *http.Request) bool {
		return r.Header.Get("Authorization") == "Bearer old_token" && r.URL.String() == "https://api.github.com/user"
	})).Return(resp401, nil).Once()

	// 2. Refresh call (triggered by 401)
	respRefresh := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(`{"access_token":"new_token","refresh_token":"new_refresh_token","expires_in":3600,"refresh_token_expires_in":3600}`)),
	}
	// The refresh call uses the same client, so it goes through mockRT too!
	mockRT.On("RoundTrip", mock.MatchedBy(func(r *http.Request) bool {
		return r.URL.String() == "https://github.com/login/oauth/access_token"
	})).Return(respRefresh, nil).Once()

	// 3. Retry call with new token
	resp200 := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("Success")),
	}
	mockRT.On("RoundTrip", mock.MatchedBy(func(r *http.Request) bool {
		return r.Header.Get("Authorization") == "Bearer new_token" && r.URL.String() == "https://api.github.com/user"
	})).Return(resp200, nil).Once()

	// Execute
	resp, err := rt.RoundTrip(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Success", string(body))

	mockRT.AssertExpectations(t)
}

func TestRoundTripper401WithAppKeyRefresh(t *testing.T) {
	mockRT := new(MockRoundTripper)

	expiry := time.Now().Add(10 * time.Minute) // Not expired (proactive refresh won't trigger)
	conn := &models.GithubConnection{
		GithubConn: models.GithubConn{
			GithubAccessToken: models.GithubAccessToken{
				AccessToken: api.AccessToken{
					Token: "old_app_token",
				},
			},
			TokenExpiresAt: &expiry,
		},
	}
	// Use tokens slice so GetToken returns the current token
	conn.UpdateToken("old_app_token", "", &expiry, nil)

	// refreshFn simulates minting a new installation token
	refreshCalled := 0
	tp := &TokenProvider{
		conn: conn,
		refreshFn: func(tp *TokenProvider) errors.Error {
			refreshCalled++
			newExpiry := time.Now().Add(1 * time.Hour)
			tp.conn.UpdateToken("new_app_token", "", &newExpiry, nil)
			return nil
		},
	}

	rt := NewRefreshRoundTripper(mockRT, tp)

	req, _ := http.NewRequest("GET", "https://api.github.com/repos/test/test", nil)

	// 1. First call returns 401
	resp401 := &http.Response{
		StatusCode: 401,
		Body:       io.NopCloser(bytes.NewBufferString("Bad credentials")),
	}
	mockRT.On("RoundTrip", mock.MatchedBy(func(r *http.Request) bool {
		return r.Header.Get("Authorization") == "Bearer old_app_token"
	})).Return(resp401, nil).Once()

	// 2. Retry call with new token (after refreshFn runs)
	resp200 := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(`{"full_name":"test/test"}`)),
	}
	mockRT.On("RoundTrip", mock.MatchedBy(func(r *http.Request) bool {
		return r.Header.Get("Authorization") == "Bearer new_app_token"
	})).Return(resp200, nil).Once()

	// Execute
	resp, err := rt.RoundTrip(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, 1, refreshCalled, "refreshFn should have been called exactly once")

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"full_name":"test/test"}`, string(body))

	mockRT.AssertExpectations(t)
}
