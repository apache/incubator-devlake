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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apache/incubator-devlake/core/config"
	coremodels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/helpers/apikeyhelper"
	contextimpl "github.com/apache/incubator-devlake/impls/context"
	"github.com/apache/incubator-devlake/impls/logruslog"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// requireUserGate simulates RequireAuth with AUTH_ENABLED=true: it rejects any
// request whose gin context does not carry an authenticated user. This is the
// exact check that caused 401s for valid REST API key requests.
func requireUserGate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := shared.GetUser(c); !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "unauthorized",
			})
			return
		}
		c.Next()
	}
}

// TestRestAuthKeyReachesHandlerWhenAuthEnabled sends a valid /rest/... Bearer
// token request through a router that also has a RequireAuth-style gate. The
// request must reach the downstream handler (200) rather than being rejected
// by the auth gate (401).
//
// Without the fix, gin's HandleContext resets c.Keys, wiping the user that
// RestAuthentication stored before rerouting, so the auth gate sees no user
// and returns 401.
func TestRestAuthKeyReachesHandlerWhenAuthEnabled(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// apikeyhelper reads ENCRYPTION_SECRET from the global viper config.
	const encryptionSecret = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" // 32 bytes
	config.GetConfig().Set("ENCRYPTION_SECRET", encryptionSecret)

	basicRes := contextimpl.NewDefaultBasicRes(config.GetConfig(), logruslog.Global, nil)
	helper := apikeyhelper.NewApiKeyHelper(basicRes, logruslog.Global)

	const plaintext = "test-api-key-plaintext"
	hashedKey, hashErr := helper.DigestToken(plaintext)
	if hashErr != nil {
		t.Fatalf("DigestToken: %v", hashErr)
	}

	// Mock DAL: when First is called, populate the destination with a valid key
	// whose AllowedPath covers the webhook endpoint under test.
	db := &mockdal.Dal{}
	db.On("First", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			dst := args.Get(0).(*coremodels.ApiKey)
			dst.ApiKey = hashedKey
			dst.AllowedPath = `/plugins/webhook/connections/1/.*`
			dst.Creator = common.Creator{Creator: "test-user"}
		}).
		Return(nil)

	basicRes = contextimpl.NewDefaultBasicRes(config.GetConfig(), logruslog.Global, db)

	router := gin.New()
	router.Use(RestAuthentication(router, basicRes))
	router.Use(requireUserGate())
	router.POST("/plugins/webhook/connections/:id/deployments", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/rest/plugins/webhook/connections/1/deployments", nil)
	req.Header.Set("Authorization", "Bearer "+plaintext)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 with valid REST API key when auth gate is active, got %d: %s",
			resp.Code, resp.Body.String())
	}
}
