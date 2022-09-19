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

package plugin

import (
	"context"
	"encoding/json"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/mocks"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/runner"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"sync"
	"testing"
)

func TestPluginRunner(t *testing.T) {
	pluginHelper := newMockPluginHelper()
	pluginHelper.PrepareTaskData(func(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
		taskCtx.GetLogger().Info("running task in %s", pluginHelper.mock.RootPkgPath())
		connectionId, ok := options["ConnectionId"]
		require.True(t, ok)
		require.Equal(t, 1.0, connectionId)
		return nil, nil
	}).Once()
	pluginHelper.SubTaskMetas(func() []core.SubTaskMeta {
		return []core.SubTaskMeta{
			{
				Name: "subtask1",
				EntryPoint: func(c core.SubTaskContext) errors.Error {
					c.GetLogger().Info("inside subtask1")
					return nil
				},
				Required:         false,
				EnabledByDefault: true,
				Description:      "desc",
				DomainTypes:      []string{"dummy_domain"},
			},
		}
	}).Once()
	response, err := runPlugin(t, "test_plugin", pluginHelper.GetPlugin(), &models.Task{
		Plugin:      "test_plugin",
		Subtasks:    toJSON([]string{"subtask1"}),
		Options:     toJSON(map[string]interface{}{"ConnectionId": 1}),
		Status:      models.TASK_CREATED,
		PipelineId:  1,
		PipelineRow: 2,
		PipelineCol: 1,
	})
	require.NoError(t, err)
	require.Equal(t, models.TaskProgressDetail{
		TotalSubTasks:    1,
		FinishedSubTasks: 1,
		TotalRecords:     0,
		FinishedRecords:  0,
		SubTaskName:      "subtask1",
		SubTaskNumber:    1,
	}, response.progressDetail)
	require.Equal(t, models.TASK_COMPLETED, response.result.Status)
}

func TestPluginRunner_AutoRunRequiredTask(t *testing.T) {
	pluginHelper := newMockPluginHelper()
	pluginHelper.PrepareTaskData(func(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
		taskCtx.GetLogger().Info("running task in %s", pluginHelper.mock.RootPkgPath())
		connectionId, ok := options["ConnectionId"]
		require.True(t, ok)
		require.Equal(t, 1.0, connectionId)
		return nil, nil
	}).Once()
	pluginHelper.SubTaskMetas(func() []core.SubTaskMeta {
		return []core.SubTaskMeta{
			{
				Name: "subtask1",
				EntryPoint: func(c core.SubTaskContext) errors.Error {
					c.GetLogger().Info("inside subtask1")
					return nil
				},
				Required:         true,
				EnabledByDefault: true,
				Description:      "desc",
				DomainTypes:      []string{"dummy_domain"},
			},
		}
	}).Once()
	response, err := runPlugin(t, "test_plugin", pluginHelper.GetPlugin(), &models.Task{
		Plugin:      "test_plugin",
		Subtasks:    toJSON([]string{}),
		Options:     toJSON(map[string]interface{}{"ConnectionId": 1}),
		Status:      models.TASK_CREATED,
		PipelineId:  1,
		PipelineRow: 2,
		PipelineCol: 1,
	})
	require.NoError(t, err)
	require.Equal(t, models.TaskProgressDetail{
		TotalSubTasks:    1,
		FinishedSubTasks: 1,
		TotalRecords:     0,
		FinishedRecords:  0,
		SubTaskName:      "subtask1",
		SubTaskNumber:    1,
	}, response.progressDetail)
	require.Equal(t, models.TASK_COMPLETED, response.result.Status)
}

