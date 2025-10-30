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
	"testing"

	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/stretchr/testify/assert"
)

func TestBitbucketConn_SetupAuthentication_ApiToken(t *testing.T) {
	// Test API token authentication
	conn := BitbucketConn{
		BasicAuth: api.BasicAuth{
			Username: "user@example.com",
			Password: "api_token_123",
		},
		UsesApiToken: true,
	}

	req, err := http.NewRequest("GET", "https://api.bitbucket.org/2.0/user", nil)
	assert.NoError(t, err)

	authErr := conn.SetupAuthentication(req)
	assert.Nil(t, authErr)
	assert.NotEmpty(t, req.Header.Get("Authorization"))
	assert.Contains(t, req.Header.Get("Authorization"), "Basic")
}

func TestBitbucketConn_SetupAuthentication_AppPassword(t *testing.T) {
	// Test app password authentication
	conn := BitbucketConn{
		BasicAuth: api.BasicAuth{
			Username: "bitbucket_username",
			Password: "app_password_123",
		},
		UsesApiToken: false,
	}

	req, err := http.NewRequest("GET", "https://api.bitbucket.org/2.0/user", nil)
	assert.NoError(t, err)

	authErr := conn.SetupAuthentication(req)
	assert.Nil(t, authErr)
	assert.NotEmpty(t, req.Header.Get("Authorization"))
	assert.Contains(t, req.Header.Get("Authorization"), "Basic")
}

func TestBitbucketConn_Sanitize(t *testing.T) {
	// Test that Sanitize removes the password
	conn := BitbucketConn{
		BasicAuth: api.BasicAuth{
			Username: "user@example.com",
			Password: "secret_password",
		},
		UsesApiToken: true,
	}

	sanitized := conn.Sanitize()
	assert.Empty(t, sanitized.Password)
	assert.Equal(t, "user@example.com", sanitized.Username)
	assert.True(t, sanitized.UsesApiToken)
}

func TestBitbucketConn_UsesApiToken_Default(t *testing.T) {
	// Test that UsesApiToken can be set correctly
	conn := BitbucketConn{
		BasicAuth: api.BasicAuth{
			Username: "user@example.com",
			Password: "token",
		},
		UsesApiToken: true,
	}

	assert.True(t, conn.UsesApiToken)

	// Test app password mode
	conn2 := BitbucketConn{
		BasicAuth: api.BasicAuth{
			Username: "username",
			Password: "password",
		},
		UsesApiToken: false,
	}

	assert.False(t, conn2.UsesApiToken)
}

func TestBitbucketConnection_Sanitize(t *testing.T) {
	// Test that BitbucketConnection.Sanitize works correctly
	connection := BitbucketConnection{
		BitbucketConn: BitbucketConn{
			BasicAuth: api.BasicAuth{
				Username: "user@example.com",
				Password: "secret_token",
			},
			UsesApiToken: true,
		},
	}

	sanitized := connection.Sanitize()
	assert.Empty(t, sanitized.Password)
	assert.Equal(t, "user@example.com", sanitized.Username)
	assert.True(t, sanitized.UsesApiToken)
}

func TestBitbucketConnection_TableName(t *testing.T) {
	connection := BitbucketConnection{}
	assert.Equal(t, "_tool_bitbucket_connections", connection.TableName())
}

func TestBitbucketConnection_MergeFromRequest_PreservesPassword(t *testing.T) {
	original := &BitbucketConnection{
		BitbucketConn: BitbucketConn{
			BasicAuth: api.BasicAuth{
				Username: "user@example.com",
				Password: "secret_token",
			},
			UsesApiToken: true,
		},
	}

	target := &BitbucketConnection{}
	*target = *original // copy

	// Update without password (empty password should preserve original)
	body := map[string]interface{}{
		"username":     "new_user@example.com",
		"usesApiToken": false,
	}

	err := original.MergeFromRequest(target, body)
	assert.NoError(t, err)
	assert.Equal(t, "new_user@example.com", target.Username)
	assert.Equal(t, "secret_token", target.Password) // Should preserve
	assert.False(t, target.UsesApiToken)
}

