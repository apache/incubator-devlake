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
)

func TestGetEndpoint_CloudDefault(t *testing.T) {
	conn := &AzuredevopsConn{}
	endpoint := conn.GetEndpoint()
	if endpoint != "https://dev.azure.com/" {
		t.Errorf("GetEndpoint() = %q; want %q", endpoint, "https://dev.azure.com/")
	}
}

func TestGetEndpoint_OnPremises(t *testing.T) {
	conn := &AzuredevopsConn{
		Endpoint: "https://tfs.company.local/DefaultCollection/",
	}
	endpoint := conn.GetEndpoint()
	if endpoint != "https://tfs.company.local/DefaultCollection/" {
		t.Errorf("GetEndpoint() = %q; want %q", endpoint, "https://tfs.company.local/DefaultCollection/")
	}
}

func TestGetEndpoint_TrailingSlash(t *testing.T) {
	conn := &AzuredevopsConn{
		Endpoint: "https://tfs.company.local/DefaultCollection",
	}
	endpoint := conn.GetEndpoint()
	expected := "https://tfs.company.local/DefaultCollection/"
	if endpoint != expected {
		t.Errorf("GetEndpoint() = %q; want %q (trailing slash should be appended)", endpoint, expected)
	}
}

func TestGetEndpoint_ConnectionWrapper(t *testing.T) {
	conn := AzuredevopsConnection{
		AzuredevopsConn: AzuredevopsConn{
			Endpoint: "https://tfs.company.local/DefaultCollection/",
		},
	}
	endpoint := conn.GetEndpoint()
	if endpoint != "https://tfs.company.local/DefaultCollection/" {
		t.Errorf("AzuredevopsConnection.GetEndpoint() = %q; want %q", endpoint, "https://tfs.company.local/DefaultCollection/")
	}
}

func TestGetEndpoint_ConnectionWrapperCloud(t *testing.T) {
	conn := AzuredevopsConnection{}
	endpoint := conn.GetEndpoint()
	if endpoint != "https://dev.azure.com/" {
		t.Errorf("AzuredevopsConnection.GetEndpoint() = %q; want %q", endpoint, "https://dev.azure.com/")
	}
}

func TestSetupAuthentication_Cloud(t *testing.T) {
	conn := &AzuredevopsConn{
		AzuredevopsAccessToken: AzuredevopsAccessToken{
			Token: "test-pat-token",
		},
		Username: "testuser",
	}

	req, _ := http.NewRequest("GET", "https://dev.azure.com/org/_apis/projects", nil)
	err := conn.SetupAuthentication(req)
	if err != nil {
		t.Fatalf("SetupAuthentication() returned error: %v", err)
	}

	user, pass, ok := req.BasicAuth()
	if !ok {
		t.Fatal("SetupAuthentication() did not set Basic Auth header")
	}
	if user != "testuser" {
		t.Errorf("Basic Auth username = %q; want %q", user, "testuser")
	}
	if pass != "test-pat-token" {
		t.Errorf("Basic Auth password = %q; want %q", pass, "test-pat-token")
	}
}

func TestSetupAuthentication_OnPremises(t *testing.T) {
	conn := &AzuredevopsConn{
		AzuredevopsAccessToken: AzuredevopsAccessToken{
			Token: "on-prem-pat-token",
		},
		Username: "domain\\admin",
		Endpoint: "https://tfs.company.local/DefaultCollection/",
	}

	req, _ := http.NewRequest("GET", "https://tfs.company.local/DefaultCollection/_apis/projects", nil)
	err := conn.SetupAuthentication(req)
	if err != nil {
		t.Fatalf("SetupAuthentication() returned error: %v", err)
	}

	user, pass, ok := req.BasicAuth()
	if !ok {
		t.Fatal("SetupAuthentication() did not set Basic Auth header")
	}
	// On-Premises: username should be empty (PAT auth with empty username)
	if user != "" {
		t.Errorf("On-Premises Basic Auth username = %q; want empty string", user)
	}
	if pass != "on-prem-pat-token" {
		t.Errorf("On-Premises Basic Auth password = %q; want %q", pass, "on-prem-pat-token")
	}
}

func TestSanitize(t *testing.T) {
	conn := &AzuredevopsConn{
		AzuredevopsAccessToken: AzuredevopsAccessToken{
			Token: "secret-token",
		},
		Username: "testuser",
		Endpoint: "https://tfs.company.local/",
	}

	sanitized := conn.Sanitize()
	if sanitized.Endpoint != "https://tfs.company.local/" {
		t.Errorf("Sanitize() should preserve Endpoint, got %q", sanitized.Endpoint)
	}
	if sanitized.Username != "testuser" {
		t.Errorf("Sanitize() should preserve Username, got %q", sanitized.Username)
	}
}
