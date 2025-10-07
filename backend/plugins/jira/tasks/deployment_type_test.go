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
	"strings"
	"testing"

	"github.com/apache/incubator-devlake/plugins/jira/models"
)

func TestDeploymentTypeComparison(t *testing.T) {
	tests := []struct {
		name           string
		deploymentType models.DeploymentType
		wantServer     bool
	}{
		{
			name:           "lowercase server",
			deploymentType: models.DeploymentServer,
			wantServer:     true,
		},
		{
			name:           "uppercase Server",
			deploymentType: "Server",
			wantServer:     true,
		},
		{
			name:           "mixed case SeRvEr",
			deploymentType: "SeRvEr",
			wantServer:     true,
		},
		{
			name:           "lowercase cloud",
			deploymentType: models.DeploymentCloud,
			wantServer:     false,
		},
		{
			name:           "uppercase Cloud",
			deploymentType: "Cloud",
			wantServer:     false,
		},
		{
			name:           "mixed case ClOuD",
			deploymentType: "ClOuD",
			wantServer:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the case-insensitive comparison logic used in the collectors
			// We compare lowercase versions because the API might return "Cloud" or "Server" with capital letters
			isServer := strings.ToLower(string(tt.deploymentType)) == strings.ToLower(string(models.DeploymentServer))
			if isServer != tt.wantServer {
				t.Errorf("DeploymentType comparison for %q: got isServer=%v, want %v", tt.deploymentType, isServer, tt.wantServer)
			}
		})
	}
}

func TestDeploymentTypeConstants(t *testing.T) {
	// Verify the constants have the expected values
	if models.DeploymentServer != "Server" {
		t.Errorf("DeploymentServer constant = %q, want %q", models.DeploymentServer, "Server")
	}
	if models.DeploymentCloud != "Cloud" {
		t.Errorf("DeploymentCloud constant = %q, want %q", models.DeploymentCloud, "Cloud")
	}
}
