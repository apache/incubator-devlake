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
	"github.com/apache/incubator-devlake/plugins/feishu/api"
	"github.com/apache/incubator-devlake/plugins/feishu/models"
	"github.com/apache/incubator-devlake/plugins/feishu/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/feishu/tasks"
)

var _ plugin.PluginMeta = (*Feishu)(nil)
var _ plugin.PluginInit = (*Feishu)(nil)
var _ plugin.PluginTask = (*Feishu)(nil)
var _ plugin.PluginApi = (*Feishu)(nil)
var _ plugin.PluginModel = (*Feishu)(nil)
var _ plugin.PluginMigration = (*Feishu)(nil)
var _ plugin.CloseablePluginTask = (*Feishu)(nil)

type Feishu struct{}

func (p Feishu) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (p Feishu) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.FeishuConnection{},
		&models.FeishuMeetingTopUserItem{},
	}
}

func (p Feishu) Description() string {
	return "To collect and enrich data from Feishu"
}

func (p Feishu) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectMeetingTopUserItemMeta,
		tasks.ExtractMeetingTopUserItemMeta,
	}
}

func (p Feishu) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.FeishuOptions
	if err := helper.Decode(options, &op, nil); err != nil {
		return nil, err
	}

	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	connection := &models.FeishuConnection{}
	err := connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, err
	}

	apiClient, err := tasks.NewFeishuApiClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}
	return &tasks.FeishuTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (p Feishu) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/feishu"
}

func (p Feishu) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Feishu) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
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
	}
}

func (p Feishu) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.FeishuTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
