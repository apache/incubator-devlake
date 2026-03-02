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

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
)

func TestValidateConnection_Success(t *testing.T) {
	connection := &models.GhCopilotConnection{
		GhCopilotConn: models.GhCopilotConn{
			Organization: "octodemo",
			Token:        "ghp_example",
		},
	}
	connection.Normalize()

	err := validateConnection(connection)
	assert.NoError(t, err)
}

func TestValidateConnection_MissingOrganization(t *testing.T) {
	connection := &models.GhCopilotConnection{
		GhCopilotConn: models.GhCopilotConn{
			Token: "ghp_example",
		},
	}

	err := validateConnection(connection)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "either enterprise or organization is required")
}

func TestValidateConnection_EnterpriseOnly(t *testing.T) {
	connection := &models.GhCopilotConnection{
		GhCopilotConn: models.GhCopilotConn{
			Enterprise: "my-enterprise",
			Token:      "ghp_example",
		},
	}
	connection.Normalize()

	err := validateConnection(connection)
	assert.NoError(t, err)
}

func TestValidateConnection_MissingToken(t *testing.T) {
	connection := &models.GhCopilotConnection{
		GhCopilotConn: models.GhCopilotConn{
			Organization: "octodemo",
			Token:        "",
		},
	}

	err := validateConnection(connection)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is required")
}

func TestValidateConnection_InvalidRateLimit(t *testing.T) {
	connection := &models.GhCopilotConnection{
		GhCopilotConn: models.GhCopilotConn{
			Organization:     "octodemo",
			Token:            "ghp_example",
			RateLimitPerHour: -1,
		},
	}

	err := validateConnection(connection)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rateLimitPerHour must be non-negative")
}
