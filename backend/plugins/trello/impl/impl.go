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
	"github.com/apache/incubator-devlake/plugins/trello/api"
	"github.com/apache/incubator-devlake/plugins/trello/models"
	"github.com/apache/incubator-devlake/plugins/trello/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/trello/tasks"
)

var _ interface {
	plugin.PluginTask
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginApi
	plugin.PluginModel
	plugin.PluginMigration
	plugin.CloseablePluginTask
	plugin.DataSourcePluginBlueprintV200
} = (*Trello)(nil)

type Trello struct{}

func (p Trello) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)

	return nil
}

func (p Trello) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.TrelloConnection{},
		&models.TrelloBoard{},
		&models.TrelloList{},
		&models.TrelloCard{},
		&models.TrelloLabel{},
		&models.TrelloMember{},
		&models.TrelloCheckItem{},
		&models.TrelloScopeConfig{},
	}
}

func (p Trello) Description() string {
	return "To collect and enrich data from Trello"
}

func (p Trello) Name() string {
	return "trello"
}

func (p Trello) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectListMeta,
		tasks.ExtractListMeta,

		tasks.CollectCardMeta,
		tasks.ExtractCardMeta,

		tasks.CollectLabelMeta,
		tasks.ExtractLabelMeta,

		tasks.CollectCheckItemMeta,
		tasks.ExtractCheckItemMeta,

		tasks.CollectMemberMeta,
		tasks.ExtractMemberMeta,
	}
}

func (p Trello) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.TrelloOptions
	err := helper.Decode(options, &op, nil)
	if err != nil {
		return nil, errors.Default.Wrap(err, "Trello plugin could not decode options")
	}
	if op.BoardId == "" {
		return nil, errors.BadInput.New("trello boardId is required")
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("trello connectionId is invalid")
	}

	connection := &models.TrelloConnection{}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
		p.Name(),
	)
	err = connectionHelper.FirstById(connection, op.ConnectionId)

	if err != nil {
		return nil, errors.Default.Wrap(err, "error getting connection for Trello plugin")
	}
	apiClient, err := tasks.CreateApiClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}
	return &tasks.TrelloTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (p Trello) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/trello"
}

func (p Trello) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Trello) Connection() dal.Tabler {
	return &models.TrelloConnection{}
}

func (p Trello) Scope() plugin.ToolLayerScope {
	return &models.TrelloBoard{}
}

func (p Trello) ScopeConfig() dal.Tabler {
	return &models.TrelloScopeConfig{}
}

func (p Trello) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"echo": {
			"POST": func(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
				return &plugin.ApiResourceOutput{Body: input.Body}, nil
			},
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
		"connections/:connectionId/scopes/:boardId": {
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

func (p Trello) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	return api.MakePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes)
}

func (p Trello) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.TrelloTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
