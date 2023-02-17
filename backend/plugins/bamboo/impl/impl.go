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
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
	"github.com/apache/incubator-devlake/plugins/bamboo/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/bamboo/tasks"
)

// make sure interface is implemented
var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginModel
	plugin.PluginMigration
	plugin.PluginBlueprintV100
	plugin.DataSourcePluginBlueprintV200
	plugin.CloseablePluginTask
	plugin.PluginSource
} = (*Bamboo)(nil)

type Bamboo struct{}

func (p Bamboo) Init(br context.BasicRes) errors.Error {
	api.Init(br)
	return nil
}

func (p Bamboo) Connection() interface{} {
	return &models.BambooConnection{}
}

func (p Bamboo) Scope() interface{} {
	return nil
}

func (p Bamboo) TransformationRule() interface{} {
	return nil
}

func (p Bamboo) MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200, syncPolicy plugin.BlueprintSyncPolicy) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
	//return api.MakePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes, &syncPolicy)
	return nil, nil, nil
}

func (p Bamboo) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.BambooConnection{},
	}
}

func (p Bamboo) Description() string {
	return "collect some Bamboo data"
}

func (p Bamboo) SubTaskMetas() []plugin.SubTaskMeta {
	// TODO add your sub task here
	return []plugin.SubTaskMeta{}
}

func (p Bamboo) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	connection := &models.BambooConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Bamboo connection by the given connection ID")
	}

	apiClient, err := tasks.NewBambooApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Bamboo API client instance")
	}

	return &tasks.BambooTaskData{
		Options:   op,
		ApiClient: apiClient,
	}, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (p Bamboo) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/bamboo"
}

func (p Bamboo) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Bamboo) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
	}
}

func (p Bamboo) MakePipelinePlan(connectionId uint64, scope []*plugin.BlueprintScopeV100) (plugin.PipelinePlan, errors.Error) {
	//return api.MakePipelinePlan(p.SubTaskMetas(), connectionId, scope)
	return nil, nil
}

func (p Bamboo) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.BambooTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
