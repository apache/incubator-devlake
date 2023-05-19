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
	"github.com/apache/incubator-devlake/core/runner"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/impls/dalgorm"
	"github.com/apache/incubator-devlake/plugins/zentao/impl"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"github.com/apache/incubator-devlake/plugins/zentao/tasks"
)

func TestZentaoDbGetDataFlow(t *testing.T) {

	var zentao impl.Zentao
	dataflowTester := e2ehelper.NewDataFlowTester(t, "zentao", zentao)
	cfg := config.GetConfig()

	taskData := &tasks.ZentaoTaskData{
		Options: &tasks.ZentaoOptions{
			ConnectionId: 1,
			ProjectId:    0,
			ProductId:    3,

			BaseDbConfigReader: runner.BaseDbConfigReader{
				DbUrl:          "mysql://root:merico@sshd-proxy:3306/zentao?charset=utf8mb4&parseTime=True",
				DbLoggingLevel: cfg.GetString("DB_LOGGING_LEVEL"),
				DbIdleConns:    cfg.GetInt("DB_IDLE_CONNS"),
				DbMaxConns:     cfg.GetInt("DB_MAX_CONNS"),
			},
		},
	}

	rgorm, err := runner.NewGormDb(&taskData.Options.BaseDbConfigReader, dataflowTester.Log)
	if err != nil {
		return
	}
	taskData.Options.RemoteDb = dalgorm.NewDalgorm(rgorm)

	// verify conversion
	dataflowTester.FlushTabler(&models.ZentaoChangelog{})
	dataflowTester.FlushTabler(&models.ZentaoChangelogDetail{})
	dataflowTester.Subtask(tasks.DBGetChangelogMeta, taskData)

}
