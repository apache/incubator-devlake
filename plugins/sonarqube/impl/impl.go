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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/sonarqube/api"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/sonarqube/tasks"
)

// make sure interface is implemented
var _ core.PluginMeta = (*Sonarqube)(nil)
var _ core.PluginInit = (*Sonarqube)(nil)
var _ core.PluginTask = (*Sonarqube)(nil)
var _ core.PluginApi = (*Sonarqube)(nil)
var _ core.PluginBlueprintV100 = (*Sonarqube)(nil)
var _ core.CloseablePluginTask = (*Sonarqube)(nil)

type Sonarqube struct{}

func (plugin Sonarqube) Description() string {
	return "collect some Sonarqube data"
}

func (plugin Sonarqube) Init(br core.BasicRes) errors.Error {
	api.Init(br)
	return nil
}

func (plugin Sonarqube) SubTaskMetas() []core.SubTaskMeta {
	// TODO add your sub task here
	return []core.SubTaskMeta{}
}

func (plugin Sonarqube) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	connection := &models.SonarqubeConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Sonarqube connection by the given connection ID")
	}

	apiClient, err := tasks.NewSonarqubeApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Sonarqube API client instance")
	}

	return &tasks.SonarqubeTaskData{
		Options:   op,
		ApiClient: apiClient,
	}, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (plugin Sonarqube) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/sonarqube"
}

func (plugin Sonarqube) MigrationScripts() []core.MigrationScript {
	return migrationscripts.All()
}

func (plugin Sonarqube) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
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

func (plugin Sonarqube) MakePipelinePlan(connectionId uint64, scope []*core.BlueprintScopeV100) (core.PipelinePlan, errors.Error) {
	return api.MakePipelinePlan(plugin.SubTaskMetas(), connectionId, scope)
}

func (plugin Sonarqube) Close(taskCtx core.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.SonarqubeTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