func TestPluginRunner_PrepareError(t *testing.T) {
	pluginHelper := newMockPluginHelper()
	pluginHelper.PrepareTaskData(func(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
		taskCtx.GetLogger().Info("running task in %s", pluginHelper.mock.RootPkgPath())
		return nil, errors.Default.New("prepare task error")
	}).Once()
	pluginHelper.SubTaskMetas(func() []core.SubTaskMeta {
		return []core.SubTaskMeta{
			{
				Name: "subtask1",
				EntryPoint: func(c core.SubTaskContext) errors.Error {
					c.GetLogger().Info("inside subtask1")
					return nil
				},
				Required:         false,
				EnabledByDefault: true,
				Description:      "desc",
				DomainTypes:      []string{"dummy_domain"},
			},
		}
	}).Once()
	response, err := runPlugin(t, "test_plugin", pluginHelper.GetPlugin(), &models.Task{
		Plugin:      "test_plugin",
		Subtasks:    toJSON([]string{"subtask1"}),
		Options:     toJSON(map[string]interface{}{"ConnectionId": 1}),
		Status:      models.TASK_CREATED,
		PipelineId:  1,
		PipelineRow: 2,
		PipelineCol: 1,
	})
	require.Error(t, err)
	require.Contains(t, err.Messages().Format(), "prepare task error")
	require.Equal(t, models.TaskProgressDetail{
		TotalSubTasks:    0,
		FinishedSubTasks: 0,
		TotalRecords:     0,
		FinishedRecords:  0,
		SubTaskName:      "",
		SubTaskNumber:    0,
	}, response.progressDetail)
	require.Equal(t, models.TASK_FAILED, response.result.Status)
}

func TestPluginRunner_EntrypointError(t *testing.T) {
	pluginHelper := newMockPluginHelper()
	pluginHelper.PrepareTaskData(func(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
		return nil, nil
	}).Once()
	pluginHelper.SubTaskMetas(func() []core.SubTaskMeta {
		return []core.SubTaskMeta{
			{
				Name: "subtask1",
				EntryPoint: func(c core.SubTaskContext) errors.Error {
					c.GetLogger().Info("inside subtask1")
					return errors.Default.New("entrypoint error")
				},
				Required:         false,
				EnabledByDefault: true,
				Description:      "desc",
				DomainTypes:      []string{"dummy_domain"},
			},
		}
	}).Once()
	response, err := runPlugin(t, "test_plugin", pluginHelper.GetPlugin(), &models.Task{
		Plugin:      "test_plugin",
		Subtasks:    toJSON([]string{"subtask1"}),
		Options:     toJSON(map[string]interface{}{"ConnectionId": 1}),
		Status:      models.TASK_CREATED,
		PipelineId:  1,
		PipelineRow: 2,
		PipelineCol: 1,
	})
	require.Error(t, err)
	require.Contains(t, err.Messages().Format(), "entrypoint error")
	require.Equal(t, models.TaskProgressDetail{
		TotalSubTasks:    1,
		FinishedSubTasks: 0,
		TotalRecords:     0,
		FinishedRecords:  0,
		SubTaskName:      "subtask1",
		SubTaskNumber:    1,
	}, response.progressDetail)
	require.Equal(t, models.TASK_FAILED, response.result.Status)
}

func TestPluginRunner_WithPrepareData(t *testing.T) {
	pluginHelper := newMockPluginHelper()
	pluginHelper.PrepareTaskData(func(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
		return "test data", nil
	}).Once()
	pluginHelper.SubTaskMetas(func() []core.SubTaskMeta {
		return []core.SubTaskMeta{
			{
				Name: "subtask1",
				EntryPoint: func(c core.SubTaskContext) errors.Error {
					c.GetLogger().Info("inside subtask1")
					require.Equal(t, "test data", c.GetData())
					return nil
				},
				Required:         false,
				EnabledByDefault: true,
				Description:      "desc",
				DomainTypes:      []string{"dummy_domain"},
			},
		}
	}).Once()
	response, err := runPlugin(t, "test_plugin", pluginHelper.GetPlugin(), &models.Task{
		Plugin:      "test_plugin",
		Subtasks:    toJSON([]string{"subtask1"}),
		Options:     toJSON(map[string]interface{}{"ConnectionId": 1}),
		Status:      models.TASK_CREATED,
		PipelineId:  1,
		PipelineRow: 2,
		PipelineCol: 1,
	})
	require.NoError(t, err)
	require.Equal(t, models.TaskProgressDetail{
		TotalSubTasks:    1,
		FinishedSubTasks: 1,
		TotalRecords:     0,
		FinishedRecords:  0,
		SubTaskName:      "subtask1",
		SubTaskNumber:    1,
	}, response.progressDetail)
	require.Equal(t, models.TASK_COMPLETED, response.result.Status)
}

