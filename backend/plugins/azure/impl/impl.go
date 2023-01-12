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
	"github.com/apache/incubator-devlake/plugins/azure/api"
	"github.com/apache/incubator-devlake/plugins/azure/models"
	"github.com/apache/incubator-devlake/plugins/azure/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/azure/tasks"
)

// make sure interface is implemented
var _ plugin.PluginMeta = (*Azure)(nil)
var _ plugin.PluginInit = (*Azure)(nil)
var _ plugin.PluginTask = (*Azure)(nil)
var _ plugin.PluginApi = (*Azure)(nil)
var _ plugin.PluginModel = (*Azure)(nil)
var _ plugin.CloseablePluginTask = (*Azure)(nil)
var _ plugin.PluginMigration = (*Azure)(nil)

// PluginEntry exports for Framework to search and load
var PluginEntry Azure //nolint

type Azure struct{}

func (p Azure) Description() string {
	return "collect some Azure data"
}

func (p Azure) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (p Azure) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.AzureBuild{},
		&models.AzureBuildDefinition{},
		&models.AzureConnection{},
		&models.AzureRepo{},
	}
}

func (p Azure) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectApiRepoMeta,
		tasks.ExtractApiRepoMeta,
		tasks.CollectApiBuildDefinitionMeta,
		tasks.ExtractApiBuildDefinitionMeta,
	}
}

func (p Azure) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, errors.Default.New("connectionId is invalid")
	}

	connection := &models.AzureConnection{}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, err
	}

	apiClient, err := tasks.CreateApiClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}
	return &tasks.AzureTaskData{
		Options:    op,
		ApiClient:  apiClient,
		Connection: connection,
	}, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (p Azure) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/azure"
}

func (p Azure) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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

func (p Azure) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Azure) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.AzureTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