func TestBitbucketConnection_MergeFromRequest_UpdatesPassword(t *testing.T) {
	original := &BitbucketConnection{
		BitbucketConn: BitbucketConn{
			BasicAuth: api.BasicAuth{
				Username: "user@example.com",
				Password: "old_token",
			},
			UsesApiToken: true,
		},
	}

	target := &BitbucketConnection{}
	*target = *original

	// Update with new password
	body := map[string]interface{}{
		"username":     "user@example.com",
		"password":     "new_token",
		"usesApiToken": true,
	}

	err := original.MergeFromRequest(target, body)
	assert.NoError(t, err)
	assert.Equal(t, "new_token", target.Password) // Should update
}

func TestBitbucketConnection_MergeFromRequest_TogglesUsesApiToken(t *testing.T) {
	original := &BitbucketConnection{
		BitbucketConn: BitbucketConn{
			BasicAuth: api.BasicAuth{
				Username: "user@example.com",
				Password: "credential",
			},
			UsesApiToken: false, // App password
		},
	}

	target := &BitbucketConnection{}
	*target = *original

	// Toggle to API token
	body := map[string]interface{}{
		"usesApiToken": true,
		"password":     "api_token_123",
	}

	err := original.MergeFromRequest(target, body)
	assert.NoError(t, err)
	assert.True(t, target.UsesApiToken)
	assert.Equal(t, "api_token_123", target.Password)
}

func TestBitbucketConn_SetupAuthentication_BasicAuthFormat(t *testing.T) {
	// Test that BOTH methods produce Basic Auth (not Bearer)
	tests := []struct {
		name         string
		username     string
		password     string
		usesApiToken bool
	}{
		{
			name:         "API Token produces Basic Auth",
			username:     "user@example.com",
			password:     "api_token_123",
			usesApiToken: true,
		},
		{
			name:         "App Password produces Basic Auth",
			username:     "bitbucket_username",
			password:     "app_password_123",
			usesApiToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := BitbucketConn{
				BasicAuth: api.BasicAuth{
					Username: tt.username,
					Password: tt.password,
				},
				UsesApiToken: tt.usesApiToken,
			}

			req, _ := http.NewRequest("GET", "https://api.bitbucket.org/2.0/user", nil)
			err := conn.SetupAuthentication(req)

			assert.Nil(t, err)
			authHeader := req.Header.Get("Authorization")
			assert.NotEmpty(t, authHeader)
			assert.Contains(t, authHeader, "Basic ", "Should use Basic auth, not Bearer")
			assert.NotContains(t, authHeader, "Bearer", "Should NOT use Bearer token")
		})
	}
}

func TestBitbucketConn_EmptyPassword(t *testing.T) {
	conn := BitbucketConn{
		BasicAuth: api.BasicAuth{
			Username: "user@example.com",
			Password: "",
		},
		UsesApiToken: true,
	}

	req, _ := http.NewRequest("GET", "https://api.bitbucket.org/2.0/user", nil)
	err := conn.SetupAuthentication(req)

	// Should still set auth header (though it won't work in practice)
	assert.Nil(t, err)
	assert.NotEmpty(t, req.Header.Get("Authorization"))
}

func TestBitbucketConn_SpecialCharactersInPassword(t *testing.T) {
	conn := BitbucketConn{
		BasicAuth: api.BasicAuth{
			Username: "user@example.com",
			Password: "p@ssw0rd!#$%&*()+=",
		},
		UsesApiToken: true,
	}

	req, _ := http.NewRequest("GET", "https://api.bitbucket.org/2.0/user", nil)
	err := conn.SetupAuthentication(req)

	assert.Nil(t, err)
	assert.NotEmpty(t, req.Header.Get("Authorization"))
}

func TestBitbucketConnection_Sanitize_PreservesUsesApiToken(t *testing.T) {
	connection := BitbucketConnection{
		BitbucketConn: BitbucketConn{
			BasicAuth: api.BasicAuth{
				Username: "user@example.com",
				Password: "secret",
			},
			UsesApiToken: true,
		},
	}

	sanitized := connection.Sanitize()

	assert.Empty(t, sanitized.Password, "Password should be removed")
	assert.Equal(t, "user@example.com", sanitized.Username, "Username should be preserved")
	assert.True(t, sanitized.UsesApiToken, "UsesApiToken flag should be preserved")
}
