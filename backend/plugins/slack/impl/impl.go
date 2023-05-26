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
	"github.com/apache/incubator-devlake/plugins/slack/api"
	"github.com/apache/incubator-devlake/plugins/slack/models"
	"github.com/apache/incubator-devlake/plugins/slack/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/slack/tasks"
)

var _ plugin.PluginMeta = (*Slack)(nil)
var _ plugin.PluginInit = (*Slack)(nil)
var _ plugin.PluginTask = (*Slack)(nil)
var _ plugin.PluginApi = (*Slack)(nil)
var _ plugin.PluginModel = (*Slack)(nil)
var _ plugin.PluginMigration = (*Slack)(nil)
var _ plugin.CloseablePluginTask = (*Slack)(nil)

type Slack struct{}

func (p Slack) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (p Slack) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.SlackConnection{},
	}
}

func (p Slack) Description() string {
	return "To collect and enrich data from Slack"
}

func (p Slack) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectChannelMeta,
		tasks.ExtractChannelMeta,

		tasks.CollectChannelMessageMeta,
		tasks.ExtractChannelMessageMeta,

		tasks.CollectThreadMeta,
		tasks.ExtractThreadMeta,
	}
}

func (p Slack) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.SlackOptions
	if err := helper.Decode(options, &op, nil); err != nil {
		return nil, err
	}

	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	connection := &models.SlackConnection{}
	err := connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, err
	}

	apiClient, err := tasks.NewSlackApiClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}
	return &tasks.SlackTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (p Slack) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/slack"
}

func (p Slack) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Slack) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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

func (p Slack) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.SlackTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