func TestPluginRunner_twoSubtasks(t *testing.T) {
	pluginHelper := newMockPluginHelper()
	pluginHelper.PrepareTaskData(func(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
		taskCtx.GetLogger().Info("running task in %s", pluginHelper.mock.RootPkgPath())
		connectionId, ok := options["ConnectionId"]
		require.True(t, ok)
		require.Equal(t, 1.0, connectionId)
		return nil, nil
	}).Once()
	var subtaskResponses []string
	pluginHelper.SubTaskMetas(func() []core.SubTaskMeta {
		return []core.SubTaskMeta{
			{
				Name: "subtask1",
				EntryPoint: func(c core.SubTaskContext) errors.Error {
					c.GetLogger().Info("inside subtask1")
					subtaskResponses = append(subtaskResponses, "subtask1")
					return nil
				},
				Required:         false,
				EnabledByDefault: true,
				Description:      "desc",
				DomainTypes:      []string{"dummy_domain"},
			},
			{
				Name: "subtask2",
				EntryPoint: func(c core.SubTaskContext) errors.Error {
					c.GetLogger().Info("inside subtask2")
					subtaskResponses = append(subtaskResponses, "subtask2")
					return nil
				},
				Required:         false,
				EnabledByDefault: true,
				Description:      "desc",
				DomainTypes:      []string{"dummy_domain"},
			},
		}
	}).Once()
	response, err := runPlugin(t, "test_plugin", pluginHelper.GetPlugin(), &models.Task{
		Plugin:      "test_plugin",
		Subtasks:    toJSON([]string{"subtask1", "subtask2"}),
		Options:     toJSON(map[string]interface{}{"ConnectionId": 1}),
		Status:      models.TASK_CREATED,
		PipelineId:  1,
		PipelineRow: 2,
		PipelineCol: 1,
	})
	require.NoError(t, err)
	require.Equal(t, models.TaskProgressDetail{
		TotalSubTasks:    2,
		FinishedSubTasks: 2,
		TotalRecords:     0,
		FinishedRecords:  0,
		SubTaskName:      "subtask2",
		SubTaskNumber:    2,
	}, response.progressDetail)
	require.Equal(t, models.TASK_COMPLETED, response.result.Status)
	require.Equal(t, []string{"subtask1", "subtask2"}, subtaskResponses)
}

func TestPluginRunner_twoSubtasks_secondErrorsOut(t *testing.T) {
	pluginHelper := newMockPluginHelper()
	pluginHelper.PrepareTaskData(func(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
		taskCtx.GetLogger().Info("running task in %s", pluginHelper.mock.RootPkgPath())
		connectionId, ok := options["ConnectionId"]
		require.True(t, ok)
		require.Equal(t, 1.0, connectionId)
		return nil, nil
	}).Once()
	var subtaskResponses []string
	pluginHelper.SubTaskMetas(func() []core.SubTaskMeta {
		return []core.SubTaskMeta{
			{
				Name: "subtask1",
				EntryPoint: func(c core.SubTaskContext) errors.Error {
					c.GetLogger().Info("inside subtask1")
					subtaskResponses = append(subtaskResponses, "subtask1")
					return nil
				},
				Required:         false,
				EnabledByDefault: true,
				Description:      "desc",
				DomainTypes:      []string{"dummy_domain"},
			},
			{
				Name: "subtask2",
				EntryPoint: func(c core.SubTaskContext) errors.Error {
					c.GetLogger().Info("inside subtask2")
					subtaskResponses = append(subtaskResponses, "subtask2")
					return errors.Default.New("subtask2 error")
				},
				Required:         false,
				EnabledByDefault: true,
				Description:      "desc",
				DomainTypes:      []string{"dummy_domain"},
			},
		}
	}).Once()
	response, err := runPlugin(t, "test_plugin", pluginHelper.GetPlugin(), &models.Task{
		Plugin:      "test_plugin",
		Subtasks:    toJSON([]string{"subtask1", "subtask2"}),
		Options:     toJSON(map[string]interface{}{"ConnectionId": 1}),
		Status:      models.TASK_CREATED,
		PipelineId:  1,
		PipelineRow: 2,
		PipelineCol: 1,
	})
	require.Error(t, err)
	require.Contains(t, err.Messages().Format(), "subtask2 error")
	require.Equal(t, models.TaskProgressDetail{
		TotalSubTasks:    2,
		FinishedSubTasks: 1,
		TotalRecords:     0,
		FinishedRecords:  0,
		SubTaskName:      "subtask2",
		SubTaskNumber:    2,
	}, response.progressDetail)
	require.Equal(t, models.TASK_FAILED, response.result.Status)
	require.Equal(t, []string{"subtask1", "subtask2"}, subtaskResponses)
}

