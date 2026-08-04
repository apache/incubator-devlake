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
	"testing"

	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func TestValidateUserTokenPrefix(t *testing.T) {
	t.Parallel()

	serverEndpoint := "https://rad-sonar.example.com/api/"
	cloudEndpoint := "https://sonarcloud.io/api/"

	tests := []struct {
		name    string
		conn    SonarqubeConn
		wantErr bool
	}{
		{
			name: "user token on server",
			conn: SonarqubeConn{
				RestConnection:       helper.RestConnection{Endpoint: serverEndpoint},
				SonarqubeAccessToken: SonarqubeAccessToken{Token: "squ_abc123"},
			},
		},
		{
			name: "global analysis token on server",
			conn: SonarqubeConn{
				RestConnection:       helper.RestConnection{Endpoint: serverEndpoint},
				SonarqubeAccessToken: SonarqubeAccessToken{Token: "sqa_abc123"},
			},
			wantErr: true,
		},
		{
			name: "project analysis token on server",
			conn: SonarqubeConn{
				RestConnection:       helper.RestConnection{Endpoint: serverEndpoint},
				SonarqubeAccessToken: SonarqubeAccessToken{Token: "sqp_abc123"},
			},
			wantErr: true,
		},
		{
			name: "unknown prefix on server",
			conn: SonarqubeConn{
				RestConnection:       helper.RestConnection{Endpoint: serverEndpoint},
				SonarqubeAccessToken: SonarqubeAccessToken{Token: "legacy-token"},
			},
			wantErr: true,
		},
		{
			name: "sonarcloud skips prefix check",
			conn: SonarqubeConn{
				RestConnection:       helper.RestConnection{Endpoint: cloudEndpoint},
				SonarqubeAccessToken: SonarqubeAccessToken{Token: "sqa_abc123"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.conn.ValidateUserTokenPrefix()
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
