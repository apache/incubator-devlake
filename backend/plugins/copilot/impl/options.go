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

package impl

import "github.com/apache/incubator-devlake/plugins/copilot/models"

const (
	// DefaultEndpoint is the GitHub REST API endpoint used for Copilot metrics.
	DefaultEndpoint = "https://api.github.com"
	// DefaultRateLimitPerHour mirrors GitHub's default rate limit for PATs.
	DefaultRateLimitPerHour = 5000
)

// NormalizeConnection ensures required defaults are set before use.
func NormalizeConnection(connection *models.CopilotConnection) {
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
