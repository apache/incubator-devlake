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

package e2e

import (
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/runner"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/impl"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/tasks"
)

func TestCopilotMetricsDataFlow(t *testing.T) {
	cfg := config.GetConfig()
	dbUrl := cfg.GetString("E2E_DB_URL")
	if dbUrl == "" {
		t.Skip("skipping e2e test: E2E_DB_URL is not set")
	}
	if err := runner.CheckDbConnection(dbUrl, 10*time.Second); err != nil {
		t.Skipf("skipping e2e test: cannot connect to E2E_DB_URL: %v", err)
	}

	var copilot impl.GhCopilot
	dataflowTester := e2ehelper.NewDataFlowTester(t, "gh-copilot", copilot)

	taskData := &tasks.GhCopilotTaskData{
		Options: &tasks.GhCopilotOptions{
			ConnectionId: 1,
			ScopeId:      "octodemo",
		},
		Connection: &models.GhCopilotConnection{
			GhCopilotConn: models.GhCopilotConn{
				RestConnection:   helper.RestConnection{Endpoint: "https://api.github.com"},
				Organization:     "octodemo",
				RateLimitPerHour: 5000,
			},
		},
	}

	dataflowTester.ImportCsvIntoRawTable("./metrics/raw_tables/_raw_copilot_metrics.csv", "_raw_copilot_org_metrics")
	dataflowTester.ImportCsvIntoRawTable("./metrics/raw_tables/_raw_copilot_seats.csv", "_raw_copilot_seats")

	dataflowTester.FlushTabler(&models.GhCopilotOrgMetrics{})
	dataflowTester.FlushTabler(&models.GhCopilotLanguageMetrics{})
	dataflowTester.FlushTabler(&models.GhCopilotSeat{})

	dataflowTester.Subtask(tasks.ExtractSeatsMeta, taskData)
	dataflowTester.Subtask(tasks.ExtractOrgMetricsMeta, taskData)

	dataflowTester.VerifyTableWithOptions(&models.GhCopilotOrgMetrics{}, e2ehelper.TableOptions{
		CSVRelPath:  "./metrics/snapshot_tables/_tool_copilot_org_daily_metrics.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.VerifyTableWithOptions(&models.GhCopilotSeat{}, e2ehelper.TableOptions{
		CSVRelPath: "./metrics/snapshot_tables/_tool_copilot_seats.csv",
		IgnoreTypes: []interface{}{
			common.RawDataOrigin{},
		},
	})

	dataflowTester.VerifyTableWithOptions(&models.GhCopilotLanguageMetrics{}, e2ehelper.TableOptions{
		CSVRelPath:  "./metrics/language_breakdown.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
