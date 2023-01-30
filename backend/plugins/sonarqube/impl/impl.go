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
	"time"

	"github.com/apache/incubator-devlake/plugins/sonarqube/api"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/sonarqube/tasks"
)

// make sure interface is implemented
var _ plugin.PluginMeta = (*Sonarqube)(nil)
var _ plugin.PluginInit = (*Sonarqube)(nil)
var _ plugin.PluginTask = (*Sonarqube)(nil)
var _ plugin.PluginApi = (*Sonarqube)(nil)
var _ plugin.PluginBlueprintV100 = (*Sonarqube)(nil)
var _ plugin.CloseablePluginTask = (*Sonarqube)(nil)

type Sonarqube struct{}

func (p Sonarqube) Description() string {
	return "collect some Sonarqube data"
}

func (p Sonarqube) Init(br context.BasicRes) errors.Error {
	api.Init(br)
	return nil
}

func (p Sonarqube) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectProjectsMeta,
		tasks.ExtractProjectsMeta,
	}
}

func (p Sonarqube) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	logger := taskCtx.GetLogger()
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
	taskData := &tasks.SonarqubeTaskData{
		Options:   op,
		ApiClient: apiClient,
	}
	var createdDateAfter time.Time
	if op.CreatedDateAfter != "" {
		createdDateAfter, err = errors.Convert01(time.Parse(time.RFC3339, op.CreatedDateAfter))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "invalid value for `createdDateAfter`")
		}
	}
	if !createdDateAfter.IsZero() {
		taskData.CreatedDateAfter = &createdDateAfter
		logger.Debug("collect data updated createdDateAfter %s", createdDateAfter)
	}
	return taskData, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (p Sonarqube) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/sonarqube"
}

func (p Sonarqube) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Sonarqube) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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

func (p Sonarqube) MakePipelinePlan(connectionId uint64, scope []*plugin.BlueprintScopeV100) (plugin.PipelinePlan, errors.Error) {
	return api.MakePipelinePlan(p.SubTaskMetas(), connectionId, scope)
}

func (p Sonarqube) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.SonarqubeTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
