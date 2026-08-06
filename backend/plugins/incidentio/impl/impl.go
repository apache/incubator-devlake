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
	"github.com/apache/incubator-devlake/plugins/incidentio/api"
	"github.com/apache/incubator-devlake/plugins/incidentio/models"
	"github.com/apache/incubator-devlake/plugins/incidentio/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/incidentio/tasks"
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
} = (*Incidentio)(nil)

type Incidentio struct{}

func (p Incidentio) Description() string {
	return "collect incident.io incident data"
}

func (p Incidentio) Name() string {
	return "incidentio"
}

func (p Incidentio) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)
	return nil
}

func (p Incidentio) Connection() dal.Tabler {
	return &models.IncidentioConnection{}
}

func (p Incidentio) Scope() plugin.ToolLayerScope {
	return &models.IncidentType{}
}

func (p Incidentio) ScopeConfig() dal.Tabler {
	return &models.IncidentioScopeConfig{}
}

func (p Incidentio) SubTaskMetas() []plugin.SubTaskMeta {
	// Convert incident types before incidents so the domain Board row
	// exists before the BoardIssue rows that reference it.
	return []plugin.SubTaskMeta{
		tasks.CollectIncidentTypesMeta,
		tasks.ExtractIncidentTypesMeta,
		tasks.CollectIncidentsMeta,
		tasks.ExtractIncidentsMeta,
		tasks.ConvertIncidentTypesMeta,
		tasks.ConvertIncidentsMeta,
	}
}

func (p Incidentio) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.IncidentType{},
		&models.Incident{},
		&models.IncidentioConnection{},
		&models.IncidentioScopeConfig{},
	}
}

func (p Incidentio) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
		p.Name(),
	)
	connection := &models.IncidentioConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get incident.io connection by the given connection ID")
	}

	client, err := helper.NewApiClientFromConnection(taskCtx.GetContext(), taskCtx, connection)
	if err != nil {
		return nil, err
	}
	asyncClient, err := helper.CreateAsyncApiClient(taskCtx, client, nil)
	if err != nil {
		return nil, err
	}
	return &tasks.IncidentioTaskData{
		Options: op,
		Client:  asyncClient,
	}, nil
}

// RootPkgPath information lost when compiled as plugin(.so)
func (p Incidentio) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/incidentio"
}

func (p Incidentio) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Incidentio) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
	}
}

func (p Incidentio) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes)
}

func (p Incidentio) Close(taskCtx plugin.TaskContext) errors.Error {
	_, ok := taskCtx.GetData().(*tasks.IncidentioTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	return nil
}
