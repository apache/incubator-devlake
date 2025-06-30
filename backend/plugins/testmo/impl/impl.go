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
	pluginhelper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/testmo/api"
	"github.com/apache/incubator-devlake/plugins/testmo/models"
	"github.com/apache/incubator-devlake/plugins/testmo/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/testmo/tasks"
)

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginApi
	plugin.PluginModel
	plugin.PluginMigration
	plugin.CloseablePluginTask
	plugin.PluginSource
	plugin.DataSourcePluginBlueprintV200
} = (*Testmo)(nil)

type Testmo struct{}

func (p Testmo) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)
	return nil
}

func (p Testmo) Connection() dal.Tabler {
	return &models.TestmoConnection{}
}

func (p Testmo) Scope() plugin.ToolLayerScope {
	return &models.TestmoProject{}
}

func (p Testmo) ScopeConfig() dal.Tabler {
	return &models.TestmoScopeConfig{}
}

func (p Testmo) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.TestmoConnection{},
		&models.TestmoProject{},
		&models.TestmoScopeConfig{},
		&models.TestmoAutomationRun{},
		&models.TestmoTest{},
		&models.TestmoMilestone{},
	}
}

func (p Testmo) Description() string {
	return "To collect and enrich data from Testmo"
}

func (p Testmo) Name() string {
	return "testmo"
}

func (p Testmo) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectProjectsMeta,
		tasks.ExtractProjectsMeta,
		tasks.CollectMilestonesMeta,
		tasks.ExtractMilestonesMeta,
		tasks.CollectAutomationRunsMeta,
		tasks.ExtractAutomationRunsMeta,
		tasks.CollectTestsMeta,
		tasks.ExtractTestsMeta,
		tasks.ConvertProjectsMeta,
		tasks.ConvertAutomationRunsMeta,
		tasks.ConvertTestsMeta,
	}
}

func (p Testmo) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeTaskOptions(options)
	if err != nil {
		return nil, err
	}

	logger := taskCtx.GetLogger()
	logger.Debug("%v", options)

	connection := &models.TestmoConnection{}
	connectionHelper := pluginhelper.NewConnectionHelper(
		taskCtx,
		nil,
		p.Name(),
	)

	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, err
	}

	apiClient, err := tasks.CreateTestmoApiClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}

	taskData := &tasks.TestmoTaskData{
		Options:   op,
		ApiClient: apiClient,
	}

	return taskData, nil
}

func (p Testmo) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/testmo"
}

func (p Testmo) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Testmo) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
			"GET":    api.GetConnection,
		},
		"connections/:connectionId/test": {
			"POST": api.TestExistingConnection,
		},
		"connections/:connectionId/remote-scopes": {
			"GET": api.RemoteScopes,
		},
		"connections/:connectionId/search-remote-scopes": {
			"GET": api.SearchRemoteScopes,
		},
		"connections/:connectionId/scopes/:scopeId": {
			"GET":    api.GetScope,
			"PATCH":  api.PatchScope,
			"DELETE": api.DeleteScope,
		},
		"connections/:connectionId/scopes/:scopeId/latest-sync-state": {
			"GET": api.GetScopeLatestSyncState,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScopes,
		},
		"connections/:connectionId/scope-configs": {
			"POST": api.CreateScopeConfig,
			"GET":  api.GetScopeConfigList,
		},
		"connections/:connectionId/scope-configs/:scopeConfigId": {
			"PATCH":  api.UpdateScopeConfig,
			"GET":    api.GetScopeConfig,
			"DELETE": api.DeleteScopeConfig,
		},
	}
}

func (p Testmo) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.TestmoTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}

func (p Testmo) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes)
}
