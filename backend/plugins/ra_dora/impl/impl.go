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
	"github.com/apache/incubator-devlake/plugins/ra_dora/api"
	"github.com/apache/incubator-devlake/plugins/ra_dora/models"
	"github.com/apache/incubator-devlake/plugins/ra_dora/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/ra_dora/tasks"
)

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginApi
	plugin.PluginModel
	plugin.PluginMigration
	plugin.CloseablePluginTask
	plugin.DataSourcePluginBlueprintV200
	plugin.PluginSource
} = (*RaDoraMetrics)(nil)

type RaDoraMetrics struct{}

func (r RaDoraMetrics) Init(br context.BasicRes) errors.Error {
	api.Init(br, r)

	return nil
}
func (r RaDoraMetrics) Description() string {
	return "Collection Argo data for DORA metrics"
}

func (r RaDoraMetrics) Name() string {
	return "ra_dora"
}

func (r RaDoraMetrics) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/ra_dora"
}

func (r RaDoraMetrics) Connection() dal.Tabler {
	return &models.ArgoConnection{}
}

// TODO
func (r RaDoraMetrics) Scope() plugin.ToolLayerScope {
	return nil
}

func (r RaDoraMetrics) ScopeConfig() dal.Tabler {
	return nil
}

func (r RaDoraMetrics) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.ArgoConnection{},
		&models.Deployment{},
	}
}

func (r RaDoraMetrics) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectDeploymentsMeta,
		tasks.ExtractDeploymentsMeta,
		//tasks.ConvertDeploymentsMeta,
	}
}

func (r RaDoraMetrics) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.ArgoOptions
	if err := helper.Decode(options, &op, nil); err != nil {
		return nil, err
	}

	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
		r.Name(),
	)
	connection := &models.ArgoConnection{}
	err := connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, err
	}

	var apiClient *helper.ApiAsyncClient
	syncPolicy := taskCtx.SyncPolicy()
	if !syncPolicy.SkipCollectors {
		newApiClient, err := tasks.NewSlackApiClient(taskCtx, connection)
		if err != nil {
			return nil, err
		}
		apiClient = newApiClient
	}
	return &tasks.ArgoTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (r RaDoraMetrics) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (r RaDoraMetrics) TestConnection(id uint64) errors.Error {
	return nil
}

func (r RaDoraMetrics) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
			"GET":    api.GetConnection,
		},
	}
}

func (p RaDoraMetrics) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope,
	skipCollectors bool,
) (pp coreModels.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes, skipCollectors)
}

func (r RaDoraMetrics) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.ArgoTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	if data != nil && data.ApiClient != nil {
		data.ApiClient.Release()
	}
	return nil
}
