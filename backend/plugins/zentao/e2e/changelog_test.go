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

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/runner"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/impls/dalgorm"
	"github.com/apache/incubator-devlake/plugins/zentao/impl"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"github.com/apache/incubator-devlake/plugins/zentao/tasks"
	"github.com/spf13/viper"
)

func TestZentaoDbGetDataFlow(t *testing.T) {

	var zentao impl.Zentao
	dataflowTester := e2ehelper.NewDataFlowTester(t, "zentao", zentao)
	cfg := config.GetConfig()

	taskData := &tasks.ZentaoTaskData{
		Options: &tasks.ZentaoOptions{
			ConnectionId: 1,
			ProjectId:    0,
		},
		Stories:     map[int64]struct{}{},
		Tasks:       map[int64]struct{}{10: {}, 11: {}, 14: {}},
		Bugs:        map[int64]struct{}{1: {}, 2: {}, 3: {}, 4: {}},
		ApiClient:   getFakeAPIClient(),
		HomePageURL: getFakeHomepage(),
	}

	dataflowTester.ImportCsvIntoTabler("./raw_tables/zt_action.csv", models.ZentaoRemoteDbAction{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/zt_history.csv", models.ZentaoRemoteDbHistory{})

	v := viper.New()
	v.Set("DB_URL", cfg.GetString(`E2E_DB_URL`))
	v.Set("DB_LOGGING_LEVEL", cfg.GetString("DB_LOGGING_LEVEL"))
	v.Set("DB_IDLE_CONNS", cfg.GetInt("DB_IDLE_CONNS"))
	v.Set("DbMaxConns", cfg.GetInt("DB_MAX_CONNS"))

	rgorm, err := runner.NewGormDb(v, dataflowTester.Log)
	if err != nil {
		return
	}
	taskData.RemoteDb = dalgorm.NewDalgorm(rgorm)

	// verify conversion
	dataflowTester.FlushTabler(&models.ZentaoChangelog{})
	dataflowTester.FlushTabler(&models.ZentaoChangelogDetail{})
	dataflowTester.Subtask(tasks.DBGetChangelogMeta, taskData)

	dataflowTester.VerifyTable(
		models.ZentaoChangelog{},
		"./snapshot_tables/_tool_zentao_changelog.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"id",
			"object_id",
			"execution",
			"actor",
			"action",
			"extra",
			"object_type",
			"project",
			"product",
			"vision",
			"comment",
			"efforted",
		),
	)

	dataflowTester.VerifyTableWithOptions(
		&models.ZentaoChangelogDetail{},
		e2ehelper.TableOptions{
			CSVRelPath:  "./snapshot_tables/_tool_zentao_changelog_detail.csv",
			IgnoreTypes: []interface{}{common.NoPKModel{}},
		})

	dataflowTester.FlushTabler(&ticket.IssueChangelogs{})
	dataflowTester.Subtask(tasks.ConvertChangelogMeta, taskData)
	dataflowTester.VerifyTableWithOptions(
		&ticket.IssueChangelogs{},
		e2ehelper.TableOptions{
			CSVRelPath:  "./snapshot_tables/issue_changelogs.csv",
			IgnoreTypes: []interface{}{common.NoPKModel{}},
		})
}
