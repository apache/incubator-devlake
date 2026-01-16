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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
)

type githubErrorPayload struct {
	Message string `json:"message"`
}

func buildGitHubApiError(status int, organization string, body []byte, retryAfter string) errors.Error {
	msg := strings.TrimSpace(string(body))
	if len(body) > 0 {
		payload := &githubErrorPayload{}
		if err := json.Unmarshal(body, payload); err == nil && payload.Message != "" {
			msg = payload.Message
		}
	}

	var prefix string
	switch status {
	case http.StatusForbidden:
		prefix = "GitHub returned 403 Forbidden. Ensure the PAT includes manage_billing:copilot and the organization has Copilot access."
	case http.StatusNotFound:
		prefix = fmt.Sprintf("GitHub returned 404 Not Found for organization '%s'. Verify the organization slug and Copilot availability.", organization)
	case http.StatusUnprocessableEntity:
		prefix = "GitHub returned 422 Unprocessable Entity. Enable Copilot metrics for the organization before running collection."
	case http.StatusTooManyRequests:
		prefix = "GitHub rate limited the request (429). Respect Retry-After guidance before retrying."
	default:
		prefix = fmt.Sprintf("GitHub API request failed with status %d.", status)
	}

	if retryAfter != "" {
		prefix = fmt.Sprintf("%s Retry-After: %s.", prefix, retryAfter)
	}
	if msg != "" {
		prefix = fmt.Sprintf("%s Details: %s", prefix, msg)
	}
	return errors.HttpStatus(status).New(strings.TrimSpace(prefix))
}
