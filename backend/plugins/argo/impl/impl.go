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
	"github.com/apache/incubator-devlake/plugins/argo/api"
	"github.com/apache/incubator-devlake/plugins/argo/models"
	"github.com/apache/incubator-devlake/plugins/argo/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/argo/tasks"
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
} = (*Argo)(nil)

type Argo struct{}

func (a Argo) Init(br context.BasicRes) errors.Error {
	api.Init(br, a)

	return nil
}
func (a Argo) Description() string {
	return "Collection Argo data for DORA metrics"
}

func (a Argo) Name() string {
	return "argo"
}

func (a Argo) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/argo"
}

func (a Argo) Connection() dal.Tabler {
	return &models.ArgoConnection{}
}

func (a Argo) Scope() plugin.ToolLayerScope {
	return &models.Project{}
}

func (a Argo) ScopeConfig() dal.Tabler {
	return &models.ArgoScopeConfig{}
}

func (a Argo) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.ArgoConnection{},
		&models.Deployment{},
	}
}

func (a Argo) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectDeploymentsMeta,
		tasks.ExtractDeploymentsMeta,
		//tasks.ConvertDeploymentsMeta,
	}
}

func (a Argo) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.ArgoOptions
	if err := helper.Decode(options, &op, nil); err != nil {
		return nil, err
	}

	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
		a.Name(),
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

func (a Argo) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (a Argo) TestConnection(id uint64) errors.Error {
	_, err := api.TestExistingConnection(helper.GenerateTestingConnectionApiResourceInput(id))
	return err
}

func (a Argo) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
	}
}

func (a Argo) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope,
	skipCollectors bool,
) (pp coreModels.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(a.SubTaskMetas(), connectionId, scopes, skipCollectors)
}

func (a Argo) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.ArgoTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	if data != nil && data.ApiClient != nil {
		data.ApiClient.Release()
	}
	return nil
}
