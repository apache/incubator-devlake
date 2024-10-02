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

import (
	"fmt"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/teambition/api"
	"github.com/apache/incubator-devlake/plugins/teambition/models"
	"github.com/apache/incubator-devlake/plugins/teambition/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/teambition/tasks"
)

// make sure interface is implemented

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginApi
	plugin.CloseablePluginTask
	plugin.PluginSource
} = (*Teambition)(nil)

type Teambition struct{}

func (p Teambition) Description() string {
	return "collect some Teambition data"
}

func (p Teambition) Name() string {
	return "teambition"
}

func (p Teambition) Init(br context.BasicRes) errors.Error {
	api.Init(br, p)

	return nil
}

func (p Teambition) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.TeambitionConnection{},
		&models.TeambitionAccount{},
		&models.TeambitionTask{},
		&models.TeambitionTaskTagTask{},
		&models.TeambitionTaskTag{},
		&models.TeambitionSprint{},
		&models.TeambitionTaskActivity{},
		&models.TeambitionTaskWorktime{},
		&models.TeambitionProject{},
		&models.TeambitionTaskFlowStatus{},
		&models.TeambitionTaskScenario{},
	}
}

func (p Teambition) Connection() dal.Tabler {
	return &models.TeambitionConnection{}
}

func (p Teambition) Scope() plugin.ToolLayerScope {
	return nil
}

func (p Teambition) ScopeConfig() dal.Tabler {
	return nil
}

func (p Teambition) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectAccountsMeta,
		tasks.ExtractAccountsMeta,
		tasks.ConvertAccountsMeta,
		tasks.CollectTasksMeta,
		tasks.ExtractTasksMeta,
		tasks.CollectTaskTagsMeta,
		tasks.ExtractTaskTagsMeta,
		tasks.ConvertTaskTagTasksMeta,
		tasks.CollectTaskActivitiesMeta,
		tasks.ExtractTaskActivitiesMeta,
		tasks.ConvertTaskCommentsMeta,
		tasks.ConvertTaskChangelogMeta,
		tasks.CollectTaskWorktimeMeta,
		tasks.ExtractTaskWorktimeMeta,
		tasks.ConvertTaskWorktimeMeta,
		tasks.CollectProjectsMeta,
		tasks.ExtractProjectsMeta,
		tasks.ConvertProjectsMeta,
		tasks.CollectSprintsMeta,
		tasks.ExtractSprintsMeta,
		tasks.ConvertSprintsMeta,
		tasks.CollectTaskFlowStatusMeta,
		tasks.ExtractTaskFlowStatusMeta,
		tasks.CollectTaskScenariosMeta,
		tasks.ExtractTaskScenariosMeta,
		tasks.ConvertTasksMeta,
	}
}

func (p Teambition) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope,
	skipCollectors bool,
) (pp coreModels.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes, skipCollectors)
}

func (p Teambition) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
		p.Name(),
	)
	connection := &models.TeambitionConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Teambition connection by the given connection ID")
	}

	var apiClient *helper.ApiAsyncClient
	syncPolicy := taskCtx.SyncPolicy()
	if !syncPolicy.SkipCollectors {
		newApiClient, err := tasks.NewTeambitionApiClient(taskCtx, connection)
		if err != nil {
			return nil, errors.Default.Wrap(err, "unable to get Teambition API client instance")
		}
		apiClient = newApiClient
	}
	taskData := &tasks.TeambitionTaskData{
		Options:   op,
		ApiClient: apiClient,
		TenantId:  connection.TenantId,
	}

	return taskData, nil
}

// RootPkgPath PkgPath information lost when compiled as plugin(.so)
func (p Teambition) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/teambition"
}

func (p Teambition) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Teambition) TestConnection(id uint64) errors.Error {
	_, err := api.TestExistingConnection(helper.GenerateTestingConnectionApiResourceInput(id))
	return err
}

func (p Teambition) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":    api.GetConnection,
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
		},
		"connections/:connectionId/test": {
			"POST": api.TestExistingConnection,
		},
	}
}

func (p Teambition) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.TeambitionTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	if data != nil && data.ApiClient != nil {
		data.ApiClient.Release()
	}
	return nil
}
