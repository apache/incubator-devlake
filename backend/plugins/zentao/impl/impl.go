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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"github.com/apache/incubator-devlake/plugins/zentao/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/zentao/tasks"
)

// make sure interface is implemented
var _ plugin.PluginMeta = (*Zentao)(nil)
var _ plugin.PluginInit = (*Zentao)(nil)
var _ plugin.PluginTask = (*Zentao)(nil)
var _ plugin.PluginApi = (*Zentao)(nil)

// var _ plugin.CompositePluginBlueprintV200 = (*Zentao)(nil)
var _ plugin.CloseablePluginTask = (*Zentao)(nil)

type Zentao struct{}

func (p Zentao) Description() string {
	return "collect some Zentao data"
}

func (p Zentao) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (p Zentao) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.ConvertProductMeta,
		tasks.ConvertProjectMeta,
		tasks.CollectExecutionMeta,
		tasks.ExtractExecutionMeta,
		tasks.ConvertExecutionMeta,
		tasks.CollectStoryMeta,
		tasks.ExtractStoryMeta,
		tasks.ConvertStoryMeta,
		tasks.CollectBugMeta,
		tasks.ExtractBugMeta,
		tasks.ConvertBugMeta,
		tasks.CollectTaskMeta,
		tasks.ExtractTaskMeta,
		tasks.ConvertTaskMeta,
		tasks.CollectAccountMeta,
		tasks.ExtractAccountMeta,
		tasks.ConvertAccountMeta,
		tasks.CollectDepartmentMeta,
		tasks.ExtractDepartmentMeta,
		tasks.ConvertDepartmentMeta,
	}
}

func (p Zentao) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, errors.Default.Wrap(err, "could not decode Zentao options")
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	connection := &models.ZentaoConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Zentao connection by the given connection ID: %v")
	}

	apiClient, err := tasks.NewZentaoApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Zentao API client instance: %v")
	}

	return &tasks.ZentaoTaskData{
		Options:   op,
		ApiClient: apiClient,
	}, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (p Zentao) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/zentao"
}

func (p Zentao) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Zentao) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
		"connections/:connectionId/product/scopes": {
			"PUT": api.PutProductScope,
		},
		"connections/:connectionId/project/scopes": {
			"PUT": api.PutProjectScope,
		},
		"connections/:connectionId/scopes/product/:scopeId": {
			"GET":   api.GetProductScope,
			"PATCH": api.UpdateProductScope,
		},
		"connections/:connectionId/scopes/project/:scopeId": {
			"GET":   api.GetProjectScope,
			"PATCH": api.UpdateProjectScope,
		},
		"connections/:connectionId/remote-scopes": {
			"GET": api.RemoteScopes,
		},
	}
}

func (p Zentao) MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200, syncPolicy plugin.BlueprintSyncPolicy) (pp plugin.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes, &syncPolicy)
}

func (p Zentao) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.ZentaoTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
