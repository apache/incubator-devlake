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

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
	"github.com/apache/incubator-devlake/plugins/clickup/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/clickup/tasks"
)

// make sure interface is implemented
var _ plugin.PluginMeta = (*Clickup)(nil)
var _ plugin.PluginInit = (*Clickup)(nil)
var _ plugin.PluginTask = (*Clickup)(nil)
var _ plugin.PluginApi = (*Clickup)(nil)
var _ plugin.CloseablePluginTask = (*Clickup)(nil)
var _ plugin.DataSourcePluginBlueprintV200 = (*Clickup)(nil)

type Clickup struct{}

func (p Clickup) Description() string {
	return "collect some Clickup data"
}

func (p Clickup) Init(br context.BasicRes) errors.Error {
	api.Init(br)
	return nil
}

func (p Clickup) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectUserMeta,
		tasks.ExtractUserMeta,
		tasks.ConvertUsersMeta,
		tasks.CollectIssueMeta,
		tasks.ExtractIssueMeta,
		tasks.ConvertIssuesMeta,
		tasks.CollectFolderMeta,
		tasks.ExtractFolderMeta,
		tasks.CollectFolderlessListMeta,
		tasks.ExtractFolderlessListMeta,
		tasks.ConvertSprintsMeta,
		tasks.ConvertSprintIssuesMeta,
		tasks.CollectTaskTimeInStatusMeta,
		tasks.ExtractTaskTimeInStatusMeta,
	}
}

func (p Clickup) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.ClickupOptions
	var err errors.Error
	// db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	logger.Debug(fmt.Sprintf("Preparing ClickUp task data with options: %v", op))
	err = helper.Decode(options, &op, nil)
	if err != nil {
		return nil, errors.Default.Wrap(err, "could not decode ClickUp options")
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("ClickUp connectionId is empty")
	}
	connection := &models.ClickupConnection{}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get clickup connection")
	}
	apiClient, err := tasks.NewClickupApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Clickup API client instance")
	}
	taskData := &tasks.ClickupTaskData{
		Options:   &op,
		ApiClient: apiClient,
		TeamId:    connection.TeamId,
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
func (p Clickup) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/clickup"
}

func (p Clickup) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Clickup) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
		"connections/:connectionId/remote-scopes": {
			"GET": api.RemoteScopes,
		},
		"connections/:connectionId/scopes/:scopeId": {
			"GET":   api.GetScope,
			"PATCH": api.UpdateScope,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScope,
		},
	}
}

func (p Clickup) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.ClickupTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}

func (p Clickup) MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200, syncPolicy plugin.BlueprintSyncPolicy) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
	return api.MakePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes, &syncPolicy)
}
