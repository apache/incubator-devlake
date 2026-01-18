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
	"net/http"
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const (
	// DefaultEndpoint is the GitHub REST API endpoint used for Copilot metrics.
	DefaultEndpoint = "https://api.github.com"
	// DefaultRateLimitPerHour mirrors GitHub's default rate limit for PATs.
	DefaultRateLimitPerHour = 5000
)

// GhCopilotConn stores GitHub Copilot connection settings.
type GhCopilotConn struct {
	helper.RestConnection `mapstructure:",squash"`

	Token            string `mapstructure:"token" json:"token"`
	Organization     string `mapstructure:"organization" json:"organization"`
	RateLimitPerHour int    `mapstructure:"rateLimitPerHour" json:"rateLimitPerHour"`
}

// SetupAuthentication implements plugin.ApiAuthenticator so helper.NewApiClientFromConnection
// can attach the Authorization header for GitHub API requests.
func (conn *GhCopilotConn) SetupAuthentication(request *http.Request) errors.Error {
	if conn == nil {
		return errors.BadInput.New("connection is required")
	}
	token := strings.TrimSpace(conn.Token)
	if token == "" {
		return errors.BadInput.New("token is required")
	}

	lower := strings.ToLower(token)
	if strings.HasPrefix(lower, "bearer ") || strings.HasPrefix(lower, "token ") {
		request.Header.Set("Authorization", token)
		return nil
	}
	request.Header.Set("Authorization", "Bearer "+token)
	return nil
}

func (conn *GhCopilotConn) Sanitize() GhCopilotConn {
	if conn == nil {
		return GhCopilotConn{}
	}
	clone := *conn
	clone.Token = utils.SanitizeString(clone.Token)
	return clone
}

// GhCopilotConnection persists connection details with metadata required by DevLake.
type GhCopilotConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	GhCopilotConn         `mapstructure:",squash"`
}

func (GhCopilotConnection) TableName() string {
	return "_tool_copilot_connections"
}

func (connection GhCopilotConnection) Sanitize() GhCopilotConnection {
	connection.GhCopilotConn = connection.GhCopilotConn.Sanitize()
	return connection
}

func (connection *GhCopilotConnection) MergeFromRequest(target *GhCopilotConnection, body map[string]interface{}) error {
	if target == nil {
		return nil
	}
	originalToken := target.Token
	if err := helper.DecodeMapStruct(body, target, true); err != nil {
		return err
	}
	sanitizedOriginal := utils.SanitizeString(originalToken)
	if target.Token == "" || target.Token == sanitizedOriginal {
		target.Token = originalToken
	}
	return nil
}

// Normalize applies default connection values where necessary.
func (connection *GhCopilotConnection) Normalize() {
	if connection == nil {
		return
	}
	if connection.Endpoint == "" {
		connection.Endpoint = DefaultEndpoint
	}
	if connection.RateLimitPerHour <= 0 {
		connection.RateLimitPerHour = DefaultRateLimitPerHour
	}
}
