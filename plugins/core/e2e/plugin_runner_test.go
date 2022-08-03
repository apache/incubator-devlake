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
	"context"
	"encoding/json"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/migration"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/runner"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"sync"
	"testing"
)

func TestPluginRunner(t *testing.T) {
	ctx := context.Background()
	plugin := TestPlugin{}
	plugin.AddSubTaskMetas(core.SubTaskMeta{
		Name: "subtask1",
		EntryPoint: func(c core.SubTaskContext) error {
			plugin.subtasksCalled = append(plugin.subtasksCalled, "subtask1")
			c.GetLogger().Info("inside subtask1")
			return nil
		},
		Required:         true,
		EnabledByDefault: true,
		Description:      "desc",
		DomainTypes:      []string{"dummy_domain"},
	})
	tester := e2ehelper.NewDataFlowTester(t, "test_plugin", &plugin)
	log := logger.Global.Nested("test")
	runMigrations(t, ctx, tester)
	task := models.Task{
		Plugin: "test_plugin",
		Subtasks: toJSON([]string{
			"subtask1",
		}),
		Options: toJSON(map[string]interface{}{
			"ConnectionId": 1,
		}),
		Status:      models.TASK_CREATED,
		PipelineId:  1,
		PipelineRow: 2,
		PipelineCol: 1,
	}
	err := tester.Db.Save(&task).Error
	require.NoError(t, err)

	progressDetail := &models.TaskProgressDetail{}
	progChan := make(chan core.RunningProgress)
	wg := &sync.WaitGroup{}
	defer close(progChan)
	go func() {
		for p := range progChan {
			runner.UpdateProgressDetail(tester.Db, log, task.ID, progressDetail, &p)
		}
		wg.Done()
	}()
	err = runner.RunTask(ctx, tester.Cfg, log, tester.Db, progChan, task.ID)
	require.NoError(t, err)
	require.Equal(t, []string{"subtask1"}, plugin.subtasksCalled)
	wg.Wait()
	require.Equal(t, models.TaskProgressDetail{
		TotalSubTasks:    1,
		FinishedSubTasks: 1,
		TotalRecords:     0,
		FinishedRecords:  0,
		SubTaskName:      "subtask1",
		SubTaskNumber:    1,
	}, *progressDetail)
}

func toJSON(obj interface{}) datatypes.JSON {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return b
}

func runMigrations(t *testing.T, ctx context.Context, tester *e2ehelper.DataFlowTester) {
	tester.DropAllTables()
	migration.Init(tester.Db)
	runner.RegisterMigrationScripts(migrationscripts.All(), "Framework", nil, logger.Global)
	err := migration.Execute(ctx)
	require.NoError(t, err)
}
