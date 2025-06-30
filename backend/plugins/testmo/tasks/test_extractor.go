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
	"regexp"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/testmo/models"
)

var ExtractTestsMeta = plugin.SubTaskMeta{
	Name:             "extractTests",
	EntryPoint:       ExtractTests,
	EnabledByDefault: true,
	Description:      "Extract raw tests data into tool layer table testmo_tests",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

func ExtractTests(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*TestmoTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:   taskCtx,
			Table: RAW_TEST_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			var apiTest struct {
				Id              uint64     `json:"id"`
				ProjectId       uint64     `json:"project_id"`
				AutomationRunId uint64     `json:"automation_run_id"`
				ThreadId        uint64     `json:"thread_id"`
				Name            string     `json:"name"`
				Key             string     `json:"key"`
				Status          int32      `json:"status"`
				StatusName      string     `json:"status_name"`
				Elapsed         *int64     `json:"elapsed"`
				Message         string     `json:"message"`
				CreatedAt       *time.Time `json:"created_at"`
				UpdatedAt       *time.Time `json:"updated_at"`
			}

			err := json.Unmarshal(row.Data, &apiTest)
			if err != nil {
				return nil, errors.Default.Wrap(err, "error unmarshaling test")
			}

			test := &models.TestmoTest{
				ConnectionId:    data.Options.ConnectionId,
				Id:              apiTest.Id,
				ProjectId:       apiTest.ProjectId,
				AutomationRunId: apiTest.AutomationRunId,
				ThreadId:        apiTest.ThreadId,
				Name:            apiTest.Name,
				Key:             apiTest.Key,
				Status:          apiTest.Status,
				StatusName:      apiTest.StatusName,
				Elapsed:         apiTest.Elapsed,
				Message:         apiTest.Message,
				TestmoCreatedAt: apiTest.CreatedAt,
				TestmoUpdatedAt: apiTest.UpdatedAt,
			}

			// Classify test types based on scope config patterns
			if data.Options.ScopeConfig != nil {
				test.IsAcceptanceTest = matchesPattern(test.Name, data.Options.ScopeConfig.AcceptanceTestPattern)
				test.IsSmokeTest = matchesPattern(test.Name, data.Options.ScopeConfig.SmokeTestPattern)
				test.Team = extractTeam(test.Name, data.Options.ScopeConfig.TeamPattern)
			}

			return []interface{}{test}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func matchesPattern(testName, pattern string) bool {
	if pattern == "" {
		return false
	}

	// Simple case-insensitive substring match for now
	// Can be enhanced to support regex patterns
	return strings.Contains(strings.ToLower(testName), strings.ToLower(pattern))
}

func extractTeam(testName, pattern string) string {
	if pattern == "" {
		return ""
	}

	// Try to extract team using regex pattern
	if regex, err := regexp.Compile(pattern); err == nil {
		matches := regex.FindStringSubmatch(testName)
		if len(matches) > 1 {
			return matches[1] // Return first capture group
		}
	}

	return ""
}
