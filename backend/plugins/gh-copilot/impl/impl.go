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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/api"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/tasks"
)

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginApi
	plugin.PluginModel
	plugin.PluginSource
	plugin.DataSourcePluginBlueprintV200
	plugin.PluginMigration
	plugin.CloseablePluginTask
} = (*Copilot)(nil)

// Copilot is the plugin entrypoint implementing DevLake interfaces.
type Copilot struct{}

func (p Copilot) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)
	return nil
}

func (p Copilot) Description() string {
	return "Collect GitHub Copilot adoption metrics (organization level)"
}

func (p Copilot) Name() string {
	return "gh-copilot"
}

func (p Copilot) Connection() dal.Tabler {
	return &models.CopilotConnection{}
}

func (p Copilot) Scope() plugin.ToolLayerScope {
	return &models.CopilotScope{}
}

func (p Copilot) ScopeConfig() dal.Tabler {
	return nil
}

func (p Copilot) GetTablesInfo() []dal.Tabler {
	return models.GetTablesInfo()
}

func (p Copilot) SubTaskMetas() []plugin.SubTaskMeta {
	return tasks.GetSubTaskMetas()
}

func (p Copilot) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.CopilotOptions
	if err := helper.Decode(options, &op, nil); err != nil {
		return nil, err
	}

	connectionHelper := helper.NewConnectionHelper(taskCtx, nil, p.Name())
	connection := &models.CopilotConnection{}
	if err := connectionHelper.FirstById(connection, op.ConnectionId); err != nil {
		return nil, err
	}

	NormalizeConnection(connection)

	taskData := &tasks.CopilotTaskData{
		Options:    &op,
		Connection: connection,
	}

	return taskData, nil
}

func (p Copilot) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return api.GetApiResources()
}

func (p Copilot) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes)
}

func (p Copilot) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/gh-copilot"
}

func (p Copilot) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Copilot) Close(taskCtx plugin.TaskContext) errors.Error {
	return nil
}
