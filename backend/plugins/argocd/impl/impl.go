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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/argocd/api"
	"github.com/apache/incubator-devlake/plugins/argocd/models"
	"github.com/apache/incubator-devlake/plugins/argocd/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/argocd/tasks"
)

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginApi
	plugin.PluginModel
	plugin.PluginMigration
	plugin.PluginSource
	plugin.DataSourcePluginBlueprintV200
	plugin.CloseablePluginTask
} = (*ArgoCD)(nil)

type ArgoCD struct{}

func (p ArgoCD) Connection() dal.Tabler {
	return &models.ArgocdConnection{}
}

func (p ArgoCD) Scope() plugin.ToolLayerScope {
	return &models.ArgocdApplication{}
}

func (p ArgoCD) ScopeConfig() dal.Tabler {
	return &models.ArgocdScopeConfig{}
}

func (p ArgoCD) Description() string {
	return "collect data from ArgoCD"
}

func (p ArgoCD) Name() string {
	return "argocd"
}

func (p ArgoCD) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/argocd"
}

func (p ArgoCD) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)
	return nil
}

func (p ArgoCD) SubTaskMetas() []plugin.SubTaskMeta {
	return tasks.CollectDataTaskMetas()
}

func (p ArgoCD) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	return tasks.PrepareTaskData(taskCtx, options)
}

func (p ArgoCD) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.ArgocdConnection{},
		&models.ArgocdApplication{},
		&models.ArgocdSyncOperation{},
		&models.ArgocdRevisionImage{},
		&models.ArgocdScopeConfig{},
	}
}

func (p ArgoCD) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p ArgoCD) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes)
}

func (p ArgoCD) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.ArgocdTaskData)
	if !ok {
		return errors.Default.New("GetData failed when try to close")
	}
	data.ApiClient.Release()
	return nil
}

func (p ArgoCD) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
		"connections/:connectionId/remote-scopes": {
			"GET": api.RemoteScopes,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScopes,
		},
		"connections/:connectionId/scopes/:scopeId": {
			"GET":    api.GetScope,
			"PATCH":  api.UpdateScope,
			"DELETE": api.DeleteScope,
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
		"scope-config/:scopeConfigId/projects": {
			"GET": api.GetProjectsByScopeConfig,
		},
	}
}

var PluginEntry ArgoCD //nolint
