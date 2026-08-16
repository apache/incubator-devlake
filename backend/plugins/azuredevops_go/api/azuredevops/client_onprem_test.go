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

package azuredevops

import (
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"testing"
)

func TestGetUserProfile_OnPremises(t *testing.T) {
	conn := &models.AzuredevopsConnection{
		BaseConnection: api.BaseConnection{},
		AzuredevopsConn: models.AzuredevopsConn{
			AzuredevopsAccessToken: models.AzuredevopsAccessToken{
				Token: "test-token",
			},
			Endpoint: "https://tfs.company.local/DefaultCollection/",
		},
	}

	client := NewClient(conn, nil, "https://tfs.company.local/DefaultCollection/")
	profile, err := client.GetUserProfile()
	if err != nil {
		t.Fatalf("GetUserProfile() returned unexpected error for On-Premises: %v", err)
	}
	if profile.DisplayName != "On-Premises User" {
		t.Errorf("GetUserProfile() DisplayName = %q; want %q", profile.DisplayName, "On-Premises User")
	}
}

func TestGetUserAccounts_OnPremises(t *testing.T) {
	conn := &models.AzuredevopsConnection{
		BaseConnection: api.BaseConnection{},
		AzuredevopsConn: models.AzuredevopsConn{
			AzuredevopsAccessToken: models.AzuredevopsAccessToken{
				Token: "test-token",
			},
			Endpoint: "https://tfs.company.local/DefaultCollection/",
		},
	}

	client := NewClient(conn, nil, "https://tfs.company.local/DefaultCollection/")
	accounts, err := client.GetUserAccounts("test-member-id")
	if err != nil {
		t.Fatalf("GetUserAccounts() returned unexpected error for On-Premises: %v", err)
	}
	if len(accounts) != 0 {
		t.Errorf("GetUserAccounts() returned %d accounts; want 0 for On-Premises", len(accounts))
	}
}

func TestGetUserProfile_Cloud(t *testing.T) {
	// For Cloud connections (no Endpoint set), GetUserProfile should NOT return
	// the placeholder profile. It should attempt to call the VSSPS API.
	conn := &models.AzuredevopsConnection{
		BaseConnection: api.BaseConnection{},
		AzuredevopsConn: models.AzuredevopsConn{
			AzuredevopsAccessToken: models.AzuredevopsAccessToken{
				Token: "test-token",
			},
			// Endpoint is empty = Cloud mode
		},
	}

	// Without a mock server, this will fail with a connection error,
	// which is expected - the important thing is it does NOT bypass to placeholder
	client := NewClient(conn, nil, "http://localhost:0")
	_, err := client.GetUserProfile()
	if err == nil {
		t.Error("GetUserProfile() for Cloud should attempt real API call and fail without a server")
	}
}

func TestGetUserAccounts_Cloud(t *testing.T) {
	conn := &models.AzuredevopsConnection{
		BaseConnection: api.BaseConnection{},
		AzuredevopsConn: models.AzuredevopsConn{
			AzuredevopsAccessToken: models.AzuredevopsAccessToken{
				Token: "test-token",
			},
			// Endpoint is empty = Cloud mode
		},
	}

	// Without a mock server, this will fail with a connection error
	client := NewClient(conn, nil, "http://localhost:0")
	_, err := client.GetUserAccounts("test-member-id")
	if err == nil {
		t.Error("GetUserAccounts() for Cloud should attempt real API call and fail without a server")
	}
}
