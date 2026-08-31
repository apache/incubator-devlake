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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
)

func decodeBody(t *testing.T, raw string) map[string]interface{} {
	t.Helper()
	var body map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &body); err != nil {
		t.Fatalf("failed to unmarshal test payload: %v", err)
	}
	return body
}

func TestValidateConnection_EnterpriseSlugOnly_ViaDecode(t *testing.T) {
	body := decodeBody(t, `{
		"name": "my-copilot-conn",
		"endpoint": "https://api.github.com",
		"organization": "",
		"enterprise": "my-company",
		"token": "ghp_example",
		"rateLimitPerHour": 5000
	}`)

	connection := &models.GhCopilotConnection{}
	assert.NoError(t, helper.Decode(body, connection, vld))

	connection.Normalize()
	assert.Equal(t, "my-company", connection.Enterprise)
	assert.True(t, connection.HasEnterprise())

	err := validateConnection(connection)
	assert.NoError(t, err)
}

func TestValidateConnection_EnterpriseSlugWithWhitespace_ViaDecode(t *testing.T) {
	body := decodeBody(t, `{
		"organization": "   ",
		"enterprise": "  my-company  ",
		"token": "ghp_example"
	}`)

	connection := &models.GhCopilotConnection{}
	assert.NoError(t, helper.Decode(body, connection, vld))

	connection.Normalize()
	assert.Equal(t, "my-company", connection.Enterprise)
	assert.Equal(t, "", connection.Organization)

	err := validateConnection(connection)
	assert.NoError(t, err)
}

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