func TestPluginRunner_twoSubtasks_onlyRunOne(t *testing.T) {
	pluginHelper := newMockPluginHelper()
	pluginHelper.PrepareTaskData(func(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
		taskCtx.GetLogger().Info("running task in %s", pluginHelper.mock.RootPkgPath())
		connectionId, ok := options["ConnectionId"]
		require.True(t, ok)
		require.Equal(t, 1.0, connectionId)
		return nil, nil
	}).Once()
	var subtaskResponses []string
	pluginHelper.SubTaskMetas(func() []core.SubTaskMeta {
		return []core.SubTaskMeta{
			{
				Name: "subtask1",
				EntryPoint: func(c core.SubTaskContext) errors.Error {
					c.GetLogger().Info("inside subtask1")
					subtaskResponses = append(subtaskResponses, "subtask1")
					return nil
				},
				Required:         false,
				EnabledByDefault: true,
				Description:      "desc",
				DomainTypes:      []string{"dummy_domain"},
			},
			{
				Name: "subtask2",
				EntryPoint: func(c core.SubTaskContext) errors.Error {
					c.GetLogger().Info("inside subtask2")
					subtaskResponses = append(subtaskResponses, "subtask2")
					return nil
				},
				Required:         false,
				EnabledByDefault: true,
				Description:      "desc",
				DomainTypes:      []string{"dummy_domain"},
			},
		}
	}).Once()
	response, err := runPlugin(t, "test_plugin", pluginHelper.GetPlugin(), &models.Task{
		Plugin:      "test_plugin",
		Subtasks:    toJSON([]string{"subtask1"}),
		Options:     toJSON(map[string]interface{}{"ConnectionId": 1}),
		Status:      models.TASK_CREATED,
		PipelineId:  1,
		PipelineRow: 2,
		PipelineCol: 1,
	})
	require.NoError(t, err)
	require.Equal(t, models.TaskProgressDetail{
		TotalSubTasks:    1,
		FinishedSubTasks: 1,
		TotalRecords:     0,
		FinishedRecords:  0,
		SubTaskName:      "subtask1",
		SubTaskNumber:    1,
	}, response.progressDetail)
	require.Equal(t, models.TASK_COMPLETED, response.result.Status)
	require.Equal(t, []string{"subtask1"}, subtaskResponses)
}

func toJSON(obj interface{}) datatypes.JSON {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return b
}

func runPlugin(t *testing.T, pluginName string, plugin *mocks.TestPlugin, task *models.Task) (*pluginResponse, errors.Error) {
	ctx := context.Background()
	tester := e2ehelper.NewDataFlowTester(t, pluginName, plugin)
	log := tester.Log.Nested("test")
	err := errors.Convert(tester.Db.Save(task).Error)
	require.NoError(t, err)
	progressDetail := &models.TaskProgressDetail{}
	progChan := make(chan core.RunningProgress)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for p := range progChan {
			runner.UpdateProgressDetail(tester.Db, log, task.ID, progressDetail, &p)
		}
		wg.Done()
	}()
	err = runner.RunTask(ctx, tester.Cfg, log, tester.Db, progChan, task.ID)
	close(progChan)
	wg.Wait()
	plugin.AssertExpectations(t)
	require.NoError(t, tester.Db.Find(task).Error)
	if err != nil {
		return &pluginResponse{tester: tester, result: *task, progressDetail: *progressDetail}, err
	}
	return &pluginResponse{tester: tester, result: *task, progressDetail: *progressDetail}, nil
}

type pluginResponse struct {
	tester         *e2ehelper.DataFlowTester
	result         models.Task
	progressDetail models.TaskProgressDetail
}
