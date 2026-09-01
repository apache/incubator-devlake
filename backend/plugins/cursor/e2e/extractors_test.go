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
	"github.com/apache/incubator-devlake/plugins/cursor/impl"
	"github.com/apache/incubator-devlake/plugins/cursor/models"
	"github.com/apache/incubator-devlake/plugins/cursor/tasks"
)

func TestCursorExtractorsDataFlow(t *testing.T) {
	cfg := config.GetConfig()
	dbUrl := cfg.GetString("E2E_DB_URL")
	if dbUrl == "" {
		t.Skip("skipping e2e test: E2E_DB_URL is not set")
	}
	if err := runner.CheckDbConnection(dbUrl, 10*time.Second); err != nil {
		t.Skipf("skipping e2e test: cannot connect to E2E_DB_URL: %v", err)
	}

	var cursorPlugin impl.Cursor
	dataflowTester := e2ehelper.NewDataFlowTester(t, "cursor", cursorPlugin)

	taskData := &tasks.CursorTaskData{
		Options: &tasks.CursorOptions{
			ConnectionId: 1,
			ScopeId:      "team",
		},
		Connection: &models.CursorConnection{
			CursorConn: models.CursorConn{
				RestConnection: helper.RestConnection{Endpoint: models.DefaultEndpoint},
			},
		},
	}

	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_cursor_usage_events.csv", "_raw_cursor_usage_events")
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_cursor_members.csv", "_raw_cursor_members")
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_cursor_user_spend.csv", "_raw_cursor_user_spend")
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_cursor_daily_usage.csv", "_raw_cursor_daily_usage")

	dataflowTester.FlushTabler(&models.CursorUsageEvent{})
	dataflowTester.FlushTabler(&models.CursorMember{})
	dataflowTester.FlushTabler(&models.CursorUserSpend{})
	dataflowTester.FlushTabler(&models.CursorDailyUsage{})

	dataflowTester.Subtask(tasks.ExtractMembersMeta, taskData)
	dataflowTester.Subtask(tasks.ExtractUsageEventsMeta, taskData)
	dataflowTester.Subtask(tasks.ExtractUserSpendMeta, taskData)
	dataflowTester.Subtask(tasks.ExtractDailyUsageMeta, taskData)

	dataflowTester.VerifyTableWithOptions(&models.CursorUsageEvent{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_cursor_usage_events.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.VerifyTableWithOptions(&models.CursorMember{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_cursor_members.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.VerifyTableWithOptions(&models.CursorUserSpend{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_cursor_user_spend.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.VerifyTableWithOptions(&models.CursorDailyUsage{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_cursor_daily_usage.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// Non-degradation: appending newer raw rows and re-extracting must never
	// shrink the tool table or move MAX(usage_date) backward.
	var beforeCount int64
	if err := dataflowTester.Db.Model(&models.CursorDailyUsage{}).Count(&beforeCount).Error; err != nil {
		t.Fatalf("count before: %v", err)
	}
	var beforeMax time.Time
	if err := dataflowTester.Db.Model(&models.CursorDailyUsage{}).Select("MAX(usage_date)").Scan(&beforeMax).Error; err != nil {
		t.Fatalf("max before: %v", err)
	}

	extraRaw := map[string]interface{}{
		"id":         uint64(2),
		"params":     `{"ConnectionId":1,"ScopeId":"team","Endpoint":"https://api.cursor.com"}`,
		"data":       []byte(`{"day":"2026-07-09","userId":"user_example123","email":"user@example.com","isActive":true,"completions":1,"premiumRequests":0,"agentRequests":1,"chatRequests":0,"composerRequests":0,"totalTabsAccepted":0,"totalTabsShown":0,"usageBasedReqs":0,"subscriptionIncludedReqs":1,"mostUsedModel":"composer-2.5-fast","clientVersion":"0.50.3","totalLinesAdded":0,"totalLinesDeleted":0,"acceptedLinesAdded":0,"acceptedLinesDeleted":0,"totalApplies":0,"totalAccepts":0,"totalRejects":0}`),
		"url":        "https://api.cursor.com/teams/daily-usage-data",
		"input":      nil,
		"created_at": time.Now().UTC(),
	}
	if err := dataflowTester.Db.Table("_raw_cursor_daily_usage").Create(extraRaw).Error; err != nil {
		t.Fatalf("insert extra raw: %v", err)
	}

	dataflowTester.Subtask(tasks.ExtractDailyUsageMeta, taskData)

	var afterCount int64
	if err := dataflowTester.Db.Model(&models.CursorDailyUsage{}).Count(&afterCount).Error; err != nil {
		t.Fatalf("count after: %v", err)
	}
	var afterMax time.Time
	if err := dataflowTester.Db.Model(&models.CursorDailyUsage{}).Select("MAX(usage_date)").Scan(&afterMax).Error; err != nil {
		t.Fatalf("max after: %v", err)
	}
	if afterCount < beforeCount {
		t.Fatalf("tool row count decreased: before=%d after=%d", beforeCount, afterCount)
	}
	if afterMax.Before(beforeMax) {
		t.Fatalf("MAX(usage_date) moved backward: before=%v after=%v", beforeMax, afterMax)
	}
}
