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
	"github.com/apache/incubator-devlake/plugins/pagerduty/api"
	"github.com/apache/incubator-devlake/plugins/pagerduty/models"
	"github.com/apache/incubator-devlake/plugins/pagerduty/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/pagerduty/tasks"
)

// make sure interface is implemented

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginApi
	plugin.PluginModel
	plugin.DataSourcePluginBlueprintV200
	plugin.CloseablePluginTask
	plugin.PluginSource
} = (*PagerDuty)(nil)

type PagerDuty struct{}

func (p PagerDuty) Description() string {
	return "collect some PagerDuty data"
}

func (p PagerDuty) Name() string {
	return "pagerduty"
}

func (p PagerDuty) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)

	return nil
}

func (p PagerDuty) Connection() dal.Tabler {
	return &models.PagerDutyConnection{}
}

func (p PagerDuty) Scope() plugin.ToolLayerScope {
	return &models.Service{}
}

func (p PagerDuty) ScopeConfig() dal.Tabler {
	return &models.PagerdutyScopeConfig{}
}

func (p PagerDuty) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectIncidentsMeta,
		tasks.ExtractIncidentsMeta,
		tasks.ConvertIncidentsMeta,
		tasks.ConvertServicesMeta,
	}
}

func (p PagerDuty) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.Service{},
		&models.Incident{},
		&models.User{},
		&models.Assignment{},
		&models.PagerDutyConnection{},
		&models.PagerdutyScopeConfig{},
	}
}

func (p PagerDuty) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
		p.Name(),
	)
	connection := &models.PagerDutyConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Pagerduty connection by the given connection ID")
	}

	// Load ScopeConfig from database if only ScopeConfigId is provided
	if op.ScopeConfig == nil && op.ScopeConfigId != 0 {
		var scopeConfig models.PagerdutyScopeConfig
		db := taskCtx.GetDal()
		err = db.First(&scopeConfig, dal.Where("id = ?", op.ScopeConfigId))
		if err != nil && !db.IsErrorNotFound(err) {
			return nil, errors.BadInput.Wrap(err, "fail to get scopeConfig")
		}
		op.ScopeConfig = &scopeConfig
	}
	// Initialize empty ScopeConfig if none provided
	if op.ScopeConfig == nil {
		op.ScopeConfig = new(models.PagerdutyScopeConfig)
	}

	client, err := helper.NewApiClientFromConnection(taskCtx.GetContext(), taskCtx, connection)

	if err != nil {
		return nil, err
	}
	asyncClient, err := helper.CreateAsyncApiClient(taskCtx, client, nil)
	if err != nil {
		return nil, err
	}
	return &tasks.PagerDutyTaskData{
		Options: op,
		Client:  asyncClient,
	}, nil
}

// RootPkgPath information lost when compiled as plugin(.so)
func (p PagerDuty) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/pagerduty"
}

func (p PagerDuty) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p PagerDuty) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
		"connections/:connectionId/search-remote-scopes": {
			"GET": api.SearchRemoteScopes,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScopes,
		},
		"connections/:connectionId/scopes/:scopeId": {
			"GET":    api.GetScope,
			"PATCH":  api.PatchScope,
			"DELETE": api.DeleteScope,
		},
		"connections/:connectionId/scopes/:scopeId/latest-sync-state": {
			"GET": api.GetScopeLatestSyncState,
		},
		"connections/:connectionId/scope-configs": {
			"POST": api.CreateScopeConfig,
			"GET":  api.GetScopeConfigList,
		},
		"connections/:connectionId/scope-configs/:scopeConfigId": {
			"PATCH":  api.PatchScopeConfig,
			"GET":    api.GetScopeConfig,
			"DELETE": api.DeleteScopeConfig,
		},
		"connections/:connectionId/scope-configs/:scopeConfigId/projects": {
			"GET": api.GetServicesByScopeConfig,
		},
	}
}

func (p PagerDuty) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes)
}

func (p PagerDuty) Close(taskCtx plugin.TaskContext) errors.Error {
	_, ok := taskCtx.GetData().(*tasks.PagerDutyTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	return nil
}
