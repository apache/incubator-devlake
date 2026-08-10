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
	"io"
	"net/http"
	"strings"
	"testing"

	devlakeerrors "github.com/apache/incubator-devlake/core/errors"
	helperapi "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

func TestGetJiraSearchEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		deploymentType models.DeploymentType
		want           string
	}{
		{name: "Cloud uses v3", deploymentType: models.DeploymentCloud, want: jiraSearchEndpointV3},
		{name: "Lowercase cloud uses v3", deploymentType: "cloud", want: jiraSearchEndpointV3},
		{name: "Server uses v2", deploymentType: models.DeploymentServer, want: jiraSearchEndpointV2},
		{name: "Data Center uses v2", deploymentType: "Data Center", want: jiraSearchEndpointV2},
		{name: "Unknown uses v2", deploymentType: "something-else", want: jiraSearchEndpointV2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getJiraSearchEndpoint(tt.deploymentType); got != tt.want {
				t.Fatalf("getJiraSearchEndpoint() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestBuildJiraV3SearchRequestBody(t *testing.T) {
	requestData := &helperapi.RequestData{
		Pager: &helperapi.Pager{Size: 100},
	}
	body := buildJiraV3SearchRequestBody("project = DEV", requestData)

	if body["jql"] != "project = DEV" {
		t.Fatalf("expected jql in body")
	}
	if body["maxResults"] != 100 {
		t.Fatalf("expected maxResults=100, got %v", body["maxResults"])
	}
	if body["expand"] != "changelog" {
		t.Fatalf("expected expand=changelog")
	}
	if body["fields"] != "*all" {
		t.Fatalf("expected fields=*all")
	}
	if _, ok := body["nextPageToken"]; ok {
		t.Fatalf("did not expect nextPageToken when CustomData is nil")
	}

	requestData.CustomData = "token-123"
	body = buildJiraV3SearchRequestBody("project = DEV", requestData)
	if body["nextPageToken"] != "token-123" {
		t.Fatalf("expected nextPageToken=token-123, got %v", body["nextPageToken"])
	}

	requestData.CustomData = 123
	body = buildJiraV3SearchRequestBody("project = DEV", requestData)
	if _, ok := body["nextPageToken"]; ok {
		t.Fatalf("did not expect nextPageToken for non-string CustomData")
	}
}

func TestGetNextPageCustomDataForV3(t *testing.T) {
	response := &http.Response{Body: io.NopCloser(strings.NewReader("{\"nextPageToken\":\"abc\",\"issues\":[]}"))}
	customData, err := getNextPageCustomDataForV3(nil, response)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token, ok := customData.(string); !ok || token != "abc" {
		t.Fatalf("expected token abc, got %v", customData)
	}

	bodyAfter, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		t.Fatalf("failed to read response body after parsing token: %v", readErr)
	}
	if string(bodyAfter) != "{\"nextPageToken\":\"abc\",\"issues\":[]}" {
		t.Fatalf("response body should remain readable after token extraction")
	}

	response = &http.Response{Body: io.NopCloser(strings.NewReader("{\"issues\":[]}"))}
	_, err = getNextPageCustomDataForV3(nil, response)
	if !devlakeerrors.Is(err, helperapi.ErrFinishCollect) {
		t.Fatalf("expected ErrFinishCollect when nextPageToken is missing, got %v", err)
	}

	response = &http.Response{Body: io.NopCloser(strings.NewReader("{"))}
	_, err = getNextPageCustomDataForV3(nil, response)
	if err == nil {
		t.Fatalf("expected unmarshal error for malformed JSON")
	}
}
