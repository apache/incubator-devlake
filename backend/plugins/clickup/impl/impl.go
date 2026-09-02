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
	"time"

	"github.com/apache/devlake/core/context"
	"github.com/apache/devlake/core/dal"
	"github.com/apache/devlake/core/errors"
	coreModels "github.com/apache/devlake/core/models"
	"github.com/apache/devlake/core/plugin"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/clickup/api"
	"github.com/apache/devlake/plugins/clickup/models"
	"github.com/apache/devlake/plugins/clickup/models/migrationscripts"
	"github.com/apache/devlake/plugins/clickup/tasks"
)

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginApi
	plugin.PluginModel
	plugin.PluginSource
	plugin.PluginMigration
	plugin.CloseablePluginTask
	plugin.DataSourcePluginBlueprintV200
} = (*ClickUp)(nil)

type ClickUp struct{}

func (p ClickUp) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)
	return nil
}

func (p ClickUp) Description() string {
	return "To collect and enrich data from ClickUp"
}

func (p ClickUp) Name() string {
	return "clickup"
}

func (p ClickUp) RootPkgPath() string {
	return "github.com/apache/devlake/plugins/clickup"
}

func (p ClickUp) Connection() dal.Tabler {
	return &models.ClickUpConnection{}
}

func (p ClickUp) Scope() plugin.ToolLayerScope {
	return &models.ClickUpFolder{}
}

func (p ClickUp) ScopeConfig() dal.Tabler {
	return &models.ClickUpScopeConfig{}
}

func (p ClickUp) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

// GetTablesInfo MUST list every model (CI `Test_GetPluginTablesInfo` fails otherwise).
func (p ClickUp) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.ClickUpConnection{},
		&models.ClickUpFolder{},
		&models.ClickUpList{},
		&models.ClickUpScopeConfig{},
		&models.ClickUpUser{},
		&models.ClickUpTask{},
		&models.ClickUpTaskComment{},
	}
}

// SubTaskMetas lists subtasks in dependency order: collect/extract before
// convert; the folder's lists (which classify sprints) before tasks; users
// before tasks (issues reference accounts); the folder-board + sprints before
// board_issues/sprint_issues.
func (p ClickUp) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectListMeta,
		tasks.ExtractListMeta,
		tasks.CollectUserMeta,
		tasks.ExtractUserMeta,
		tasks.CollectTaskMeta,
		tasks.ExtractTaskMeta,
		tasks.ConvertFolderMeta,
		tasks.ConvertSprintMeta,
		tasks.ConvertUserMeta,
		tasks.ConvertTaskMeta,
	}
}

func (p ClickUp) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.ClickUpOptions
	if err := helper.Decode(options, &op, nil); err != nil {
		return nil, errors.Default.Wrap(err, "could not decode ClickUp options")
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("clickup connectionId is invalid")
	}
	if op.FolderId == "" {
		return nil, errors.BadInput.New("clickup folderId is required")
	}

	connection := &models.ClickUpConnection{}
	connectionHelper := helper.NewConnectionHelper(taskCtx, nil, p.Name())
	if err := connectionHelper.FirstById(connection, op.ConnectionId); err != nil {
		return nil, errors.Default.Wrap(err, "error getting connection for ClickUp plugin")
	}

	apiClient, err := tasks.CreateApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to create ClickUp API client")
	}

	// Resolve the scope config. Default to an empty (non-nil) config so subtasks
	// can rely on it being present.
	scopeConfig := &models.ClickUpScopeConfig{}
	if op.ScopeConfigId != 0 {
		if err := taskCtx.GetDal().First(scopeConfig, dal.Where("id = ?", op.ScopeConfigId)); err != nil {
			return nil, errors.Default.Wrap(err, "error getting scope config for ClickUp plugin")
		}
	}

	taskData := &tasks.ClickUpTaskData{
		Options:     &op,
		ApiClient:   apiClient,
		ScopeConfig: scopeConfig,
	}
	if op.TimeAfter != "" {
		timeAfter, errConv := errors.Convert01(time.Parse(time.RFC3339, op.TimeAfter))
		if errConv != nil {
			return nil, errors.BadInput.Wrap(errConv, "invalid timeAfter")
		}
		taskData.TimeAfter = &timeAfter
	}
	return taskData, nil
}

func (p ClickUp) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
		"connections/:connectionId/proxy/rest/*path": {
			"GET": api.Proxy,
		},
		"connections/:connectionId/scope-configs": {
			"POST": api.PostScopeConfig,
			"GET":  api.GetScopeConfigList,
		},
		"connections/:connectionId/scope-configs/:scopeConfigId": {
			"PATCH":  api.PatchScopeConfig,
			"GET":    api.GetScopeConfig,
			"DELETE": api.DeleteScopeConfig,
		},
		"connections/:connectionId/scopes/:scopeId": {
			"GET":    api.GetScope,
			"PATCH":  api.PatchScope,
			"DELETE": api.DeleteScope,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScopes,
		},
		"scope-config/:scopeConfigId/projects": {
			"GET": api.GetProjectsByScopeConfig,
		},
	}
}

func (p ClickUp) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	return api.MakePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes)
}

func (p ClickUp) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.ClickUpTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	if data.ApiClient != nil {
		data.ApiClient.Release()
	}
	return nil
}
